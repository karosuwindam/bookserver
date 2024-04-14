package table_test

import (
	"bookserver/config"
	"bookserver/table"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestVersionTable(t *testing.T) {
	os.Setenv("DB_ROOTPASS", "./db/")
	config.Init()
	defer func() {
		os.RemoveAll("./db")
	}()
	db, _ := openSqlite3()
	config.Version = "0.0.0"
	if table.VersionChack() != table.V0Table {
		t.Fatal("error chack version")
	}
	config.Version = "0.9.0"
	if table.VersionChack() != table.V1Table {
		t.Fatal("error chack version")
	}
	if err := table.InitVTable(db); err != nil {
		t.Fatal(err)
	}
	if err := table.InitVTable(db); err != nil {
		t.Fatal(err)
	}
	config.Version = "0.0.0"

	if err := table.InitVTable(db); err == nil {
		t.Fatal("Fatal chack table version")
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
