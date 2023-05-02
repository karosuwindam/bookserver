package upload

import (
	"bookserver/config"
	"bookserver/transform/writetable"
	"bookserver/webserver"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"golang.org/x/sync/errgroup"
)

// 基本セットアップ
func TestUploadSetupt(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)

	t.Log("----------------- Set up --------------------------")
	u := &setupdata
	_, err := Setup()
	if err != nil {
		t.Errorf("setup error")
		t.FailNow()
	}
	if u.Pdf != pdfpass {
		t.Errorf("ENV PDF_FILEPASS %v = %v", u.Pdf, pdfpass)
		t.FailNow()
	}
	if u.Zip != zippass {
		t.Errorf("ENV ZIP_FILEPASS %v = %v", u.Zip, zippass)
		t.FailNow()
	}
	if !u.flag {
		t.Errorf("setup falg %v", u.flag)
		t.FailNow()
	}
	t.Log("----------------- Set up OK --------------------------")

}

// ファイルアップロードテスト
func TestUploadServerPost(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)
	t.Log("----------------- upload Server --------------------------")

	web, _ := Setup()

	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	webserver.Config(ss, web, "")
	s, _ := ss.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return s.Run(ctx)
	})
	uploadfile("test.zip", t)
	uploadfile("test.pdf", t)
	if name, err := GetUploadName(); err != nil {
		t.Fatalf("GtuploadName get Error:%v", err.Error())
		t.FailNow()
	} else {
		fmt.Println(name)
	}
	if name, err := GetUploadName(); err != nil {
		t.Fatalf("GtuploadName get Error:%v", err.Error())
		t.FailNow()
	} else {
		fmt.Println(name)
	}
	if _, err := GetUploadName(); err != nil {
		fmt.Println("GtuploadName get Error:", err.Error())
	} else {
		t.FailNow()
	}
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	os.Remove("pdf/test.pdf")
	t.Log("----------------- upload Server OK --------------------------")

}

// list機能確認
func TestUploadServerList(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)
	t.Log("----------------- upload Server --------------------------")

	web, _ := Setup()

	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	webserver.Config(ss, web, "")
	s, _ := ss.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return s.Run(ctx)
	})
	//処理
	listFileget("http://localhost:8080/upload/pdf", t)
	listFileget("http://localhost:8080/upload/zip", t)
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	t.Log("----------------- upload Server OK --------------------------")

}

func TestUloadServerGet(t *testing.T) {

	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)
	t.Setenv("DBROOTPASS", "./")
	t.Log("----------------- upload Server --------------------------")

	web, _ := Setup()

	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	webserver.Config(ss, web, "")
	writetable.Setup(cfg)
	s, _ := ss.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return s.Run(ctx)
	})
	//処理
	getUrldata("http://localhost:8080/upload/test.pdf", t, http.StatusOK)
	getUrldata("http://localhost:8080/upload/test.zip", t, http.StatusOK)
	getUrldata("http://localhost:8080/upload/bb.zip", t, http.StatusOK)
	getUrldata("http://localhost:8080/upload/bbb", t, http.StatusBadRequest)
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	t.Log("----------------- upload Server OK --------------------------")
}

/*--------------------------------------------------------------------------*/
//テスト用関数

func uploadfile(filename string, t *testing.T) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("Failed to create file writer. %s", err)
	}

	readFile, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open file. %s", err)
	}
	defer readFile.Close()
	io.Copy(fileWriter, readFile)
	writer.Close()

	res, err := http.Post("http://localhost:8080/upload", writer.FormDataContentType(), &buffer)
	if err != nil {
		t.Fatalf("Failed to POST request. %s", err)
	}
	// API レスポンス検証
	message, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read HTTP response body. %s", err)
	}
	t.Logf("message:%v", string(message))

}

func listFileget(url string, t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("LIST", url, nil)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status :%v", resp.StatusCode)
	}
	t.Logf("%s", b)

}

/*---------------確認用関数----------------*/
func getUrldata(url string, t *testing.T, code int) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != code {
		t.Fail()
	}
	fmt.Printf("%s\n", b)

}

func uploadFileget(url string, sendbyte string, t *testing.T) int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(sendbyte)))
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status :%v", resp.StatusCode)
	}
	t.Logf("%s", b)
	return resp.StatusCode

}
