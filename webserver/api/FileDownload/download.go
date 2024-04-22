package filedownload

import (
	"bookserver/table/filelists"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func GetDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	filetype := r.PathValue("filetype")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	filename, err := readFileName(id, filetype)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if b, err := readFileData(filename, filetype); err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/"+filetype)
		// ファイルの長さ
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		// bodyに書き込み
		w.Write(b)

	}

}

func readFileName(id int, filetype string) (string, error) {
	var out string
	d, err := filelists.GetId(id)
	if err != nil {
		return out, errors.Wrap(err, fmt.Sprintf("filelists.GetId(%v)", id))
	}
	switch strings.ToLower(filetype) {
	case "zip":
		out = d.Zippass
	case "pdf":
		out = d.Pdfpass
	default:
		return out, errors.New(fmt.Sprintf("eroor file type:%v", filetype))
	}
	return out, nil
}

func readFileData(filename, filetype string) ([]byte, error) {
	var buffer []byte
	pass, err := createPass(filename, filetype)
	if err != nil {
		return buffer, errors.Wrap(err, fmt.Sprintf("createPass(%v,%v)", filename, filetype))
	}
	if file, err := os.Open(pass); err != nil {
		return buffer, errors.Wrap(err, fmt.Sprintf("os.Open(%v)", pass))
	} else {
		defer file.Close()
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				log.Println("errors:", err)
				break
			}
			buffer = append(buffer, buf[:n]...)
		}
	}

	return buffer, nil
}

func createPass(name, filetype string) (string, error) {
	var out string
	switch strings.ToLower(filetype) {
	case "zip":
		out = zipFolder + name
	case "pdf":
		out = pdfFolder + name
	default:
		return out, errors.New(fmt.Sprintf("eroor file type:%v", filetype))

	}
	return out, nil
}
