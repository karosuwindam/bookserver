package view

import (
	"archive/zip"
	"bookserver/api/common"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	CASH_MAX     int     = 512 * 1024 * 1024
	CASH_TIMEMAX float64 = 1 //維持時間:1分
	CASH_CH_MAX  int     = 3
)

type ZipFile struct { // 読み込んだZipFIleの情報
	name      string
	path      string   //zipファイルパス
	flag      bool     //有効の有無
	DataName  []string `json:"Name"`  //zip内のデータファイル名
	FileCount int      `json:"Count"` //zip内のファイル数
}

// openfile(name) = ZipFile
//
// zipファイルを読み込みその情報を返す
//
// name string: 読み込むzipファイル名
func openfile(name string) ZipFile {
	output := ZipFile{name: name}
	if strings.Index(strings.ToLower(name), ".zip") <= 0 {
		return output
	}
	if !common.Exists(zippath + name) {
		return output
	}
	output.path = zippath + name
	if err := output.createfilelist(); err != nil {
		return output
	}
	output.flag = true

	return output
}

// (*ZipFile) createfilelist() = error
//
// zipファイル内のファイルリストを用意する
func (f *ZipFile) createfilelist() error {
	if f.path == "" {
		return errors.New("path is not input")
	}
	if r, err := zip.OpenReader(f.path); err != nil {
		return err
	} else {
		defer r.Close()
		i := 0
		tmp := []string{}
		for _, f := range r.File {
			tmp = append(tmp, f.Name)
			i++
		}
		f.FileCount = i
		f.DataName = tmp
	}

	return nil
}

// (*ZipFIle) openZipRead(name) = (*bytes.Buffer, error)
//
// zipファイル内のファイルから一致するファイル名の情報を取り出す
//
// name string: zipファイル内のファイル名
func (f *ZipFile) openZipRead(name string) (*bytes.Buffer, error) {
	if b := readZipCash(f.name, name); b != nil {
		return b, nil
	}
	buf := new(bytes.Buffer)
	if !f.flag {
		return buf, errors.New("not chack file")
	}
	if r, err := zip.OpenReader(f.path); err != nil {
		return buf, err
	} else {
		defer r.Close()
		if rc, err := r.Open(name); err != nil {
			return buf, err
		} else {
			_, err := io.Copy(buf, rc)
			if err != nil {
				return buf, err
			}
			rc.Close()
		}
	}
	return buf, nil
}

type cashZipFile struct { // キャッシュとして保存する構造体
	buf map[string]*bytes.Buffer //ファイル名とbufデータ
}

var chname chan string //キャッシュを作成するzpiファイルについて

// Loop(ctx)
//
// キャッシュ作成処理のループ処理
//
// ctx context.Context: ループ処理を制御するcontextの親情報
func Loop(ctx context.Context) {
	if chname == nil {
		return
	}
	fmt.Println("zip chash loop start")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case name := <-chname:
			readZipFileAll(name)
		case <-time.After(time.Second):
			clearZipFileCash()
		}
	}
	close(chname)
	fmt.Println("zip chash loop end")
}

// Add(name) = error
//
// キャッシュを読み込むファイル名を登録する
//
// name string: キャッシュの読み取り対象のファイル名
func Add(name string) error {
	select {
	case chname <- name:
		break
	case <-time.After(time.Second * 5):
		return errors.New("Time Out Add Loop")
	}
	return nil
}

var cashZip map[string]*cashZipFile = nil //キャッシュとして保存するデータ
var cashZipTime map[string]time.Time      //キャッシュの更新時間
var cashZipSize map[string]int            //キャッシュの保存サイズ

var cashch chan bool = make(chan bool, CASH_CH_MAX) //キャッシュ読み取りのスレッド数
var clearMu sync.Mutex                              //キャッシュ読み書きのロック

// createAddCashZip(filename)
//
// ファイル名に一致した空キャッシュを作成する
//
// filename string: キャッシュを作成するファイル名
func createAddCashZip(filename string) {
	if cashZip != nil {
		if cashZip[filename] == nil {
			clearMu.Lock()
			defer clearMu.Unlock()
			cashZip[filename] = &cashZipFile{buf: make(map[string]*bytes.Buffer)}
			cashZipSize[filename] = 0
			cashZipTime[filename] = time.Time{}
		}
	}
}

// addCashZip(filename, zipName, size, buf)
//
// キャッシュにZip内の読み取ったファイルを登録する
//
// filename string: Zipファイル名
// zipName string: Zipファイル内の対象ファイル名
// size int:Zipファイル内の対象ファイルのサイズ
// buf *bytes.Buffer:Zipファイル内の対象ファイルデータ
func addCashZip(filename, zipName string, size int, buf *bytes.Buffer) {
	if cashZip != nil {
		if cashZip[filename] != nil && size != 0 {
			clearMu.Lock()
			defer clearMu.Unlock()
			cashZip[filename].buf[zipName] = buf
			cashZipSize[filename] += size
		}
	}
}

// renewCashZipTime(filename)
//
// # Zipファイルを読み取った時刻を更新する
//
// filename string:対象となるzipファイル名
func renewCashZipTime(filename string) {
	clearMu.Lock()
	defer clearMu.Unlock()
	cashZipTime[filename] = time.Now()

}

// clearAddCashZip(filename)
//
// 対象のファイル名のキャッシュをクリアする
//
// filename string:ファイル名
func clearAddCashZip(filename string) {
	if cashZip != nil {
		if cashZip[filename] != nil {
			clearMu.Lock()
			defer clearMu.Unlock()
			cashZip[filename] = nil
			cashZipSize[filename] = 0
			cashZipTime[filename] = time.Time{}

		}
	}
}

// readZipFileAll(name)
//
// 対象のファイル内ファイルをキャッシュに登録する
//
// name string: キャッシュ登録するzipファイル
func readZipFileAll(name string) {
	if cashZip == nil {
		return
	}
	if cashZip[name] == nil {
		createAddCashZip(name)
		f := openfile(name)
		if f.flag {
			if r, err := zip.OpenReader(f.path); err != nil {
				log.Println(err)
				return
			} else {
				defer r.Close()
				select {
				case cashch <- true:
					break
				case <-time.After(time.Second):
					return
				}

				if getZipFileCashSize() > CASH_MAX {
					clearZipFileCash()
				}
				for _, zipname := range f.DataName {
					buf := new(bytes.Buffer)
					if rc, err := r.Open(zipname); err != nil {
						log.Println(err)
						continue
					} else {
						if f, err := io.Copy(buf, rc); err != nil {

						} else {
							addCashZip(name, zipname, int(f), buf)
						}
						rc.Close()

					}
				}
				<-cashch
			}
			renewCashZipTime(name)
		} else {
			cashZip[name] = nil
		}
	} else {
		renewCashZipTime(name)
	}
}

// readZipCash(zipName, filename) = *bytes.Buffer
//
// キャッシュからファイルを読み取り、キャッシュの時刻を更新する
//
// zipName string: 対象のzipファイル名
// filename string: zip内の対象ファイル
func readZipCash(zipName, filename string) *bytes.Buffer {
	if cashZip != nil {
		clearMu.Lock()
		defer clearMu.Unlock()
		if cashZip[zipName] != nil {
			cashZipTime[zipName] = time.Now()
			return cashZip[zipName].buf[filename]
		}
	}
	return nil
}

// getZipFileCashSize() = int
//
// 読み込んだファイルキャッシュの合計サイズを戻り値とする
func getZipFileCashSize() int {
	clearMu.Lock()
	defer clearMu.Unlock()
	size := 0
	for _, i := range cashZipSize {
		size += i
	}
	return size
}

func clearZipFileCash() {
	size := getZipFileCashSize()
	if size > CASH_MAX {
		str := ""
		bt := time.Time{}
		for s, t := range cashZipTime {
			if str == "" {
				str = s
				bt = t
			} else if t.Second() < bt.Second() && t.Second() != 0 {
				str = s
				bt = t
			}
		}
		fmt.Println("Clear Cash", str)
		clearAddCashZip(str)
	}
	for s, t := range cashZipTime {
		if t.Minute() != 0 {
			if time.Now().Sub(t).Minutes() > CASH_TIMEMAX {
				fmt.Println("Clear Cash", s)
				clearAddCashZip(s)
			}
		}
	}
}

func (f *ZipFile) convertjson() string {
	json, _ := json.Marshal(f)
	return string(json)
}
