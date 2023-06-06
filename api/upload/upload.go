package upload

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/dirread"
	"bookserver/table"
	"bookserver/transform/writetable"
	"bookserver/webserver"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var apiname string = "upload" //api名
var maxMultiMemory int64

const (
	MIN_MULTI_MEMORY = 32 << 20  // 32MB
	MAX_MULTI_MEMORY = 256 << 20 // 256MB
)

type UploadPass struct {
	Pdf  string `env:"PDF_FILEPASS" envDefault:"./upload/pdf"` //PDFのアップロード先フォルダ
	Zip  string `env:"ZIP_FILEPASS" envDefault:"./upload/zip"` //ZIPのアップロード先フォルダ
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
	Overwrite bool `json:Overwrite`
	// データあり
	Register bool `json:"Register"`
	// 書き換え名前
	Name       string              `json:"Name"`
	ChangeName writetable.PdftoZip `json:"ChangeName"`
}

// upload_filses
//
// 複数のファイルアップロード
func upload_files(w http.ResponseWriter, r *http.Request) common.Result {
	var wg sync.WaitGroup
	var mux sync.Mutex
	msg := common.Result{Code: http.StatusOK, Date: time.Now()}
	// 複数のファイルを取得する
	r.ParseMultipartForm(maxMultiMemory)
	files := r.MultipartForm.File["file"]
	for _, file := range files {
		wg.Add(1)
		fmt.Println("receive file for", file.Filename, ",size:", file.Size)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			filename := file.Filename
			// ファイルをオープンする
			f, err := file.Open()
			if err != nil {
				mux.Lock()
				msg.Code = http.StatusAccepted
				msg.Result = err.Error()
				mux.Unlock()
				log.Println(err)
			}
			defer f.Close()
			// ファイルを保存する
			savepass := ""
			if strings.Index(strings.ToLower(filename), "pdf") > 0 {
				savepass = setupdata.Pdf + "/"
			} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
				savepass = setupdata.Zip + "/"
			} else {
				log.Println("err file:", filename)
			}
			out, err := os.Create(savepass + "/" + filename)
			if err != nil {
				mux.Lock()
				msg.Code = http.StatusAccepted
				msg.Result = err.Error()
				mux.Unlock()
				log.Println(err)

			}
			defer out.Close()
			// ファイルをコピーする
			if _, err = f.Seek(0, 0); err != nil {
				mux.Lock()
				msg.Code = http.StatusAccepted
				msg.Result = err.Error()
				mux.Unlock()
				log.Println(err)
				defer os.Remove(savepass + "/" + filename)

			}
			if _, err = io.Copy(out, f); err != nil {
				mux.Lock()
				msg.Code = http.StatusAccepted
				msg.Result = err.Error()
				mux.Unlock()
				log.Println(err)
				defer os.Remove(savepass + "/" + filename)
			}
			mux.Lock()
			uploadname <- filename
			mux.Unlock()
			fmt.Println("Create File End for:", filename)
		}(file)
	}
	wg.Wait()
	if msg.Code != http.StatusAccepted {
		msg.Result = "OK"
	}
	return msg
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
	uploadname <- filename
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
	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""
	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+1 {
			tName = sUrl[i+1]
			break
		}
	}
	//jsonデータの取り出し
	var out UploadFileGet
	out.Name = tName
	if strings.Index(strings.ToLower(out.Name), "pdf") > 0 {
		// savepass = setupdata.Pdf + "/"
		out.ChangeName, _ = writetable.CreatePdfToZip(out.Name)
		out.Register = common.Exists(setupdata.Pdf + "/" + out.Name)
	} else if strings.Index(strings.ToLower(out.Name), "zip") > 0 {
		out.Register = common.Exists(setupdata.Zip + "/" + out.Name)
		// savepass = setupdata.Zip + "/"
	} else {
		msg.Result = errors.New("input pass filename:" + out.Name).Error()
		msg.Code = http.StatusBadRequest
		return msg
	}
	//sqlテーブルの検索
	if jout, err := sql.Search(table.FILELIST, tName); err == nil && jout != "[]" {
		out.Overwrite = true
	}
	msg.Result = out
	return msg
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
			msg = upload_files(w, r)
		case "LIST":
			msg = upload_list(w, r)
		case "GET":
			msg = upload_get(w, r)
		case "PUT":
			msg = upload_get(w, r)
		}
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.CommonBack(msg, w)
}

// routeのベースフォルダ
var route []webserver.WebConfig = []webserver.WebConfig{
	{"/" + apiname, fileupload},
	{"/" + apiname + "/", fileupload},
}

// GetuploadName = string,error
//
// Uploadされたファイル名を取得
func GetUploadName() (string, error) {
	output := ""
	var err error = nil
	select {
	case output = <-uploadname:
	case <-time.After(100 * time.Millisecond):
		err = errors.New("time out")
	}
	return output, err
}

var setupdata UploadPass = UploadPass{}
var uploadname chan string = make(chan string, 30)

// Setup() = []webserver.WebConfig
//
// セットアップして、HTMLのルートフォルダを用意する
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	uploadname = make(chan string, 30)
	// if err := env.Parse(&setupdata); err != nil {

	// }
	setupdata.Pdf = cfg.Folder.Pdf
	setupdata.Zip = cfg.Folder.Zip
	if err := os.MkdirAll(setupdata.Pdf, 0777); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(setupdata.Zip, 0777); err != nil {
		return nil, err
	}
	//マルチアップロードの確保メモリ最大値入力
	if tmp, err := setupMaxMultiMemory(cfg.Upload.MAX_MULTI_MEMORY); err != nil || tmp < MIN_MULTI_MEMORY {
		if err != nil {
			fmt.Println(err)
		}
		maxMultiMemory = MIN_MULTI_MEMORY
	} else {
		maxMultiMemory = tmp
	}
	if sqlcfg, err := sqlSetup(cfg); err == nil {
		sql = sqlcfg
	} else {
		return nil, err
	}
	setupdata.flag = true
	return route, nil
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
