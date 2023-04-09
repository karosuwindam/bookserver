package upload

import (
	"bookserver/api/common"
	"bookserver/dirread"
	"bookserver/table"
	"bookserver/webserverv2"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

var apiname string = "upload" //api名

type UploadPass struct {
	Pdf  string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"` //PDFのアップロード先フォルダ
	Zip  string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"` //ZIPのアップロード先フォルダ
	Sql  *table.SQLStatus
	flag bool
}

type UploadFilelist struct {
	// ファイル名
	Name string `json:"Name"`
	// ファイルのサイズ
	Size int64 `json:"Size"`
}

type UploadFileGet struct {
	// // 上書き
	// Overwrite bool `json:overwrite`
	// データあり
	Register bool `json:"Register"`
	// 書き換え名前
	Name string `json:"Name"`
}

// upload_file
//
// アップロードのメイン動作
func upload_file(w http.ResponseWriter, r *http.Request) common.Result {
	file, fileHeader, e := r.FormFile("file")
	msg := common.Result{Code: http.StatusOK, Date: time.Now()}
	if e != nil {
		log.Println(e.Error())
		msg.Code = http.StatusAccepted
		msg.Result = e.Error()
		return msg
	}
	defer file.Close()
	msg.Option = r.Method + ":" + "file:" + fileHeader.Filename
	filename := fileHeader.Filename
	savepass := ""
	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		savepass = setupdata.Pdf + "/"
	} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
		savepass = setupdata.Zip + "/"
	} else {

	}
	fp, err := os.Create(savepass + filename)
	if err != nil {
		log.Println(err.Error() + "\t" + "not create file:" + savepass + filename)
		msg.Result = err.Error() + "\t" + "not create file:" + savepass + filename
		return msg
	}
	defer fp.Close()
	fmt.Println("Create File:" + savepass + filename)
	msg.Result = "Create File:" + savepass + filename

	var data []byte = make([]byte, 1024)
	var tmplength int64 = 0

	for {
		n, e := file.Read(data)
		if n == 0 {
			break
		}
		if e != nil {
			return msg
		}
		fp.WriteAt(data, tmplength)
		tmplength += int64(n)
	}
	fmt.Println("Create File End")
	msg.Result = "OK"
	return msg
}

// pdfGetList
func pdfGetList() ([]UploadFilelist, error) {
	out := []UploadFilelist{}
	dr, err := dirread.Setup(setupdata.Pdf)
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
func zipGetList() ([]UploadFilelist, error) {
	out := []UploadFilelist{}
	dr, err := dirread.Setup(setupdata.Zip)
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
func upload_list(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	sUrl := common.UrlAnalysis(r.URL.Path)
	for i, url := range sUrl {
		if url == apiname {
			switch sUrl[i+1] {
			case "pdf":
				if rst, err := pdfGetList(); err != nil {
					msg.Code = http.StatusBadRequest
					msg.Result = err.Error()
				} else {
					msg.Result = rst
				}
			case "zip":
				if rst, err := zipGetList(); err != nil {
					msg.Code = http.StatusBadRequest
					msg.Result = err.Error()
				} else {
					msg.Result = rst
				}
			}
			break
		}
	}
	return msg
}

// upload_get
// 既存ファイルの確認
// ToDo
func upload_get(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	//jsonデータの取り出し
	var out UploadFileGet
	if err := json.NewDecoder(r.Body).Decode(&out); err != nil {
		msg.Result = err.Error()
		msg.Code = http.StatusBadRequest
		return msg
	}
	//jsonの名前から存在確認
	filename := out.Name
	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		// savepass = setupdata.Pdf + "/"
		out.Register = common.Exists(setupdata.Pdf + "/" + filename)
	} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
		out.Register = common.Exists(setupdata.Zip + "/" + filename)
		// savepass = setupdata.Zip + "/"
	} else {
		msg.Result = errors.New("input pass filename:" + filename).Error()
		msg.Code = http.StatusBadRequest
		return msg
	}
	msg.Result = out
	return msg
	// sUrl := urlAnalysis(r.URL.Path)
	// urlPoint := 0
	// for ; urlPoint < len(sUrl); urlPoint++ {
	// 	if sUrl[urlPoint] == "pdf" {
	// 		break
	// 	} else if sUrl[urlPoint] == "zip" {
	// 		break
	// 	}
	// }
	// if urlPoint+1 >= len(sUrl) || sUrl[urlPoint+1] == "" {
	// 	t.rst.Code = http.StatusBadRequest
	// 	t.rst.Result = out
	// 	return
	// } else {
	// 	out.Name, out.Register = filedbchkeck(sUrl[urlPoint+1])
	// 	switch sUrl[urlPoint] {
	// 	case "pdf":
	// 		if rst, err := t.pdfGetList(); err != nil {
	// 			t.rst.Code = http.StatusBadRequest
	// 			t.rst.Result = err.Error()
	// 		} else {
	// 			for _, filename := range rst {
	// 				if filename.Name == sUrl[urlPoint+1] {
	// 					out.Overwrite = true
	// 					break
	// 				}
	// 			}
	// 			t.rst.Result = out
	// 		}
	// 	case "zip":
	// 		if rst, err := t.zipGetList(); err != nil {
	// 			t.rst.Code = http.StatusBadRequest
	// 			t.rst.Result = err.Error()
	// 		} else {
	// 			for _, filename := range rst {
	// 				if filename.Name == sUrl[urlPoint+1] {
	// 					out.Overwrite = true
	// 					break
	// 				}
	// 			}
	// 			t.rst.Result = out
	// 		}
	// 	}
	// }

}

// fuileupload(w r)
// 本モジュールのhttp動作の判断部分
// Method別処理
// /update の処理
func fileupload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	if common.CkLogin(&msg, w, r) {
		switch r.Method {
		case "POST":
			msg = upload_file(w, r)
		case "LIST":
			msg = upload_list(w, r)
		case "GET":
			msg = upload_get(w, r)
		}
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.CommonBack(msg, w)
}

// routeのベースフォルダ
var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{"/" + apiname, fileupload},
	{"/" + apiname + "/", fileupload},
}

var setupdata UploadPass = UploadPass{}

// Setup() = []webserverv2.WebConfig
//
// セットアップして、HTMLのルートフォルダを用意する
func Setup() ([]webserverv2.WebConfig, error) {
	if err := env.Parse(&setupdata); err != nil {

	}
	if err := os.MkdirAll(setupdata.Pdf, 0777); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(setupdata.Zip, 0777); err != nil {
		return nil, err
	}
	setupdata.flag = true
	return route, nil
}
