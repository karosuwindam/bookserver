package pdftozip

import (
	"bookserver/config"
	"bookserver/table"
	"bookserver/table/booknames"
	"bookserver/table/filelists"
	"fmt"
	"os"
	"testing"
)

func TestPdfToZip(t *testing.T) {
	os.Setenv("TMP_FILEPASS", "./test/tmp")
	os.Setenv("IMG_FILEPASS", "./test/img")
	os.Setenv("PDF_FILEPASS", "./pdf")
	os.Setenv("ZIP_FILEPASS", "./test/zip")
	os.Setenv("DB_ROOTPASS", "./test/")
	defer func() {
		os.RemoveAll("./test/")
	}()
	config.Init()
	table.Init()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if d, err := ConvertPdfToZipChack("testout.pdf"); err != nil {
		t.Fatal(err)
	} else {
		if d.Name != "testout" || d.Pdfpass != "testout.pdf" || d.Zippass != "testout.zip" {
			t.Fatal(fmt.Sprintf("errdata %v %v %v ", d.Name, d.Pdfpass, d.Zippass))
		}
	}
	if err := ConvertPdfToZip("testout.pdf"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("./test/zip/testout.zip"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("./test/img/testout.jpg"); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.Search("testout"); err != nil {
		t.Fatal(err)
	} else {
		if d[0].Pdfpass != "testout.pdf" || d[0].Zippass != "testout.zip" || d[0].Tag != "" {
			t.Fatal("Error write table data")
		}
	}
}

func TestPdfToZipChangeName(t *testing.T) {
	os.Setenv("TMP_FILEPASS", "./test/tmp")
	os.Setenv("IMG_FILEPASS", "./test/img")
	os.Setenv("PDF_FILEPASS", "./pdf")
	os.Setenv("ZIP_FILEPASS", "./test/zip")
	os.Setenv("DB_ROOTPASS", "./test/")
	defer func() {
		os.RemoveAll("./test/")
	}()
	config.Init()
	table.Init()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	b := booknames.Booknames{
		Name:     "testout",
		Title:    "aa",
		Writer:   "bb",
		Booktype: "ee",
		Burand:   "gg",
		Ext:      "cc,dd",
	}
	b.Add()
	if err := ConvertPdfToZip("testout.pdf"); err != nil {
		t.Fatal(err)
	}
	if d, err := ConvertPdfToZipChack("testout.pdf"); err != nil {
		t.Fatal(err)
	} else {
		if d.Name != "testout" || d.Pdfpass != "testout.pdf" || d.Zippass != "[bb]aa.zip" {
			t.Fatal(fmt.Sprintf("errdata %v %v %v ", d.Name, d.Pdfpass, d.Zippass))
		}
	}
	if _, err := os.Stat("./test/zip/[bb]aa.zip"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("./test/img/testout.jpg"); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.Search("testout"); err != nil {
		t.Fatal(err)
	} else {
		if d[0].Pdfpass != "testout.pdf" || d[0].Zippass != "[bb]aa.zip" || d[0].Tag != "aa,bb,gg,ee,cc,dd" {
			t.Fatal("Error write table data")
		}
	}
	if err := ConvertPdfToZip("testout.pdf"); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.GetAll(); err != nil {
		t.Fatal(err)
	} else if len(d) != 1 {
		t.Fatal("add over table")
	}
}

func TestPdfToZipChangeName_1(t *testing.T) {
	os.Setenv("TMP_FILEPASS", "./test/tmp")
	os.Setenv("IMG_FILEPASS", "./test/img")
	os.Setenv("PDF_FILEPASS", "./pdf")
	os.Setenv("ZIP_FILEPASS", "./test/zip")
	os.Setenv("DB_ROOTPASS", "./test/")
	defer func() {
		os.RemoveAll("./test/")
	}()
	config.Init()
	table.Init()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	b := booknames.Booknames{
		Name:     "testout",
		Title:    "aa",
		Writer:   "bb",
		Booktype: "ee",
		Burand:   "gg",
		Ext:      "cc,dd",
	}
	b.Add()
	if d, err := ConvertPdfToZipChack("testout1.pdf"); err != nil {
		t.Fatal(err)
	} else {
		if d.Name != "testout1" || d.Pdfpass != "testout1.pdf" || d.Zippass != "[bb]aa01.zip" {
			t.Fatal(fmt.Sprintf("errdata %v %v %v ", d.Name, d.Pdfpass, d.Zippass))
		}
	}
	if err := ConvertPdfToZip("testout1.pdf"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("./test/zip/[bb]aa01.zip"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("./test/img/testout1.jpg"); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.Search("testout"); err != nil {
		t.Fatal(err)
	} else {
		if d[0].Pdfpass != "testout1.pdf" || d[0].Zippass != "[bb]aa01.zip" || d[0].Tag != "aa,bb,gg,ee,cc,dd" {
			t.Fatal("Error write table data")
		}
	}
	//異なる名前のファイルが登録されたとき
	if err := ConvertPdfToZip("testout.pdf"); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.GetAll(); err != nil {
		t.Fatal(err)
	} else if len(d) != 2 {
		t.Fatal("add over table")
	}
}
