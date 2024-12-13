package booknames

import (
	"bookserver/config"
	"context"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBookNameTable(t *testing.T) {
	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	db, _ := openSqlite3()
	if err := Init(db); err != nil {
		t.Fatal(err)
	}
	tmp := Booknames{
		Id:   100,
		Name: "test1",
	}
	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	tmp.Name = "test2"
	if err := tmp.Add(); err != nil {
		t.Fatal(err)
	}
	if d, err := GetAll(context.TODO()); err != nil {
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
	}
	if _, err := GetId(0); err == nil {
		t.Fatal("Fail Read error")
	}
	if d, err := GetId(1); err != nil {
		t.Fatal(err)
	} else {
		d.Title = "test2"
		if err = d.Update(); err != nil {
			t.Fatal(err)
		}
		if d, err = GetId(1); err != nil {
			t.Fatal(err)
		} else {
			if d.Title != "test2" {
				t.Fatal("Title not input data for test2")
			}
		}
	}
	if d, err := GetName("test2"); err != nil {
		t.Fatal(err)
	} else {
		if d.Id != 2 {
			t.Fatal("Fail Not input data", d)
		}
	}
	if err := Delete(2); err != nil {
		t.Fatal(err)
	}
	if _, err := GetId(2); err == nil {
		t.Fatal("Fail Read error")
	}
	tmp.Name = "test5"
	tmp.Add()
	if d, err := Search("t5"); err != nil {
		t.Fatal(err)
	} else if len(d) != 1 {
		t.Fatal("Search Error Output by 1")
	} else {
		if d[0].Name != "test5" {
			t.Fatal("input serch Name error")
		}
	}
	tmp.Add()
	tmp.Add()
	tmp.Add()
	if b, err := ReadRandData(2); err != nil {
		t.Fatal(err)
	} else if len(b) != 2 {
		t.Fatal("err rand data")
	}
	if b, err := ReadRandData(10); err != nil {
		t.Fatal(err)
	} else if len(b) != 5 {
		t.Fatal("err rand data")
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
