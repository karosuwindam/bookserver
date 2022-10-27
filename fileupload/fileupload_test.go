package fileupload

import (
	"bookserver/config"
	"bookserver/message"
	"bookserver/webserver"
	"bookserver/webserver/common"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestUrlAnalysis(t *testing.T) {
	url := "/v1/test/1"
	out := urlAnalysis(url)
	tmp := []string{"v1", "test", "1"}
	if len(out) != len(tmp) {
		t.Error("Split error")
		t.FailNow()
	}
	for i := 0; i < len(out); i++ {
		if out[i] != tmp[i] {
			t.Errorf("data error %v = %v", out[i], tmp[i])
			t.FailNow()
		}
	}
	t.Log("----------------- urlAnalysis OK --------------------------")
}

func TestUploadSetupt(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)

	t.Log("----------------- Set up --------------------------")
	u, err := Setup()
	if err != nil {
		t.Errorf("setup erro")
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
	if u.Name() != "upload" {
		t.Errorf("input message err %v", u.Name())
	}
	t.Log(u.MessageJsonOut())
	t.Log("----------------- Set up OK --------------------------")

}

func TestUploadServerPost(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)
	t.Log("----------------- upload Server --------------------------")

	u, _ := Setup()

	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	ss.AddV1(common.GUEST, "/upload", u, FIleupload)
	s, _ := ss.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return runServer(ctx, s)
	})
	uploadfile("test.zip", t)
	uploadfile("test.pdf", t)
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	t.Log("----------------- upload Server OK --------------------------")

}

func TestUploadServerList(t *testing.T) {
	pdfpass := "pdf"
	zippass := "zip"
	t.Setenv("PDF_FILEPASS", pdfpass)
	t.Setenv("ZIP_FILEPASS", zippass)
	t.Log("----------------- upload Server --------------------------")

	u, _ := Setup()

	cfg, _ := config.EnvRead()
	ss, _ := webserver.NewSetup(cfg)
	ss.AddV1(common.GUEST, "/upload", u, FIleupload)
	s, _ := ss.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return runServer(ctx, s)
	})
	//処理
	listFileget("http://localhost:8080/v1/upload/pdf", t)
	listFileget("http://localhost:8080/v1/upload/zip", t)
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	t.Log("----------------- upload Server OK --------------------------")

}

func runServer(ctx context.Context, s *webserver.Server) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Srv.Serve(s.L); err != nil &&
			err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	<-ctx.Done()
	message.Println("server stop")
	s.Srv.Shutdown(context.Background())
	return eg.Wait()
}

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

	res, err := http.Post("http://localhost:8080/v1/upload", writer.FormDataContentType(), &buffer)
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
