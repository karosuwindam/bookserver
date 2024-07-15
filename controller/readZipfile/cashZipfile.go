package readzipfile

import (
	"bookserver/config"
	"bytes"
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var zipPath string //アップロードされたzipファイルの保存先
var shutdown chan bool
var shutdown_back chan bool
var loopflag bool

// 初期化処理
func Init() error {
	//zipフォルダの設定
	zipPath = config.BScfg.Zip
	if zipPath[len(zipPath):] != "/" {
		zipPath += "/"
	}
	if err := os.MkdirAll(zipPath, 0777); err != nil {
		return err
	}
	if err := cashInit(); err != nil { //キャッシュ機能の初期化
		return err
	}
	//ループ処理
	shutdown = make(chan bool, 1)
	shutdown_back = make(chan bool, 1)
	return nil
}

// 定期処理
func Run(ctx context.Context) {
	var wg sync.WaitGroup
	log.Println("info:", "zip cash loop start")
	loopflag = true
loop:
	for {
		select {
		case <-shutdown:
			shutdown_back <- true
			break loop
		case filename := <-chname: //キャッシュ作成処理
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				if err := readZipFileAll(filename, context.TODO()); err != nil {
					log.Println("error:", err)
				}
			}(filename)
		case <-time.After(1 * time.Second):
			//１秒ごとの処理
			if err := clearZipFileCash(context.TODO()); err != nil { //キャッシュ定期削除処理
				log.Println("error:", err)
			}
		}
	}
	wg.Wait()
	loopflag = false
	log.Println("info:", "zip cash loop end")

}

// シャットダウン処理
func Stop() error {
	if loopflag {
		shutdown <- true
		select {
		case <-shutdown_back:
			break
		case <-time.After(1 * time.Second):
			return errors.New("Stop error timeout")
		}
	}
	return nil
}

// Zipファイルの中身をキャッシュへ登録する処理
func AddCash(filename string) error {
	select {
	case chname <- filename:
		break
	case <-time.After(1 * time.Second):
		return errors.New("Time Out")
	}

	return nil
}

// Zipファイルの中身を読み取る
func ReadZipfile(filename, zipName string) (*bytes.Buffer, error) {
	buf, err := ReadCashData(filename, zipName)
	if err != nil {
		log.Println("error:", "ReadCashData(", filename, zipName, ") not data cash", err)
		var wg sync.WaitGroup
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			if err := AddCash(filename); err != nil {
				log.Println("error:", err)
			}
		}(filename)
		buf, err = readZipFileData(filename, zipName)
		wg.Wait()
	}
	return buf, err
}
