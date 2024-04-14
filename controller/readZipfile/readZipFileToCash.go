package readzipfile

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	CASH_MAX     int     = 512 * 1024 * 1024 //キャッシュの読み取り最大数512M
	CASH_TIMEMAX float64 = 1                 //維持時間:1分
	CASH_CH_MAX  int     = 3                 //キャッシュの読み取りスレッドの最大数
)

type cashZipFile struct { // キャッシュとして保存する構造体
	buf map[string]*bytes.Buffer //ファイル名とbufデータ
}

type cashStore struct { //キャッシュ用のデータストア構造体
	cashZip     map[string]*cashZipFile //キャッシュとして保存するデータ
	cashZipTime map[string]time.Time    //キャッシュの保存時刻
	cashZipSize map[string]int          //キャッシュの保存サイズ
	mu          sync.Mutex              //キャッシュ読み書きのロック
}

var cashch chan bool = make(chan bool, CASH_CH_MAX) //キャッシュ読み取りのスレッド数
var chname chan string                              //キャッシュを作成するzpiファイルについて

var dataStore cashStore

// ReadCashData(filename, zipName string) (*bytes.Buffer, error)
//
// キャッシュ内のデータを読み取る
//
// filename string: Zipファイル名
// zipName string: Zipファイル内の対象ファイル名
// buf *bytes.Buffer:Zipファイル内の対象ファイルデータ
func ReadCashData(filename, zipName string) (*bytes.Buffer, error) {
	if dataStore.cashZip != nil {
		if dataStore.cashZip[filename] != nil {
			renewCashZipTime(filename)
			dataStore.mu.Lock()
			defer dataStore.mu.Unlock()
			if buf := dataStore.cashZip[filename].buf[zipName]; buf != nil {
				return buf, nil
			}
			return nil, errors.New("not cash file for")
		}
	}
	return nil, errors.New("dataStore not setup")
}

// キャッシュの保存データの初期化
func cashInit() error {
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return err
	}
	dataStore.cashZip = map[string]*cashZipFile{}
	dataStore.cashZipTime = map[string]time.Time{}
	dataStore.cashZipSize = map[string]int{}
	chname = make(chan string, 1)
	return nil
}

// createAddCashZip(filename)
//
// ファイル名に一致した空キャッシュを作成する
//
// filename string: キャッシュを作成するファイル名
func createAddCashZip(filename string) {
	if dataStore.cashZip != nil {
		if dataStore.cashZip[filename] == nil {
			dataStore.mu.Lock()
			defer dataStore.mu.Unlock()
			dataStore.cashZip[filename] = &cashZipFile{buf: make(map[string]*bytes.Buffer)}
			dataStore.cashZipSize[filename] = 0
			dataStore.cashZipTime[filename] = time.Time{}
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
	if dataStore.cashZip != nil {
		if dataStore.cashZip[filename] != nil && size != 0 {
			dataStore.mu.Lock()
			defer dataStore.mu.Unlock()
			dataStore.cashZip[filename].buf[zipName] = buf
			dataStore.cashZipSize[filename] += size
		}
	}
}

// renewCashZipTime(filename)
//
// # Zipファイルを読み取った時刻を更新する
//
// filename string:対象となるzipファイル名
func renewCashZipTime(filename string) {
	dataStore.mu.Lock()
	defer dataStore.mu.Unlock()
	dataStore.cashZipTime[filename] = time.Now()
}

// clearAddCashZip(filename)
//
// 対象のファイル名のキャッシュをクリアする
//
// filename string:ファイル名
func clearAddCashZip(filename string) {
	if dataStore.cashZip != nil {
		if dataStore.cashZip[filename] != nil {
			dataStore.mu.Lock()
			defer dataStore.mu.Unlock()
			dataStore.cashZip[filename] = nil
			dataStore.cashZipSize[filename] = 0
			dataStore.cashZipTime[filename] = time.Time{}
		}
	}
}

// readZipFileAll(name) error
//
// 対象のファイル内ファイルをキャッシュに登録する
//
// name string: キャッシュ登録するzipファイル
func readZipFileAll(name string) error {
	if dataStore.cashZip == nil {
		return errors.New("error clear cash")
	}
	if dataStore.cashZip[name] == nil {
		//キャッシュ作成
		createAddCashZip(name)
		pass := zipPath + name

		if _, err := os.Stat(pass); err == nil {
			//zipファイルの読み取り
			if r, err := zip.OpenReader(pass); err != nil {
				return err
			} else {
				defer r.Close()
				select {
				case cashch <- true:
					break
				case <-time.After(1 * time.Second):

					return errors.New("cash over thread for " + fmt.Sprintf("%v > %v", len(cashch), CASH_CH_MAX))
				}
				//キャッシュ容量を確認してクリア
				if getZipFileCashSize() > CASH_MAX {
					clearZipFileCash()
				}
				//キャッシュファイル登録
				for _, f := range r.File {
					buf := new(bytes.Buffer)
					if rc, err := r.Open(f.Name); err != nil {
						log.Println(err)
						continue
					} else {
						if count, err := io.Copy(buf, rc); err != nil {

						} else {
							addCashZip(name, f.Name, int(count), buf)
						}
						rc.Close()

					}
				}
				<-cashch
			}
			//キャッシュの登録時刻更新
			renewCashZipTime(name)
		} else {
			dataStore.cashZip[name] = nil
		}
	} else {
		//キャッシュの登録時刻更新
		renewCashZipTime(name)
	}
	return nil
}

// readZipFileData(zipName, filename string) (*bytes.Buffer, error)
//
// 対象のファイル内でファイル名を指定して読み取る
//
// zipName string: 読み取るzipファイル
// filename string: 読み取るzip内のファイル
func readZipFileData(zipName, filename string) (*bytes.Buffer, error) {

	buf := new(bytes.Buffer)
	pass := zipPath + zipName
	if r, err := zip.OpenReader(pass); err != nil {
		return buf, err
	} else {
		defer r.Close()

		if rc, err := r.Open(filename); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("zip Open: %v", pass))
		} else {
			if _, err := io.Copy(buf, rc); err != nil {
				return buf, errors.Wrap(err, fmt.Sprintf("zip file Copy error: %v %v", pass, filename))
			}
			rc.Close()
		}
	}
	return buf, nil
}

// getZipFileCashSize() = int
//
// 読み込んだファイルキャッシュの合計サイズを戻り値とする
func getZipFileCashSize() int {
	dataStore.mu.Lock()
	defer dataStore.mu.Unlock()
	size := 0
	for _, i := range dataStore.cashZipSize {
		size += i
	}
	return size
}

// clearZipFileCash
//
// 古いキャッシュを削除する
func clearZipFileCash() error {
	size := getZipFileCashSize()
	if size > CASH_MAX {
		str := ""
		bt := time.Time{}
		for s, t := range dataStore.cashZipTime {
			if str == "" {
				str = s
				bt = t
			} else if t.Second() < bt.Second() && t.Second() != 0 {
				str = s
				bt = t
			}
		}
		log.Println("info:", "Clear Cash", str)
		clearAddCashZip(str)
	}
	for s, t := range dataStore.cashZipTime {
		if t.Minute() != 0 {
			if time.Now().Sub(t).Minutes() > CASH_TIMEMAX {
				log.Println("info:", "Clear Cash", s)
				clearAddCashZip(s)
			}
		}
	}
	return nil
}
