package readzipfile

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCashTest(t *testing.T) {
	if err := cashInit(); err == nil {
		t.Fatal("set up error for zipPath")
	}
	zipPath = "./zip/"
	if err := cashInit(); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCashData("test.zip", "testout-000.jpg"); err == nil {
		t.Fatal("not add table")
	}
	if _, err := readZipFileData("test.zip", "testout-000.jpg"); err != nil {
		t.Fatal(err)
	}
	if _, err := readZipFileData("test.zip", "testout-100.jpg"); err == nil {
		t.Fatal("error data")
	} else {
		log.Println("info:", err)
	}

	if err := readZipFileAll("test.zip", context.TODO()); err != nil {
		t.Fatal(err)
	}
	if dataStore.cashZipSize["test.zip"] != 23002 {
		t.Fatal("add cash data err")
	}

	if _, err := ReadCashData("test.zip", "testout-000.jpg"); err != nil {
		t.Fatal("not add table")
	}
	tmp := []string{}
	for s, _ := range dataStore.cashZip["test.zip"].buf {
		tmp = append(tmp, s)
	}
	ary := []string{"testout-000.jpg", "testout-001.jpg", "testout-002.jpg"}
	for i, tt := range ary {
		if tmp[i] != tt {
			t.Fatal(fmt.Sprintf("error data %v != %v", tmp[i], tt))
		}
	}
	if err := clearZipFileCash(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if dataStore.cashZip["test.zip"] == nil {
		t.Fatal("time Clear cash error")
	}

	// time.Sleep(1 * time.Minute) //1分経過待ち
	dataStore.cashZipTime["test.zip"] = time.Now().Add(-2 * time.Minute)

	if err := clearZipFileCash(context.TODO()); err != nil {
		t.Fatal(err)
	}
	if dataStore.cashZip["test.zip"] != nil {
		t.Fatal("error not clear")
	}
}
