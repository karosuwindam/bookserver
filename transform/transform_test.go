package transform

import (
	"bookserver/config"
	"context"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Setenv("PDF_FILEPASS", "./pdftozip/pdf")
	t.Setenv("ZIP_FILEPASS", "./pdftozip/zip")
	t.Setenv("TMP_FILEPASS", "./pdftozip/tmp")
	cfg, _ := config.EnvRead()
	if err := Setup(cfg); err != nil {
		t.FailNow()
	}
	ctx, cancel := context.WithCancel(context.Background())
	go Run(ctx)
	samplefile := PdftoZip{
		InputFile:  "testout.pdf",
		OutputFile: "test.zip",
	}
	if err := Add(samplefile); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
	if err := Add("b"); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
	time.Sleep(time.Second)
	cancel()
	if err := Wait(); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
}
