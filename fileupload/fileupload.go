package fileupload

import (
	"bookserver/message"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type UploadPass struct {
	Pdf  string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"`
	Zip  string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"`
	rst  message.Result
	flag bool
}

// URLの解析
func urlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

// Setup()
// Uploadモジュールを使用するためのメイン処理
func Setup() (*UploadPass, error) {
	output := &UploadPass{}
	if err := env.Parse(output); err != nil {
		return nil, err
	}
	output.rst = message.Result{
		Name:   "upload",
		Code:   http.StatusOK,
		Option: "",
		Date:   time.Now(),
		Result: "OK",
	}
	output.flag = true
	return output, nil
}

//メッセージのバック
func (t *UploadPass) outputmessage(w http.ResponseWriter) {
	w.WriteHeader(t.rst.Code)
	fmt.Fprintf(w, "%v", t.rst.Output())
}

//名前の確認
func (t *UploadPass) Name() string {
	return t.rst.Name
}

// MessageJsonOut
// 保存しているメッセージをjson出力する
func (t *UploadPass) MessageJsonOut() string {
	return t.rst.Output()
}

// upload_file
//
// アップロードのメイン動作

func (t *UploadPass) upload_file(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, e := r.FormFile("file")
	if e != nil {
		log.Println(e.Error())
		t.rst.Code = 202
		t.rst.Result = e.Error()
		return
	}
	defer file.Close()
	t.rst.Option += "," + "file:" + fileHeader.Filename
	filename := fileHeader.Filename
	savepass := ""
	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		savepass = t.Pdf + "/"
	} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
		savepass = t.Zip + "/"
	} else {

	}
	fp, err := os.Create(savepass + filename)
	if err != nil {
		log.Println(err.Error() + "\t" + "not create file:" + savepass + filename)
		t.rst.Result = err.Error() + "\t" + "not create file:" + savepass + filename
		return
	}
	defer fp.Close()
	message.Println("Create File:" + savepass + filename)
	t.rst.Result = "Create File:" + savepass + filename

	var data []byte = make([]byte, 1024)
	var tmplength int64 = 0

	for {
		n, e := file.Read(data)
		if n == 0 {
			break
		}
		if e != nil {
			return
		}
		fp.WriteAt(data, tmplength)
		tmplength += int64(n)
	}
	message.Println("Create File End")
	t.rst.Result = "Create File End"
	t.rst.Result = "OK"
	t.rst.Code = http.StatusOK

}

// upload_list
// リスト情報取得
// ToDo
func (t *UploadPass) upload_list(w http.ResponseWriter, r *http.Request) {

}

// upload_defult
// モジュール動作後に戻す処理
func (t *UploadPass) upload_defult(w http.ResponseWriter, r *http.Request) {
	t.outputmessage(w)
}

// fileupload(t w r)
//
// 本モジュールのhttp動作の判断部分
// Method別処理
func fileupload(t *UploadPass, w http.ResponseWriter, r *http.Request) {
	t.rst.Code = http.StatusOK
	t.rst.Date = time.Now()
	t.rst.Option = r.Method + ":" + r.URL.Host
	switch r.Method {
	case "POST":
		t.upload_file(w, r)
	case "LIST":
		t.upload_list(w, r)
	}
	t.upload_defult(w, r)
}

// Fileupload(interface{},w,r)
//
// webserverでv1のアプリとして登録するための関数
//
// it(interface{}) : 本Uploadの設定ファイル
func FIleupload(it interface{}, w http.ResponseWriter, r *http.Request) {
	switch it.(type) {
	case *UploadPass:
		t := it.(*UploadPass)
		fileupload(t, w, r)

	default:
		log.Println("input point type err")
		w.WriteHeader(400)
		fmt.Fprintf(w, "input point type err")
	}

}
