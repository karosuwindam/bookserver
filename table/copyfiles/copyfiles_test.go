package copyfiles_test

import (
	"bookserver/config"
	"bookserver/table/copyfiles"
	"context"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCopyfiles(t *testing.T) {

	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	db, _ := openSqlite3()
	if err := copyfiles.Init(db); err != nil {
		t.Fatal(err)
	}
	tmp := copyfiles.Copyfile{
		Id:       100,
		Zippass:  "test1",
		Filesize: 200,
	}

	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	tmp.Zippass = "test2"
	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	if d, err := copyfiles.GetAll(context.TODO()); err != nil {
		t.Fatal(err)
	} else {
		if len(d) != 2 {
			t.Fatal("Fail Add for table")
		}
		if d[0].Id != 1 || d[0].Zippass != "test1" {
			t.Fatal("Fail Not input data", d[0])
		}
		if d[1].Id != 2 || d[1].Zippass != "test2" {
			t.Fatal("Fail Not input data", d[1])
		}
		if d[0].Filesize != d[1].Filesize && d[0].Filesize != 200 {
			t.Fatal("Fail Change File size data")
		}
		if d[0].Copyflag != d[1].Copyflag && d[0].Copyflag != 0 {
			t.Fatal("Fail Copy flag")
		}
		if _, err := copyfiles.GetId(0); err == nil {
			t.Fatal("Fail Read error")
		}
		if d, err := copyfiles.GetId(1); err != nil {
			t.Fatal(err)
		} else {
			d.Zippass = "test3"
			if err = d.Update(); err != nil {
				t.Fatal(err)
			}
			if err = d.ON(); err != nil {
				t.Fatal(err)
			}
			if d, err = copyfiles.GetId(1); err != nil {
				t.Fatal(err)
			} else {
				if d.Zippass != "test3" {
					t.Fatal("Pdfpass not input data for test3")
				}
				if d.Copyflag != 1 {
					t.Fatal("CopyFlag change flag")
				}
			}
		}
	}
	if d, err := copyfiles.OnOFFSearch(copyfiles.ON); err != nil {
		t.Fatal(err)
	} else {
		if d[0].Zippass != "test3" {
			t.Fatal("Search error data by ON")
		}
	}
	if d, err := copyfiles.OnOFFSearch(copyfiles.OFF); err != nil {
		t.Fatal(err)
	} else {
		if d[0].Zippass != "test2" {
			t.Fatal("Search error data by OFF")
		}
	}
	if err := copyfiles.Delete(2); err != nil {
		t.Fatal(err)
	}
	if _, err := copyfiles.GetId(2); err == nil {
		t.Fatal("Fail Read error")
	}
	tmp.Zippass = "test5"
	tmp.Add()
	if d, err := copyfiles.Search("t5"); err != nil {
		t.Fatal(err)
	} else if len(d) != 1 {
		t.Fatal("Search Error Output by 1")
	} else {
		if d[0].Zippass != "test5" {
			t.Fatal("input serch Name error")
		}
	}
	tmp.Zippass = "ああああ"
	tmp.Add()
	if d, err := copyfiles.GetZipName("ああああ"); err != nil {
		t.Fatal(err)
	} else {
		if d.Zippass != "ああああ" {
			t.Fatal("Fail Read error")
		}
	}
	if _, err := copyfiles.GetZipName("いいいい"); err == nil {
		t.Fatal("not add")
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
