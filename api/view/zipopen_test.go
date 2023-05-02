package view

import (
	"bookserver/config"
	"context"
	"fmt"
	"testing"
	"time"
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

func TestZipCash(t *testing.T) {
	t.Setenv("ZIP_FILEPASS", "./zip")
	cfg, _ := config.EnvRead()
	Setup(cfg)

	filename := "test.zip"
	ctx, cancel := context.WithCancel(context.Background())
	t.Log("------------------------- zip open file cash start -------------------------------")
	if err := Add(filename); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	t.Log("Add Check OK")
	if err := Add(filename); err == nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	t.Log("Add Check Time OUT OK")
	go Loop(ctx)
	if err := Add(filename); err != nil {
		t.Fatalf("%v", err)
	}
	now := time.Now()
	time.Sleep(time.Second * 3)
	if cashZip[filename] == nil {
		t.Fatalf(filename)
		cancel()
		t.FailNow()
	}
	t.Log("Check Added OK")
	fmt.Printf("Sleep Start")
	for i := 0; time.Now().Sub(now).Minutes() < 1; i++ {
		time.Sleep(time.Second)
		fmt.Printf("%d\t", i)
	}
	fmt.Printf("\n")

	if cashZip[filename] != nil {
		t.Fatalf(filename)
	}
	t.Log("chash Clear OK")
	cancel()
	t.Log("------------------------- zip open file cash end -------------------------------")
}
