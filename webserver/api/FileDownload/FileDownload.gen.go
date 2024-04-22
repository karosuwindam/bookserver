package filedownload

import (
	"bookserver/config"
	"net/http"
	"os"
)

var pdfFolder string //pdfファイルのアップロード先
var zipFolder string //zipファイルのアップロード先

func Init(url string, mux *http.ServeMux) error {
	//アップロードフォルダの存在確認
	pdfFolder = config.BScfg.Pdf
	if pdfFolder[len(pdfFolder)-1:] != "/" {
		pdfFolder += "/"
	}
	zipFolder = config.BScfg.Zip
	if zipFolder[len(zipFolder)-1:] != "/" {
		zipFolder += "/"
	}
	if err := os.MkdirAll(pdfFolder, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(zipFolder, 0777); err != nil {
		return err
	}

	mux.HandleFunc("GET "+url+"/{filetype}/{id}", GetDownload)
	return nil
}
