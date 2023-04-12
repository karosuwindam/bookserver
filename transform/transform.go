package transform

import (
	"bookserver/api/upload"
	"bookserver/config"
	"bookserver/transform/pdftozip"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

func Setup(cfg *config.Config) error {
	ch1 = make(chan interface{}, 5)
	shutdown = make(chan bool)
	if err := pdftozip.SetUp(cfg); err != nil {
		return err
	}
	return nil
}

type PdftoZip struct {
	InputFile  string
	OutputFile string
}

var ch1 chan interface{} //処理に向けてデータを
var shutdown chan bool

// 実行について
func Run(ctx context.Context) {
	fmt.Println("Start: transform loop")
	var wp sync.WaitGroup
	wp.Add(2)
	go func(ctx context.Context) { //uploadからデータの取り出し
		defer wp.Done()
	uploadloop:
		for {
			select {
			case <-ctx.Done():
				break uploadloop
			case <-time.After(time.Microsecond * 100):
				if name, err := upload.GetUploadName(); err == nil {
					fmt.Println("transform send:", name)
					Add(name)
				}
			}
		}
	}(ctx)
	go func(ctx context.Context) { //ch1の処理
		defer wp.Done()
	ch1loop:
		for {
			select {
			case <-ctx.Done():
				break ch1loop
			case tmp := <-ch1:
				switch tmp.(type) {
				case PdftoZip: //PDFをZIPへ変換処理
					data, _ := tmp.(PdftoZip)
					if err := pdftozip.Pdftoimage(data.InputFile, data.OutputFile); err != nil {
						fmt.Println(err)
					}
					fmt.Println("reseav:", data)
				default:
					fmt.Println("transform errdata:", tmp)
				}
			}
		}
		shutdown <- true
	}(ctx)
	wp.Wait()
	log.Println("Close: transform loop")
}

// 処理の追加
func Add(data interface{}) error {
	select {
	case ch1 <- data:
		fmt.Println("transform add:", data)
		return nil
	case <-time.After(1 * time.Second):
		return errors.New("time out")
	}
}

// シャットダウン待ち
func Wait() error {
	select {
	case <-shutdown:
		log.Println("Shutdown transform loop")
		return nil
	case <-time.After(10 * time.Second):
		log.Println("Shutdown transform loop time out")
		return errors.New("time out")

	}
}
