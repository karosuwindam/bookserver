package view

import (
	"bookserver/config"
	"fmt"
	"testing"
)

func TestZipOpen(t *testing.T) {

	t.Setenv("ZIP_FILEPASS", "./zip")
	cfg, _ := config.EnvRead()
	Setup(cfg)
	filename := "test.zip"
	t.Log("------------------------- zip open start -------------------------------")
	f := openfile(filename)
	if !f.flag {
		t.Fatalf("Not Open%v", filename)
		t.FailNow()
	}
	fmt.Println("pass:", f.path, "data:", f.DataName, "datacount", f.FileCount)
	t.Log("------------------------- zip open end -------------------------------")

}

func TestZipOpenFile(t *testing.T) {

	t.Setenv("ZIP_FILEPASS", "./zip")
	cfg, _ := config.EnvRead()
	Setup(cfg)
	filename := "test.zip"
	t.Log("------------------------- zip open file start -------------------------------")
	f := openfile(filename)
	for _, name := range f.DataName {
		if _, err := f.openZipRead(name); err != nil {
			t.Fatalf("Not Open%v", name)
			t.FailNow()
		}

	}
	fmt.Println(f.convertjson())
	t.Log("------------------------- zip open file end -------------------------------")

}
