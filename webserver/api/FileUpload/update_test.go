package fileupload_test

import (
	"bookserver/config"
	"bookserver/table"
	"bookserver/table/uploadtmp"
	fileupload "bookserver/webserver/api/FileUpload"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestGetUploadTest(t *testing.T) {

	os.Setenv("DB_ROOTPASS", "./db/")
	os.Setenv("PDF_FILEPASS", "./pdf")
	os.Setenv("ZIP_FILEPASS", "./zip")
	config.Init()
	table.Init()
	defer func() {
		os.RemoveAll("./db/")
	}()

	mux := http.NewServeMux()
	fileupload.Init("/upload", mux)

	srv := &http.Server{
		Addr:    config.Web.Hostname + ":" + config.Web.Port,
		Handler: mux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer srv.Shutdown(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	if code, b := getUpload("/upload/pdf/test.pdf"); code != 200 {
		t.Fatal("not url")
	} else {
		var tmp message
		if err := json.Unmarshal([]byte(b), &tmp); err != nil {
			t.Fatal(err)
		} else if tmp.Message != "ok" {
			t.Fatal("Meesage is ok != " + tmp.Message)
		}
	}

	if code, b := getUpload("/upload/zip/test.zip"); code != 200 {
		t.Fatal("not url")
	} else {
		var tmp message
		if err := json.Unmarshal([]byte(b), &tmp); err != nil {
			t.Fatal(err)
		} else if tmp.Message != "ok" {
			t.Fatal("Meesage is ok != " + tmp.Message)
		}
	}
	if code, _ := getUpload("/upload/zip/test.pdf"); code != 400 {
		t.Fatal("not url")
	}
	if code, b := getUpload("/upload/zip/test1.zip"); code != 200 {
		t.Fatal("not url")
	} else {
		var tmp message
		if err := json.Unmarshal([]byte(b), &tmp); err != nil {
			t.Fatal(err)
		} else if tmp.Message != "ng" {
			t.Fatal("Meesage is ng != " + tmp.Message)
		}
	}
	cancel()
}

func TestPostUploadTest(t *testing.T) {

	os.Setenv("DB_ROOTPASS", "./db/")
	os.Setenv("PDF_FILEPASS", "./tmp/pdf")
	os.Setenv("ZIP_FILEPASS", "./tmp/zip")
	config.Init()
	table.Init()
	defer func() {
		os.RemoveAll("./db/")
		os.RemoveAll("./tmp/")
	}()
	mux := http.NewServeMux()
	fileupload.Init("/upload", mux)

	srv := &http.Server{
		Addr:    config.Web.Hostname + ":" + config.Web.Port,
		Handler: mux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer srv.Shutdown(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()
	if code, err := postUpload("test.pdf"); err != nil {
		t.Fatal(err)
	} else if code != 200 {
		t.Fatal("upload file not now")
	}
	if code, err := postUpload("./test.zip"); err != nil {
		t.Fatal(err)
	} else if code != 200 {
		t.Fatal("upload file not now")
	}
	if code, err := postUpload("./test.txt"); err != nil {
		t.Fatal(err)
	} else if code != 200 {
		t.Fatal("upload file not now")
	}

	if bs, err := uploadtmp.GetAllbyFlag(false); err != nil {
		t.Fatal(err)
	} else {
		if len(bs) != 2 {
			t.Fatal("error upload data")
		}
		if bs[0].Name != "test.pdf" {
			t.Fatal("Not Table Write test.pdf")
		} else {
			if bs[0].SaveZip != "" || bs[0].SavePdf != "./tmp/pdf/test.pdf" {
				t.Fatal("Not input Table data")
			}
		}
		if bs[1].Name != "test.zip" {
			t.Fatal("Not Table Write test.zip")
		} else {
			if bs[1].SavePdf != "" || bs[1].SaveZip != "./tmp/zip/test.zip" {
				t.Fatal("Not input Table data")
			}
		}
	}
	if bs, err := uploadtmp.GetAllbyFlag(true); err != nil {
		t.Fatal(err)
	} else {
		if len(bs) != 1 {
			t.Fatal("error upload data")
		}
		if bs[0].Name != "test.txt" {
			t.Fatal("Not Table Write test.txt")
		}
	}
	cancel()
}

type message struct {
	Message string `json:message`
}

func getUpload(url string) (int, string) {
	res, err := http.Get("http://localhost:8080/" + url)
	if err != nil {
		return 404, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, ""
	}
	return res.StatusCode, string(body)

}

func postUpload(filename string) (int, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return 0, err
	}
	readFile, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	io.Copy(fileWriter, readFile)
	readFile.Close()
	writer.Close()

	res, err := http.Post("http://localhost:8080/upload", writer.FormDataContentType(), &buffer)
	if err != nil {
		return res.StatusCode, err
	}
	defer res.Body.Close()
	m, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}
	log.Println("info:", string(m))

	return res.StatusCode, nil

}
