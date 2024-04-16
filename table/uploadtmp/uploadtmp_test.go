package uploadtmp_test

import (
	"bookserver/config"
	"bookserver/table/uploadtmp"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUploadTmpTable(t *testing.T) {
	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	db, _ := openSqlite3()
	if err := uploadtmp.Init(db); err != nil {
		t.Fatal(err)
	}
	idata := uploadtmp.UploadTmp{
		Id:   100,
		Name: "test1",
	}
	if err := idata.Add(); err != nil {
		t.Fatal(err)
	}
	if idata.Id != 1 {
		t.Fatal("Error Read Id")
	}
	if tmp, err := uploadtmp.GetAllbyFlag(false); err != nil {
		t.Fatal(err)
	} else {
		if len(tmp) != 1 {
			t.Fatal("read tmp error")
		}
		if tmp[0].Name != "test1" {
			t.Fatal("Not input data Name test1")
		}
		if err = tmp[0].SetPdfPath("test2"); err != nil {
			t.Fatal(err)
		}
		if err = tmp[0].SetZipPath("test3"); err != nil {
			t.Fatal(err)
		}
		if tmp[0].SavePdf != "test2" || tmp[0].SaveZip != "test3" {
			t.Fatal("Error input SaveName")
		}
		if err = tmp[0].CountUp(); err != nil {
			t.Fatal(err)
		}
		if err = tmp[0].FlagOn(); err != nil {
			t.Fatal(err)
		}
	}
	if tmp, err := uploadtmp.GetAllbyFlag(true); err != nil {
		t.Fatal(err)
	} else {
		if len(tmp) != 1 {
			t.Fatal("read tmp error")
		}
		if tmp[0].Name != "test1" || tmp[0].SavePdf != "test2" || tmp[0].SaveZip != "test3" || tmp[0].Count != 1 {
			t.Fatal("Error input SaveName")
		}
	}

}

func openSqlite3() (*gorm.DB, error) {
	filepass := config.DB.DBROOTPASS
	if filepass[len(filepass)-1:] != "/" {
		filepass += "/"
	}
	if f, err := os.Stat(filepass); os.IsNotExist(err) || !f.IsDir() {
		os.MkdirAll(filepass, 0766)
	}
	filepass += config.DB.DBFILE
	return gorm.Open(sqlite.Open(filepass), &gorm.Config{})

}
