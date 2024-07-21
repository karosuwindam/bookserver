package fileupload

import (
	"bookserver/config"
	"bookserver/controller/convert/pdftozip"
	"bookserver/table/uploadtmp"
	"bookserver/webserver/api/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ファイルのアップロード処理
func PostUploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := config.TracerS(ctx, "PostUploadFile", r.URL.Path)
	defer span.End()

	var wg sync.WaitGroup
	//複数のファイルを取得する
	if err := r.ParseMultipartForm(maxMultiMemory); err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Println("error:", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	files := r.MultipartForm.File["file"]
	list := []string{}
	for _, file := range files {
		list = append(list, fmt.Sprintf("file=%v size=%v", file.Filename, file.Size))
	}
	span.SetAttributes(attribute.StringSlice("list", list))
	log.Println("info:", r.URL, r.Method, "File=", list)
	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			_, span2 := config.TracerS(ctx, "saveFileData "+file.Filename, "saveFileData")
			defer span2.End()
			defer wg.Done()
			if err := saveFileData(file); err != nil {
				span2.SetStatus(codes.Error, err.Error())
				log.Println("Not Save Error :", file.Filename)
			}
			span2.SetAttributes(attribute.String("filename", file.Filename))
			span2.SetAttributes(attribute.Int64("fileSize", file.Size))
		}(file)
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)

}

// アップロードしたときのファイル保存処理
func saveFileData(file *multipart.FileHeader) error {
	filename := file.Filename
	//テーブルへの追加
	sqlTmp := uploadtmp.UploadTmp{
		Name: filename,
	}
	if err := sqlTmp.CheckName(); err != nil {
		return err
	}
	if sqlTmp.Id > 0 {
		if err := sqlTmp.CountClear(); err != nil {
			return err
		}
	} else {
		if err := sqlTmp.Add(); err != nil {
			return err
		}
	}

	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		//pdfの処理
		//ファイルを特定フォルダへ保存
		if err := savePdfFile(file); err != nil {
			return err
		}
		//テーブルを書き換える
		if err := sqlTmp.SetPdfPath(pdfFolder + filename); err != nil {
			return err
		}
	} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
		//zipの処理
		//ファイルを特定フォルダへ保存
		if err := saveZipFile(file); err != nil {
			return err
		}
		//テーブルを書き換える
		if err := sqlTmp.SetZipPath(zipFolder + filename); err != nil {
			return err
		}
	} else {
		//定期処理を無効にする処理
		if err := sqlTmp.FlagOn(); err != nil {
			return err
		}
	}
	return nil
}

// PDFファイルの保存を実施する。
func savePdfFile(file *multipart.FileHeader) error {
	savePath := pdfFolder + file.Filename
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = f.Seek(0, 0); err != nil {
		defer os.Remove(savePath)
		return err
	}
	if _, err = io.Copy(out, f); err != nil {
		defer os.Remove(savePath)
		return err
	}
	log.Println("info:", "Create File for:", savePath)
	return nil
}

// Zipファイルの保存を実施する。
func saveZipFile(file *multipart.FileHeader) error {
	savePath := zipFolder + file.Filename
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = f.Seek(0, 0); err != nil {
		defer os.Remove(savePath)
		return err
	}
	if _, err = io.Copy(out, f); err != nil {
		defer os.Remove(savePath)
		return err
	}
	log.Println("info:", "Create File for:", savePath)
	return nil
}

type ChangeDataS struct {
	// // 上書き
	Overwrite bool `json:Overwrite`
	// データあり
	Register bool `json:"Register"`
	// 書き換え名前
	Name        string      `json:"Name"`
	ChanageName ChangeNameS `json:"ChangeName"`
}

type ChangeNameS struct {
	Name       string `json:"Name"` //登録用の名称(巻数情報も含む)
	InputFile  string `json:"Pdf"`  //入力ファイル(pdf)
	OutputFile string `json:"Zip"`  //出力ファイル(zip)
	Tag        string `json:"Tag"`  //検索用のタグ情報
}

// アップロード後の変換文字列を確認
// その結果filelistに登録予定のファイルを作成
func GetUploadFileChangeData(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	output := ChangeDataS{}
	filename := r.PathValue("filename")
	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		b, err := pdftozip.ConvertPdfToZipChack(filename)
		if err != nil {

		} else {
			output.Name = filename
			output.ChanageName = ChangeNameS{
				Name:       b.Name,
				InputFile:  b.Pdfpass,
				OutputFile: b.Zippass,
				Tag:        b.Tag,
			}

			if ok, err := checkByFIle("zip", b.Zippass); err != nil {
				log.Println("debug:", r.Method, r.URL, "message:", err)
				output.Register = false
			} else {
				output.Register = ok
			}

			if ok, err := checkByFIle("pdf", b.Pdfpass); err != nil {
				log.Println("debug:", r.Method, r.URL, "message:", err)
				output.Overwrite = false
			} else {
				output.Overwrite = ok
			}
		}
	} else if strings.Index(strings.ToLower(filename), "zip") > 0 {
		output.Name = filename
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url status error"))
		return
	}
	msg := common.Message(output)
	if b, errj := json.Marshal(&msg); errj != nil {
		log.Println(errj)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url status error"))
	} else {
		w.Write(b)
	}
}

// アップロードフォルダ内に既存ファイルの存在を確認
func GetUplodFileCheck(w http.ResponseWriter, r *http.Request) {
	filetype := r.PathValue("filetype")
	filename := r.PathValue("filename")
	log.Println("info:", r.URL, r.Method)
	var output string
	if ok, err := checkByFIle(filetype, filename); err != nil {
		log.Println("debug:", r.Method, r.URL, "message:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url status error"))
		return
	} else if ok {
		output = "{\"message\":\"ok\"}"
	} else {
		output = "{\"message\":\"ng\"}"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

func checkByFIle(filetype, filename string) (bool, error) {
	switch filetype {
	case "zip":
		return zipCheckFilePass(filename)
	case "pdf":
		return pdfCheckFilePass(filename)
	}
	return false, errors.New("Error FIle Type")
}

func pdfCheckFilePass(filename string) (bool, error) {
	checkpass := pdfFolder
	if strings.Index(strings.ToLower(filename), "pdf") > 0 {
		checkpass += filename
	} else {
		return false, errors.New("Error File type")
	}
	//ファイル確認
	_, err := os.Stat(checkpass)
	return err == nil, nil
}

func zipCheckFilePass(filename string) (bool, error) {
	checkpass := zipFolder
	if strings.Index(strings.ToLower(filename), "zip") > 0 {
		checkpass += filename
	} else {
		return false, errors.New("Error File type")
	}
	//ファイル確認
	_, err := os.Stat(checkpass)
	return err == nil, nil
}
