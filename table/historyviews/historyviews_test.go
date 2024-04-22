package historyviews_test

import (
	"bookserver/config"
	"bookserver/table/filelists"
	"bookserver/table/historyviews"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHistoryView(t *testing.T) {

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
		Name: "",
	}
	for i := 1; i < 100; i++ {
		tmp.Name = fmt.Sprintf("test%v", i)
		if err := tmp.Add(); err != nil {
			t.Fatal(err)
		}
	}
	historyviews.Init(db)
	tmp3 := []int{}
	tmp2 := historyviews.HistoryViews{}
	for i := 0; i < 100; i++ {
		tmp2.FileId = rand.Intn(99) + 1
		if err := tmp2.Add(); err != nil {
			t.Fatal(err)
		}
		tmp3 = append(tmp3, tmp2.FileId)
	}
	if ds, err := historyviews.GetHistory(10); err != nil {
		t.Fatal(err)
	} else {
		if len(ds) != 10 {
			t.Fatal("error")
		}
		counts := []int{}
		for i := 0; i < len(tmp3); i++ {
			flag := false
			for _, count := range counts {
				if count == tmp3[len(tmp3)-1-i] {
					flag = true
				}
			}
			if !flag {
				counts = append(counts, tmp3[len(tmp3)-1-i])
			}
		}
		for i := 0; i < len(ds); i++ {
			if ds[i].Id != uint(counts[i]) {
				t.Fatal(fmt.Sprintf("%v,%v", ds[i].Id, counts[i]))
			}
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
