package transform

import (
	"bookserver/api/upload"
	"bookserver/config"
	"bookserver/table"
	"bookserver/transform/pdftozip"
	"bookserver/transform/writetable"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// sqlのパスのセットアップ
// func sqlSetup(cfg *config.Config) (*table.SQLStatus, error) {
// 	var err error
// 	if sqlcfg, err := table.Setup(cfg); err == nil {
// 		return sqlcfg, err
// 	}
// 	return nil, err
// }

// var sql *table.SQLStatus

func Setup(cfg *config.Config) error {
	ch1 = make(chan interface{}, 5)
	ch2 = make(chan table.Filelists, 5)
	shutdown = make(chan bool)
	// if sqlcfg, err := sqlSetup(cfg); err == nil {
	// 	sql = sqlcfg
	// } else {
	// 	return err
	// }
	if err := writetable.Setup(cfg); err != nil {
		return err
	}
	if err := pdftozip.SetUp(cfg); err != nil {
		return err
	}
	return nil
}

type PdftoZip writetable.PdftoZip

var ch1 chan interface{} //処理に向けてデータを
var ch2 chan table.Filelists
var shutdown chan bool

// // :ToDo writetableに移動する予定
// // テーブルからファイル名の用意
// func createOutFileName(tabledata table.Booknames, count int) (string, string) {
// 	tmpname := ""
// 	tmptag := ""
// 	tmp := []string{tabledata.Title, tabledata.Writer, tabledata.Booktype, tabledata.Ext}
// 	if tabledata.Title != "" {
// 		tmpname = tabledata.Title
// 	}
// 	if tabledata.Writer != "" && tmpname != "" {
// 		tmpname = "[" + tabledata.Writer + "]" + tmpname
// 	}
// 	for _, str := range tmp {
// 		if str != "" {
// 			if tmptag == "" {
// 				tmptag = str
// 			} else {
// 				tmptag += "," + str
// 			}
// 		}
// 	}
// 	if tmpname == "" {
// 		return tmpname, tmptag
// 	}
// 	if count <= -1 {
// 		return fmt.Sprintf("%s.zip", tmpname), tmptag
// 	}
// 	return fmt.Sprintf("%s%02d.zip", tmpname, count), tmptag
// }

// 実行について
func Run(ctx context.Context) {
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
				if name, err := upload.GetUploadName(); err == nil {
					fmt.Println("transform send:", name)
					// //テーブルから作成zipファイル名を作成
					// //:ToDo

					// if i := strings.Index(name, ".pdf"); i > 0 {
					// 	outname := ""
					// 	tmpname := name[:i]
					// 	tag := ""
					// 	if jdata, err := sql.ReadName(table.BOOKNAME, name[:i]); err == nil && jdata != "[]" {
					// 		if jout, ok := table.JsonToStruct(table.BOOKNAME, []byte(jdata)).([]table.Booknames); ok {
					// 			outname, tag = createOutFileName(jout[0], -1)
					// 		}
					// 	} else {
					// 		count := -1
					// 		for j := 3; j > 0; j-- {
					// 			tt := tmpname[len(tmpname)-j:]
					// 			if cc, err := strconv.Atoi(tt); err == nil {
					// 				count = cc
					// 				tmpname = tmpname[:len(tmpname)-j]
					// 				break
					// 			}
					// 		}
					// 		if jdata, err := sql.ReadName(table.BOOKNAME, tmpname); err == nil && jdata != "[]" {
					// 			if jout, ok := table.JsonToStruct(table.BOOKNAME, []byte(jdata)).([]table.Booknames); ok {
					// 				outname, tag = createOutFileName(jout[0], count)
					// 			}
					// 		}
					// 	}
					// 	if outname == "" {
					// 		outname = name[:i] + ".zip"
					// 	}
					// 	outdata := PdftoZip{Name: name[:i], InputFile: name, OutputFile: outname, Tag: tag}
					if outdata, err := writetable.CreatePdfToZip(name); err == nil {
						Add(outdata)

					} else {
						log.Println(err)
					}
					// }
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
					ch2 <- table.Filelists{Name: data.Name, Pdfpass: data.InputFile, Zippass: data.OutputFile, Tag: data.Tag}
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
	go func(ctx context.Context) { //ch1の処理
		defer wp.Done()
	ch2loop:
		for {
			select {
			case <-ctx.Done():
				break ch2loop
			case tmp := <-ch2:
				if err := writetable.AddFileTable(&tmp); err != nil {
					log.Println(err)
				}
			}
		}
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