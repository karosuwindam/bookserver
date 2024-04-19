package convert

import (
	"bookserver/config"
	"bookserver/controller/convert/pdftozip"
	"bookserver/controller/convert/ziptopdf"
	"bookserver/table/uploadtmp"
	"fmt"
	"log"
)

const ERROR_COUNT_MAX = 3

var count_error int

func Init() error {
	if err := pdftozip.Init(); err != nil {
		return err
	}
	if err := ziptopdf.Init(); err != nil {
		return err
	}
	if err := DataStoreInit(); err != nil {
		return err
	}
	count_error = config.BScfg.ConvertCountMax
	if count_error == 0 {
		count_error = ERROR_COUNT_MAX
	}
	return nil

}

func CheckCovertData() error {
	//処理中にステータスを切り替える
	statusData.On()
	defer statusData.Off()
	statusData.Clear()
	//uploadtmpテーブルから未処理のリストを取得
	lists, err := uploadtmp.GetAllbyFlag(false)
	if err != nil {
		return err
	}
	//処理予定ファイルを登録
	for _, list := range lists {
		statusData.Add(list.Name)
	}
	for _, list := range lists {
		if list.SavePdf != "" && list.SaveZip == "" { //ZipからPDFへ変換する
			if err := pdftozip.ConvertPdfToZip(list.Name); err != nil {
				log.Panicln("error:", err)
			} else {
				//コンバート成功時の処理
				data, _ := pdftozip.ConvertPdfToZipChack(list.Name)
				if err := list.SetZipPath(data.Zippass); err != nil {
					log.Println("debug:", err)
				}
			}
		}
		if list.SavePdf == "" && list.SaveZip != "" { //PDFからZipへ変換する

		}
		//変換処理が無事に完了したらuolodtmpのテーブルを処理済みにする
		if list.SavePdf != "" && list.SaveZip != "" { //正常処理済みなので完了させる
			statusData.Change(list.Name) //ステータスの切り替え
			if err := list.FlagOn(); err != nil {
				log.Println("error:", err)
			}
		} else {
			if err := list.CountUp(); err != nil {
				log.Println("error:", err)
			} else {
				if list.Count > count_error {
					log.Println("info:", fmt.Sprintf("error Count OVER %v by file %v", count_error, list.Name))
					if err := list.FlagOn(); err != nil {
						log.Println("error:", err)
					}
				}
			}
		}

	}
	return nil
}
