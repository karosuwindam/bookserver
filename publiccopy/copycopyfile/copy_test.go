package copycopyfile

import (
	"bookserver/config"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Setenv("PDF_FILEPASS", "./pdf")
	t.Setenv("ZIP_FILEPASS", "./zip")
	cfg, _ := config.EnvRead()
	t.Log("------------------ Copy Start -------------------------")
	if err := Setup(cfg); err != nil {
		t.FailNow()
	}
	filename := "test.zip"
	if err := CopyFile(filename); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}
	if !Exists(publicPass + filename) {
		t.Fatalf("Not File %s", publicPass+filename)
		t.FailNow()
	}
	if err := RemoveFile(filename); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}
	filename = "test.pdf"
	if err := CopyFile(filename); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}
	if !Exists(publicPass + filename) {
		t.Fatalf("Not File %s", publicPass+filename)
		t.FailNow()
	}
	if err := RemoveFile(filename); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}
	t.Log("------------------ Copy End -------------------------")

}
