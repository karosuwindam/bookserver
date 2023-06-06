package ziptopdf

import (
	"bookserver/config"
	"os"
	"testing"
)

func TestZipToPdf(t *testing.T) {

	t.Setenv("PDF_FILEPASS", "./pdf")
	t.Setenv("ZIP_FILEPASS", "./zip")
	t.Setenv("TMP_FILEPASS", "./tmp")
	t.Setenv("IMG_FILEPASS", "./img")
	cfg, _ := config.EnvRead()
	//
	if err := SetUp(cfg); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	zipfile := "test.zip"
	pdffile := "test.pdf"
	z, err := unzip(zipfile)
	if err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	if _, err := z.imgCopy("test"); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()

	}
	if err := z.img2pdf(pdffile); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	if err := z.removeFolder(); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	if _, err := os.Stat(imgPass + "test.jpg"); err != nil {
		t.Fatalf("%v", err)
	} else {
		os.Remove(imgPass + "test.jpg")
	}
	if _, err := os.Stat(pdfPass + pdffile); err != nil {
		t.Fatalf("%v", err)
	} else {
		os.Remove(imgPass + pdffile)
	}
	if _, err := os.Stat(tmpPass + "test"); err == nil {
		t.Fatalf("not delete folder")
	}
}
