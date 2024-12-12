package convert

import (
	"bookserver/config"
	"bookserver/controller/convert/pdftozip"
	"bookserver/controller/convert/ziptopdf"
	"bookserver/table/uploadtmp"
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx := context.TODO()
	ctx, span := config.TracerS(ctx, "Convet", "CheckCovertData")
	defer span.End()
	slog.DebugContext(ctx,
		"CheckCovertData",
	)
	//処理中にステータスを切り替える
	statusData.On()
	defer statusData.Off()
	statusData.Clear()
	//uploadtmpテーブルから未処理のリストを取得
	lists, err := uploadtmp.GetAllbyFlag(false)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	//処理予定ファイルを登録
	for _, list := range lists {
		statusData.Add(list.Name)
	}
	for _, list := range lists {

		ctx, spanList := config.TracerS(ctx, "CheckCovertData", "Listdata")
		defer spanList.End()
		spanList.SetAttributes(attribute.String("Pdf", list.SavePdf))
		spanList.SetAttributes(attribute.String("Zip", list.SaveZip))
		slog.DebugContext(ctx,
			fmt.Sprintf("CheckCovertData file=%v", list.Name),
		)

		if list.SavePdf != "" && list.SaveZip == "" { //ZipからPDFへ変換する
			_, spanConvert := config.TracerS(ctx, "CheckCovertData", "ConvertPdfToZip")
			defer spanConvert.End()
			if err := pdftozip.ConvertPdfToZip(list.Name); err != nil {
				spanConvert.SetStatus(codes.Error, err.Error())
				slog.ErrorContext(ctx,
					fmt.Sprintf("ConvertPdfToZip file=%v", list.Name),
					"file", list.Name,
					"Error", err,
				)

			} else {
				//コンバート成功時の処理
				data, _ := pdftozip.ConvertPdfToZipChack(list.Name)
				if err := list.SetZipPath(data.Zippass); err != nil {
					spanConvert.SetStatus(codes.Error, err.Error())
					slog.WarnContext(ctx,
						fmt.Sprintf("ConvertPdfToZip SetZipPath file=%v", list.Name),
						"file", list.Name,
						"Error", err,
					)
				}
				spanConvert.SetAttributes(attribute.String("Name", data.Name))
				spanConvert.SetAttributes(attribute.String("Pdf", data.Pdfpass))
				spanConvert.SetAttributes(attribute.String("Zip", data.Zippass))
				spanConvert.SetAttributes(attribute.String("Tag", data.Tag))
			}
		}
		if list.SavePdf == "" && list.SaveZip != "" { //PDFからZipへ変換する

		}
		//変換処理が無事に完了したらuolodtmpのテーブルを処理済みにする
		if list.SavePdf != "" && list.SaveZip != "" { //正常処理済みなので完了させる
			statusData.Change(list.Name) //ステータスの切り替え
			if err := list.FlagOn(); err != nil {
				slog.ErrorContext(ctx,
					fmt.Sprintf("FlagOn file=%v", list.Name),
					"file", list.Name,
					"Error", err,
				)
			}
		} else {
			if err := list.CountUp(); err != nil {
				slog.ErrorContext(ctx,
					fmt.Sprintf("CountUp file=%v", list.Name),
					"file", list.Name,
					"Error", err,
				)
			} else {
				if list.Count > count_error {
					slog.InfoContext(ctx,
						fmt.Sprintf("error Count OVER %v by file %v", count_error, list.Name),
						"file", list.Name,
						"count", list.Count,
						"count_error", count_error,
					)
					if err := list.FlagOn(); err != nil {
						slog.ErrorContext(ctx,
							fmt.Sprintf("FlagOn file=%v", list.Name),
							"file", list.Name,
							"Error", err,
						)
					}
				}
			}
		}

	}
	return nil
}
