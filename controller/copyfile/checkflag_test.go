package copyfile_test

import (
	"bookserver/config"
	"bookserver/controller/copyfile"
	"bookserver/table"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCheckFlag(t *testing.T) {

	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	table.Init()
	filelist := filelists.Filelists{
		Name:    "test1",
		Pdfpass: "test1.pdf",
		Zippass: "test1.zip",
		Tag:     "",
	}
	filelist.Add()
	filelist = filelists.Filelists{
		Name:    "test2",
		Pdfpass: "test2.pdf",
		Zippass: "test2.zip",
		Tag:     "",
	}
	filelist.Add()
	copyfiledata := copyfiles.Copyfile{
		Zippass:  "test1.zip",
		Filesize: 123,
		Copyflag: 1,
	}
	copyfiledata.Add()

	if d, err := copyfile.ReadCopyFIleFlagById(1); err != nil {
		t.Fatal(err)
	} else {
		if d.Zippass != copyfiledata.Zippass || d.Filesize != copyfiledata.Filesize || d.Copyflag != copyfiledata.Copyflag {
			t.Fatal(fmt.Sprintln(d, copyfiledata))
		}
		d, _ := json.Marshal(&d)
		log.Println(string(d))
	}
	if d, err := copyfile.ReadCopyFIleFlagById(2); err != nil {
		t.Fatal(err)
	} else {
		if d.Zippass != filelist.Zippass || d.Copyflag != 0 {
			t.Fatal(fmt.Sprintln(d))
		}
		d, _ := json.Marshal(&d)
		log.Println(string(d))
	}
	if _, err := copyfile.ReadCopyFIleFlagById(3); err == nil {
		t.Fatal("Not data error")
	}
}
