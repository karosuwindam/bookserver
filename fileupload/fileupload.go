package fileupload

import (
	"bookserver/dirread"
	"bookserver/message"
	"bookserver/table"
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
	Sql  *table.SQLStatus
	flag bool
}

type UploadFilelist struct {
	// ファイル名
	Name string `json:name`
	// ファイルのサイズ
	Size int64 `json:size`
}

type UploadFileGet struct {
	// 上書き
	Overwrite bool `json:overwrite`
	// データあり
	Register bool `json:register`
	// 書き換え名前
	Name string `json:name`
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

	if err := os.MkdirAll(output.Pdf, 0777); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(output.Zip, 0777); err != nil {
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

// pdfGetList
func (t *UploadPass) pdfGetList() ([]UploadFilelist, error) {
	out := []UploadFilelist{}
	dr, err := dirread.Setup(t.Pdf)
	if err != nil {
		return out, err
	}
	if err := dr.Read("./"); err == nil {
		for _, filedata := range dr.Data {
			tmp := UploadFilelist{Name: filedata.Name, Size: filedata.Size}
			out = append(out, tmp)
		}
	} else {
		return out, err
	}
	return out, nil
}

// zipGetList
func (t *UploadPass) zipGetList() ([]UploadFilelist, error) {
	out := []UploadFilelist{}
	dr, err := dirread.Setup(t.Zip)
	if err != nil {
		return out, err
	}
	if err := dr.Read("./"); err == nil {
		for _, filedata := range dr.Data {
			tmp := UploadFilelist{Name: filedata.Name, Size: filedata.Size}
			out = append(out, tmp)
		}
	} else {
		return out, err
	}
	return out, nil
}

// upload_list
// リスト情報取得
// ToDo
func (t *UploadPass) upload_list(w http.ResponseWriter, r *http.Request) {
	t.rst.Result = ""
	sUrl := urlAnalysis(r.URL.Path)
	urlPoint := 0
	for ; urlPoint < len(sUrl); urlPoint++ {
		if sUrl[urlPoint] == "pdf" {
			break
		} else if sUrl[urlPoint] == "zip" {
			break
		}
	}
	if urlPoint >= len(sUrl) {
		t.rst.Code = http.StatusBadRequest
		return
	} else {
		switch sUrl[urlPoint] {
		case "pdf":
			if rst, err := t.pdfGetList(); err != nil {
				t.rst.Code = http.StatusBadRequest
				t.rst.Result = err.Error()
			} else {
				t.rst.Result = rst
			}
		case "zip":
			if rst, err := t.zipGetList(); err != nil {
				t.rst.Code = http.StatusBadRequest
				t.rst.Result = err.Error()
			} else {
				t.rst.Result = rst
			}
		}
	}
}

// filedbchkeck
// ファイル名の変更確認
func filedbchkeck(str string) (string, bool) {
	//ToDO
	return str, true
}

// upload_get
// 既存ファイルの確認
// ToDo
func (t *UploadPass) upload_get(w http.ResponseWriter, r *http.Request) {
	t.rst.Result = ""
	out := UploadFileGet{}
	sUrl := urlAnalysis(r.URL.Path)
	urlPoint := 0
	for ; urlPoint < len(sUrl); urlPoint++ {
		if sUrl[urlPoint] == "pdf" {
			break
		} else if sUrl[urlPoint] == "zip" {
			break
		}
	}
	if urlPoint+1 >= len(sUrl) || sUrl[urlPoint+1] == "" {
		t.rst.Code = http.StatusBadRequest
		t.rst.Result = out
		return
	} else {
		out.Name, out.Register = filedbchkeck(sUrl[urlPoint+1])
		switch sUrl[urlPoint] {
		case "pdf":
			if rst, err := t.pdfGetList(); err != nil {
				t.rst.Code = http.StatusBadRequest
				t.rst.Result = err.Error()
			} else {
				for _, filename := range rst {
					if filename.Name == sUrl[urlPoint+1] {
						out.Overwrite = true
						break
					}
				}
				t.rst.Result = out
			}
		case "zip":
			if rst, err := t.zipGetList(); err != nil {
				t.rst.Code = http.StatusBadRequest
				t.rst.Result = err.Error()
			} else {
				for _, filename := range rst {
					if filename.Name == sUrl[urlPoint+1] {
						out.Overwrite = true
						break
					}
				}
				t.rst.Result = out
			}
		}
	}

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
	case "GET":
		t.upload_get(w, r)
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
