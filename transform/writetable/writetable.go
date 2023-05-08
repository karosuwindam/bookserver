package writetable

import (
	"bookserver/config"
	"bookserver/table"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	PDF string = ".pdf" //対象となる拡張子
	ZIP string = ".zip"
)

// PdftoZip
//
// データ登録用で、pdfからzipファイルを作成する
type PdftoZip struct {
	Name       string `json:"Name"` //登録用の名称(巻数情報も含む)
	InputFile  string `json:"Pdf"`  //入力ファイル(pdf)
	OutputFile string `json:"Zip"`  //出力ファイル(zip)
	Tag        string `json:"Tag"`  //検索用のタグ情報
}

// createOutFileNmae(tabledata , count) = (string, string)
//
// # Booknamesのデータから、Zip出力用のファイル名とタグ情報を作る
//
// tabledata: Booknamesの配列データ
// count: 巻数データ -1:のときはつけない設定
func createOutFileNmae(tabledata *table.Booknames, count int) (string, string) {
	tmpname := ""
	tmptag := ""
	tmpcount := ""
	if tabledata == nil {
		return tmpname, tmptag
	}
	if count > -1 {
		tmpcount = fmt.Sprintf("%02d", count)
	}
	tmp := []string{tabledata.Title + tmpcount, tabledata.Writer, tabledata.Burand, tabledata.Booktype, tabledata.Ext}
	if tabledata.Title != "" {
		tmpname = tabledata.Title
	}
	if tabledata.Writer != "" && tmpname != "" {
		tmpname = "[" + tabledata.Writer + "]" + tmpname
	}
	for _, str := range tmp {
		if str != "" {
			if tmptag == "" {
				tmptag = str
			} else {
				tmptag += "," + str
			}
		}
	}
	if tmpname == "" {
		return tmpname, tmptag
	}
	if tmptag == "" {
		tmptag = tmpname + tmpcount
	}
	if count <= -1 {
		return fmt.Sprintf("%s.zip", tmpname), tmptag
	}
	return fmt.Sprintf("%s%02d.zip", tmpname, count), tmptag
}

// createBooknamesCount(name) = *table.Booknames
//
// ファイル名からBooknamesのデータ抽出と巻数情報の取り出し
//
// name: Booknamesからname列の文字列と一致する情報の取り出し
func createBooknamesCount(name string) *table.Booknames {
	if jdata, err := sql.ReadName(table.BOOKNAME, name); err == nil && jdata != "[]" {
		if jout, ok := table.JsonToStruct(table.BOOKNAME, []byte(jdata)).([]table.Booknames); ok {
			return &jout[0]
		}
	}
	return nil
}

// CreatePdfToZip(name) = PdfToZip, error
//
// ファイル名からPdftozipの情報を作成する。
func CreatePdfToZip(name string) (PdftoZip, error) {
	var output PdftoZip
	if name == "" {
		return output, errors.New("name is not data")
	}
	output.InputFile = name
	tmpname := strings.ToLower(name)
	if i := strings.Index(tmpname, PDF); i > 0 {
		tmpname = name[:i]
		output.Name = tmpname
		count := -1
		if tmpst := createBooknamesCount(tmpname); tmpst != nil {
			output.OutputFile, output.Tag = createOutFileNmae(tmpst, count)
		} else {
			for j := 3; j > 0; j-- {
				if len(tmpname) < j {
					continue
				}
				tt := tmpname[len(tmpname)-j:]
				if cc, err := strconv.Atoi(tt); err == nil {
					count = cc
					tmpname = tmpname[:len(tmpname)-j]
					break
				}
			}
			if tmpst := createBooknamesCount(tmpname); tmpst != nil {
				output.OutputFile, output.Tag = createOutFileNmae(tmpst, count)
			}
		}
		if output.OutputFile == "" {
			output.OutputFile = name[:i] + ZIP
			output.Tag = name[:i]
		}
	} else {
		return output, errors.New("input name is not PDF")
	}
	return output, nil
}

// 既存ファイル名チェック
func AddFileTable(tmp *table.Filelists) error {
	if jout, err := sql.ReadName(table.FILELIST, tmp.Name); err == nil && jout != "[]" {
		if tmpr, ok := table.JsonToStruct(table.FILELIST, []byte(jout)).([]table.Filelists); ok {
			tmp.Id = tmpr[0].Id
			if _, err := sql.Edit(table.FILELIST, tmp, tmp.Id); err != nil {
				return err
			}
		} else {
			return errors.New("file read error")
		}
	} else {
		if err := sql.Add(table.FILELIST, tmp); err != nil {
			return err
		}

	}
	return nil
}

// sqlのパスのセットアップ
func sqlSetup(cfg *config.Config) (*table.SQLStatus, error) {
	var err error
	if sqlcfg, err := table.Setup(cfg); err == nil {
		return sqlcfg, err
	}
	return nil, err
}

var sql *table.SQLStatus

func Setup(cfg *config.Config) error {

	if sqlcfg, err := sqlSetup(cfg); err == nil {
		sql = sqlcfg
	} else {
		return err
	}
	return nil
}
