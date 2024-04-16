package fileupload

import (
	"bookserver/config"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var pdfFolder string //pdfファイルのアップロード先
var zipFolder string //zipファイルのアップロード先

var maxMultiMemory int64 //ファイル読み取りの確保メモリ

const (
	MIN_MULTI_MEMORY = 32 << 20  // 32MB
	MAX_MULTI_MEMORY = 512 << 20 // 512MB
)

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

	//マルチアップロードの確保メモリ最大値入力
	if tmp, err := setupMaxMultiMemory(config.BScfg.MAX_MULTI_MEMORY); err != nil || tmp < MIN_MULTI_MEMORY {
		if err != nil {
			log.Println(err)
		}
		if tmp == 0 { //0の場合は予定最大値を設定
			maxMultiMemory = MAX_MULTI_MEMORY
		} else {
			maxMultiMemory = MIN_MULTI_MEMORY
		}
	} else {
		maxMultiMemory = tmp
	}
	mux.HandleFunc("POST "+url, PostUploadFile)
	mux.HandleFunc("GET "+url+"/{filename}", GetUploadFileChangeData)
	mux.HandleFunc("GET "+url+"/{filetype}/{filename}", GetUplodFileCheck)
	return nil
}

// setupMaxMultiMemory(str) = int,error
//
// 設定パラメータから確保するメモリ量を設定
func setupMaxMultiMemory(str string) (int64, error) {
	output := -1
	ary := map[string]int{
		"k": 10,
		"m": 20,
		"g": 30,
		"t": 40,
	}
	for s, j := range ary {
		if i := strings.Index(strings.ToLower(str), s); i > 0 {
			tstr := str[:i]
			if num, err := strconv.Atoi(tstr); err == nil {
				return int64(num) << j, nil
			}
		}
	}
	if num, err := strconv.Atoi(str); err == nil {
		return int64(num), nil
	}
	return int64(output), errors.New("read error data for:" + str)
}
