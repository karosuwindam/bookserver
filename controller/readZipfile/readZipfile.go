package readzipfile

import (
	"archive/zip"
	"bookserver/config"
	"bookserver/table/filelists"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func GetZiplist(ctx context.Context) (ZipFile, error) {
	var output ZipFile

	ctx, span := config.TracerS(ctx, "GetZiplist", "Get Zip List")
	defer span.End()
	id, ok := contextReadZipId(ctx)
	if !ok {
		err := fmt.Errorf("context Not Input id")
		config.TracerError(span, err)
		return output, err
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("GetZiplist id=%v", id),
		"id", id,
	)
	//idを元にテーブルからzipのファイル名を取得
	if d, err := filelists.GetId(id); err != nil {
		config.TracerError(span, err)
		return output, errors.Wrap(err, "sql read by booknamaes")
	} else {
		//zipファイル名からファイルリストを作成
		lists, err := openfile(d.Zippass)
		if err != nil {
			config.TracerError(span, err)

			return output, errors.Wrap(err, "read zipfile error")

		}
		//zipフォルダのキャッシュを作成依頼
		if err := AddCash(d.Zippass); err != nil {
			config.TracerError(span, err)
			slog.ErrorContext(ctx,
				fmt.Sprintf("GetZiplist AddChash file=%v", d.Zippass),
				"file", d.Zippass,
				"Error", err,
			)
		}
		output = lists

	}
	return output, nil
}

type ZipFile struct { // 読み込んだZipFIleの情報
	DataName       []string `json:"Name"`     //zip内のデータファイル名
	FileCount      int      `json:"Count"`    //zip内のファイル数
	ImageFileCount int      `json:"ImgCount"` //zip内のイメージファイルカウント
}

// openfile(name) = ZipFile
//
// zipファイルを読み込みファイルリストと数を返す
//
// filepass string: 読み込むzipファイル名
func openfile(filepass string) (ZipFile, error) {
	output := ZipFile{}
	pass := zipPath + filepass
	if strings.Index(strings.ToLower(filepass), ".zip") <= 0 {
		return output, errors.New("not zip file")
	}
	if _, err := os.Stat(pass); err != nil {
		return output, errors.Wrap(err, "Not file")
	}
	if r, err := zip.OpenReader(pass); err != nil {
		return output, errors.Wrap(err, "Zip not file")
	} else {
		defer r.Close()
		i := 0
		j := 0
		tmp := []string{}
		for _, f := range r.File {
			tmp = append(tmp, f.Name)
			j += ImageCheck(f.Name)
			i++
		}
		output.FileCount = i
		output.ImageFileCount = j
		output.DataName = tmp
	}

	return output, nil
}

func ImageCheck(filename string) int {
	imgType := []string{".jpg", ".png"}
	tmp := strings.ToLower(filename)
	for _, imgT := range imgType {
		if strings.Index(tmp, imgT) > 0 {
			return 1
		}
	}
	return 0
}
