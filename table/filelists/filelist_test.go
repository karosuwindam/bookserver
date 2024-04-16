package filelists_test

import (
	"bookserver/config"
	"bookserver/table/filelists"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFilelists(t *testing.T) {
	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	db, _ := openSqlite3()
	if err := filelists.Init(db); err != nil {
		t.Fatal(err)
	}
	tmp := filelists.Filelists{
		Id:   100,
		Name: "test1",
		Tag:  "test1",
	}
	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	tmp.Name = "test2"
	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	if d, err := filelists.GetAll(); err != nil {
		t.Fatal(err)
	} else {
		if len(d) != 2 {
			t.Fatal("Fail Add for table")
		}
		if d[0].Id != 1 || d[0].Name != "test1" {
			t.Fatal("Fail Not input data", d[0])
		}
		if d[1].Id != 2 || d[1].Name != "test2" {
			t.Fatal("Fail Not input data", d[1])
		}
		if d[0].Tag != d[1].Tag {
			t.Fatal("Fail Change tag data")
		}
		if _, err := filelists.GetId(0); err == nil {
			t.Fatal("Fail Read error")
		}
		if d, err := filelists.GetId(1); err != nil {
			t.Fatal(err)
		} else {
			d.Pdfpass = "test3"
			if err = d.Update(); err != nil {
				t.Fatal(err)
			}
			if d, err = filelists.GetId(1); err != nil {
				t.Fatal(err)
			} else {
				if d.Pdfpass != "test3" {
					t.Fatal("Pdfpass not input data for test3")
				}
			}
		}
	}
	if err := filelists.Delete(2); err != nil {
		t.Fatal(err)
	}
	if _, err := filelists.GetId(2); err == nil {
		t.Fatal("Fail Read error")
	}
	tmp.Name = "test5"
	tmp.Add()
	if d, err := filelists.Search("t5"); err != nil {
		t.Fatal(err)
	} else if len(d) != 1 {
		t.Fatal("Search Error Output by 1")
	} else {
		if d[0].Name != "test5" {
			t.Fatal("input serch Name error")
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
