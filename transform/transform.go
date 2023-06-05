package transform

import (
	"bookserver/api/upload"
	"bookserver/config"
	"bookserver/health/healthmessage"
	"bookserver/table"
	"bookserver/transform/pdftozip"
	"bookserver/transform/writetable"
	"bookserver/transform/ziptopdf"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type fileType int

type fileData struct {
	file     fileType
	fileData string
}

const (
	otherType fileType = 0
	pdf       fileType = 1
	zip       fileType = 1 << 1
)

var fData []fileData = []fileData{
	fileData{file: pdf, fileData: ".pdf"},
	fileData{file: zip, fileData: ".zip"},
}

func Setup(cfg *config.Config) error {
	ch1 = make(chan interface{}, 5)
	ch2 = make(chan table.Filelists, 5)
	shutdown = make(chan bool)
	if err := writetable.Setup(cfg); err != nil {
		return err
	}
	if err := pdftozip.SetUp(cfg); err != nil {
		return err
	}
	if err := ziptopdf.SetUp(cfg); err != nil {
		return err
	}
	return nil
}

var ch1 chan interface{} //処理に向けてデータを
var ch2 chan table.Filelists
var shutdown chan bool
var message healthmessage.HealthMessage //状態記録用

// 動作確認について
func Health() healthmessage.HealthMessage {
	return message
}

// 実行について
func Run(ctx context.Context) {
	hMessage := healthmessage.Create("TransForm Loop")
	if ch1 == nil || ch2 == nil || shutdown == nil {
		return
	}
	hMessage.ChangeMessage("TransForm Loop Start", true)
	message = hMessage.ChangeOut()
	fmt.Println("Start: transform loop")
	var wp sync.WaitGroup
	wp.Add(3)
	go func(ctx context.Context) { //uploadからデータの取り出し
		defer wp.Done()
	uploadloop:
		for {
			select {
			case <-ctx.Done():
				break uploadloop
			case <-time.After(time.Microsecond * 100):

				hMessage.ChangeMessage("Get Upload Data")
				message = hMessage.ChangeOut()
				if name, err := upload.GetUploadName(); err == nil {
					fmt.Println("transform send:", name)
					//ファイルタイプの確認
					switch checkFileType(name) {
					case pdf: //pdfのときの処理
						if outdata, err := writetable.CreatePdfToZip(name); err == nil {
							Add(outdata) //テーブルへ登録処理
						} else {
							log.Println(err)
						}
					case zip: //zipのときの処理
						if outdata, err := writetable.CreateZipToPdf(name); err == nil {
							Add(outdata)
						} else {
							log.Println(err)
						}
					case otherType: //対象外
						log.Println(name, "is not ")
					}
				}
				hMessage.ChangeMessage("OK")
				message = hMessage.ChangeOut()
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
				case writetable.PdftoZip: //PDFをZIPへ変換処理
					hMessage.ChangeMessage("Change PDF to Zip")
					message = hMessage.ChangeOut()
					data, _ := tmp.(writetable.PdftoZip)
					ch2 <- table.Filelists{Name: data.Name, Pdfpass: data.InputFile, Zippass: data.OutputFile, Tag: data.Tag}
					//PDFからZIPを作成する処理
					if err := pdftozip.Pdftoimage(data.InputFile, data.OutputFile); err != nil {
						fmt.Println(err)
					}
					fmt.Println("reseav:", data)
					hMessage.ChangeMessage("OK")
					message = hMessage.ChangeOut()
				case writetable.ZipToPdf: //ZIPをPDFへ変換処理
					hMessage.ChangeMessage("Change Zip to PDF")
					message = hMessage.ChangeOut()
					data, _ := tmp.(writetable.ZipToPdf)
					ch2 <- table.Filelists{Name: data.Name, Zippass: data.InputFile, Pdfpass: data.OutputFile, Tag: data.Tag}
					//ZIPからPDFを作成する処理
					// To Do
					fmt.Println("reseav:", data)
					hMessage.ChangeMessage("OK")
					message = hMessage.ChangeOut()

				default:
					fmt.Println("transform errdata:", tmp)
				}
			}
		}
	}(ctx)
	go func(ctx context.Context) { //ch1の処理
		defer wp.Done()
	ch2loop:
		for {
			select {
			case <-ctx.Done():
				break ch2loop
			case tmp := <-ch2:
				hMessage.ChangeMessage("Table Renew Data")
				message = hMessage.ChangeOut()
				//テーブルへデータを登録する
				if err := writetable.AddFileTable(&tmp); err != nil {
					log.Println(err)
				}
				hMessage.ChangeMessage("OK")
				message = hMessage.ChangeOut()
			}
		}
	}(ctx)
	wp.Wait()
	hMessage.ChangeMessage("TransForm Loop End", false)
	message = hMessage.ChangeOut()
	log.Println("Close: transform loop")
	shutdown <- true
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

// ファイルタイプの確認処理
func checkFileType(name string) fileType {
	for _, f := range fData {
		if strings.Index(strings.ToLower(name), f.fileData) > 0 {
			return f.file
		}
	}
	return otherType
}
