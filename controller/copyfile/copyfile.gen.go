package copyfile

import (
	"bookserver/config"
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

var shutdown chan bool
var shutdown_back chan bool
var loopflag bool
var copythread chan CopyFIleData

var zippass string
var publicpass string

func Init() error {
	zippass = config.BScfg.Zip
	publicpass = config.BScfg.Public
	if zippass[len(zippass)-1:] != "/" {
		zippass += "/"
	}
	if publicpass[len(publicpass)-1:] != "/" {
		publicpass += "/"
	}
	if err := os.MkdirAll(zippass, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(publicpass, 0777); err != nil {
		return err
	}
	//ループ処理
	shutdown = make(chan bool, 1)
	shutdown_back = make(chan bool, 1)
	copythread = make(chan CopyFIleData, 1)
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	log.Println("info:", "copy file loop start")
	loopflag = true
loop:
	for {
		select {
		case <-shutdown:
			shutdown_back <- true
			break loop
		case data := <-copythread:
			if err := data.AddTable(); err != nil {
				log.Println("error:", err)
			}
		case <-time.After(20 * time.Second):
			//20秒ごとの処理
			//テーブルを確認して、有効なファイルが公開フォルダに登録があるか確認
			if err := ChackCopyFileTableDataAll(); err != nil {
				log.Println("error:", err)
			}
		}

	}
	wg.Wait()
	loopflag = false
	log.Println("info:", "copy file loop end")
	return nil
}

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

// 特定場所からファイルコピーする処理
func Add(id int, flag bool) error {
	select {
	case copythread <- CopyFIleData{
		id:   id,
		flag: flag,
	}:
	case <-time.After(1 * time.Second):
		return errors.New("Stop error timeout")
	}
	return nil
}
