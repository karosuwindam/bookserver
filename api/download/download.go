package download

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/table"
	"bookserver/webserver"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var apiname string = "download" //api名

var zippath string //zipのフォルダパス
var pdfpath string //pdfのフォルダパス
var sql *table.SQLStatus

// dataDownload(filename, pass, typedata, w)
//
// データをダウンロードする処理
//
// filename: ダウンロード時のファイル名
// pass: 開くファイルパス
// typedata: コンテンツの種類
// w: 書き込むhttpファイル
func dataDownload(filename, pass, typedata string, w http.ResponseWriter) {

	fmt.Println("Download pass:", pass)
	if file, err := os.Open(pass); err != nil {
		log.Println(err)
	} else {
		defer file.Close()
		buf := make([]byte, 1024)
		var buffer []byte
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				// Readエラー処理
				break
			}
			buffer = append(buffer, buf[:n]...)
		}
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/"+typedata)
		// ファイルの長さ
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
		// bodyに書き込み
		w.Write(buffer)
	}
}

// webDownload(w, r)
//
// /download/:type/:idの動作
func webDownload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	sUrl := common.UrlAnalysis(r.URL.Path)
	tType := ""
	tId := ""
	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+2 {
			tType = sUrl[i+1]
			tId = sUrl[i+2]
			break
		}
	}
	if tId == "" || tType == "" {
		w.WriteHeader(http.StatusNotFound)
	} else if id, err := strconv.Atoi(tId); err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		if json, err := sql.ReadID(table.FILELIST, id); err != nil {
			w.WriteHeader(http.StatusNotFound)

		} else {
			if data, ok := table.JsonToStruct(table.FILELIST, []byte(json)).([]table.Filelists); ok {
				switch tType {
				case "zip":
					filename := data[0].Zippass
					pass := zippath + filename
					dataDownload(filename, pass, tType, w)
				case "pdf":
					filename := data[0].Pdfpass
					pass := pdfpath + filename
					dataDownload(filename, pass, tType, w)
				default:
					w.WriteHeader(http.StatusNotFound)

				}

			}

		}

	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname + "/", Handler: webDownload},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	zippath = cfg.Folder.Zip
	if zippath[len(zippath)-1:] != "/" {
		zippath += "/"
	}
	pdfpath = cfg.Folder.Pdf
	if pdfpath[len(pdfpath)-1:] != "/" {
		pdfpath += "/"
	}
	if sqlcfg, err := table.Setup(cfg); err != nil {
		return nil, err
	} else {
		sql = sqlcfg
	}
	return route, nil
}
