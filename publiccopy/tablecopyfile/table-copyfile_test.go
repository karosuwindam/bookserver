package tablecopyfile

import (
	"bookserver/config"
	"bookserver/table"
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {

	t.Setenv("DB_ROOTPASS", "./")
	t.Setenv("PDF_FILEPASS", "./")
	t.Setenv("ZIP_FILEPASS", "./")
	t.Setenv("PUBLIC_FILEPASS", "./")
	cfg, _ := config.EnvRead()
	if err := Setup(cfg); err != nil {
		t.FailNow()
	}
	sql.CreateTable()
	data := []table.Filelists{
		{Name: "a", Pdfpass: "a.pdf", Zippass: "a.zip"},
		{Name: "b", Pdfpass: "b.pdf", Zippass: "b.zip"},
		{Name: "test", Pdfpass: "test.pdf", Zippass: "test.zip"},
	}
	sql.Add(table.FILELIST, &data)

	data1 := []table.Copyfile{
		{Zippass: "a.zip", Filesize: 100, Copyflag: 1},
		{Zippass: "b.zip", Filesize: 100, Copyflag: 1},
	}
	sql.Add(table.COPYFILE, &data1)
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)
	t.Log("------------------------- Table Copy FIle start-------------------------------")
	if str, err := GetTableByName(1, PDF); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	} else if str != "a.pdf" {
		t.Fatalf("%s != a.pdf", str)
		t.FailNow()
	}
	if str, err := GetTableByName(2, ""); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	} else if str != "b.zip" {
		t.Fatalf("%s != b.zip", str)
		t.FailNow()
	}
	if str, err := GetTableByName(3, ZIP); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	} else {
		if err := CkCopyFileAll(); err != nil {
			t.Fatalf("%v+", err)
			t.FailNow()
		}
		if str, err := sql.ReadAll(table.COPYFILE); err != nil {
			t.Fatalf("%v+", err)
			t.FailNow()
		} else {
			if ary, ok := table.JsonToStruct(table.COPYFILE, []byte(str)).([]table.Copyfile); !ok || ary[0].Copyflag != 0 || ary[1].Copyflag != 0 {
				t.Fatalf("%s", str)
				t.FailNow()
			}
		}
		t.Log("No file change OK")
		AddCopyFIle("a.zip", true)
		AddCopyFIle("b.zip", true)
		if str, err := sql.ReadAll(table.COPYFILE); err != nil || str != "[]" {
			t.Fatalf("%v,%v+", str, err)
			t.FailNow()
		}
		t.Log("Not file chack OK")
		AddCopyFIle(str, true)
		if str, err := sql.ReadID(table.COPYFILE, 1); err != nil {
			t.Fatalf("%s != a.pdf", str)
			t.FailNow()
		} else {
			if ary, ok := table.JsonToStruct(table.COPYFILE, []byte(str)).([]table.Copyfile); !ok || ary[0].Copyflag != 1 {
				t.FailNow()
			}
		}
		t.Log("Add new row line OK")
		AddCopyFIle(str, false)
		if str, err := sql.ReadID(table.COPYFILE, 1); err != nil {
			t.Fatalf("%s != a.pdf", str)
			t.FailNow()
		} else {
			if ary, ok := table.JsonToStruct(table.COPYFILE, []byte(str)).([]table.Copyfile); !ok || ary[0].Copyflag != 0 {
				t.FailNow()
			}
		}
		t.Log("Edit new row line OK")
	}

	t.Log("------------------------- Table Copy FIle End-------------------------------")

}
