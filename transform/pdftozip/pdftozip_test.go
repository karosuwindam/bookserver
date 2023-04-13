package pdftozip

import (
	"bookserver/config"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Setenv("PDF_FILEPASS", "./pdf")
	t.Setenv("ZIP_FILEPASS", "./zip")
	t.Setenv("TMP_FILEPASS", "./tmp")
	cfg, _ := config.EnvRead()
	t.Log("------------- Set up ----------------------")
	if err := SetUp(cfg); err != nil {
		t.Errorf("setup error")
		t.FailNow()
	}
	if tmpPass != "./tmp/" {
		t.Errorf("ENV TMP_FILEPASS %v = %v", tmpPass, "./tmp/")
		t.FailNow()
	}
	if zipPass != "./zip/" {
		t.Errorf("ENV ZIP_FILEPASS %v = %v", zipPass, "./zip/")
		t.FailNow()
	}
	if pdfPass != "./pdf/" {
		t.Errorf("ENV PDF_FILEPASS %v = %v", pdfPass, "./pdf/")
		t.FailNow()
	}

}

func TestPdftoimages(t *testing.T) {
	t.Setenv("PDF_FILEPASS", "./pdf")
	t.Setenv("ZIP_FILEPASS", "./zip")
	t.Setenv("TMP_FILEPASS", "./tmp")
	cfg, _ := config.EnvRead()
	if err := SetUp(cfg); err != nil {
		t.Errorf("setup error")
		t.FailNow()
	}
	t.Log("------------- Pdf to images ----------------------")
	if err := Pdftoimage("testout.pdf", "test.zip"); err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	t.Log("------------- Pdf to images End ----------------------")

}
