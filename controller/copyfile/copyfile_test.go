package copyfile

import (
	"bookserver/config"
	"bookserver/table"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"context"
	"fmt"
	"os"
	"testing"
)

func TestCopyfile(t *testing.T) {
	os.Setenv("ZIP_FILEPASS", "./zip")
	config.Init()
	Init()
	defer func() {
		os.RemoveAll(publicpass)
	}()
	if i := checkZipFileSize("test.zip"); i != 23444 {
		t.Fatal(fmt.Sprintln("file size 23444 !=", i))
	}
	if i := checkZipFileSize("test1.zip"); i != 0 {
		t.Fatal(fmt.Sprintln("file size 0 !=", i))
	}
	if err := copyFileForZipToPublic("test.zip"); err != nil {
		t.Fatal(err)
	}
	if err := copyFileForZipToPublic("test.zip"); err != nil {
		t.Fatal(err)
	}
	if err := copyFileForZipToPublic("test1.zip"); err == nil {
		t.Fatal(fmt.Sprintln("error"))
	}
	if err := removeFileFromPublic("test.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestCopyFileTable(t *testing.T) {
	os.Setenv("ZIP_FILEPASS", "./zip")
	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	table.Init()
	Init()
	defer func() {
		os.RemoveAll(publicpass)
		os.RemoveAll("./db/")
	}()
	filelist := filelists.Filelists{
		Name:    "test",
		Pdfpass: "test.pdf",
		Zippass: "test.zip",
		Tag:     "",
	}
	filelist.Add()
	filelist = filelists.Filelists{
		Name:    "test1",
		Pdfpass: "test1.pdf",
		Zippass: "test1.zip",
		Tag:     "",
	}
	filelist.Add()
	cp := CopyFIleData{
		id:   1,
		flag: true,
	}
	if err := cp.AddTable(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if err := ChackCopyFileTableDataAll(); err != nil {
		t.Fatal(err)
	}
	if d, err := copyfiles.OnOFFSearch(copyfiles.ON); err != nil {
		t.Fatal(err)
	} else {
		if len(d) != 1 {
			t.Fatal(fmt.Sprintf("%v", len(d)))
		}
	}
	removeFileFromPublic("test.zip")
	if err := ChackCopyFileTableDataAll(); err != nil {
		t.Fatal(err)
	}
	if d, err := copyfiles.OnOFFSearch(copyfiles.OFF); err != nil {
		t.Fatal(err)
	} else {
		if len(d) != 1 {
			t.Fatal(fmt.Sprintf("%v", len(d)))
		}
	}
}
