package fileupload

import (
	"bookserver/message"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
)

type UploadPass struct {
	Pdf  string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"`
	Zip  string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"`
	msg  message.Message
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

//セットアップ
func Setup() (*UploadPass, error) {
	output := &UploadPass{}
	if err := env.Parse(output); err != nil {
		return nil, err
	}
	output.msg = message.Message{Name: "upload", Status: "OK", Code: http.StatusOK}
	output.flag = true
	return output, nil
}

//メッセージのバック
func (t *UploadPass) outputmessage(w http.ResponseWriter) {
	w.WriteHeader(t.msg.Code)
	fmt.Fprintf(w, "%v", t.msg.Output())
}

//名前の確認
func (t *UploadPass) Name() string {
	return t.msg.Name
}

//メッセージの確認
func (t *UploadPass) Message() string {
	return t.msg.Status
}

//アップロード処理

func (t *UploadPass) upload_file(w http.ResponseWriter, r *http.Request) {
	t.msg = message.Message{Status: "OK", Code: http.StatusOK}
	file, fileHeader, e := r.FormFile("file")
	if e != nil {
		t.msg.Status = e.Error()
		t.msg.Code = 202
		t.outputmessage(w)
		return
	}
	defer file.Close()
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
		t.msg.InputMessage(err.Error()+"\t"+"not create file:"+savepass+filename, message.LOGOUTPUT_ON)
		t.outputmessage(w)
		return
	}
	defer fp.Close()
	t.msg.InputMessage("Create File:"+savepass+filename, message.LOGOUTPUT_ON)

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
	t.msg.InputMessage("Create File End", message.LOGOUTPUT_ON)
	t.msg.InputMessage("OK", message.LOGOUTPUT_OFF)
	t.msg.Code = http.StatusOK

}

//リスト情報取得
func (t *UploadPass) upload_list(w http.ResponseWriter, r *http.Request) {

}

//基本処置
func (t *UploadPass) upload_defult(w http.ResponseWriter, r *http.Request) {
	// msg := Message{Status: "OK", Code: 200}
	// bytes, _ := json.Marshal(msg)
	// w.WriteHeader(msg.Code)
	// fmt.Fprintf(w, "%v", string(bytes))
	t.outputmessage(w)
}

//Method別処理
func fIleupload(t *UploadPass, w http.ResponseWriter, r *http.Request) {
	t.msg.Name = "upload"
	t.msg.Code = 200
	switch r.Method {
	case "POST":
		t.upload_file(w, r)
	case "LIST":
		t.upload_list(w, r)
	default:
		t.upload_defult(w, r)
	}
}

func FIleupload(it interface{}, w http.ResponseWriter, r *http.Request) {
	switch it.(type) {
	case *UploadPass:
		t := it.(*UploadPass)
		fIleupload(t, w, r)

	default:
		log.Println("input point type err")
		w.WriteHeader(400)
		fmt.Fprintf(w, "input point type err")
	}

}
