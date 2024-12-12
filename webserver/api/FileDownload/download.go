package filedownload

import (
	"bookserver/table/filelists"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func GetDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	filetype := r.PathValue("filetype")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetDownload strconv.Atoi id=%v", tmpid),
			"Id", tmpid,
			"Error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	filename, err := readFileName(id, filetype)
	if err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetDownload readFileName id=%v filetype=%v", id, filetype),
			"Id", id,
			"Filetype", filetype,
			"Error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if b, err := readFileData(filename, filetype); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetDownload readFileData filename=%v filetype=%v", filename, filetype),
			"Filename", filename,
			"Filetype", filetype,
			"Error", err,
		)
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
	ctx := context.TODO()
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
				slog.ErrorContext(ctx,
					fmt.Sprintf("readFileData file.Read pass=%v", pass),
					"Pass", pass,
					"Error", err,
				)
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
