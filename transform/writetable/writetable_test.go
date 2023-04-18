package writetable

import (
	"bookserver/config"
	"bookserver/table"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	if err := Setup(cfg); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
	t.Log("OK")
}

func TestWriteTable(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	if err := Setup(cfg); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
	t.Log("---------------Create Table Start ---------------")
	str := "test.pdf"
	tt := PdftoZip{Name: "test", InputFile: str, OutputFile: "[テスト1]テスト.zip", Tag: "テスト,テスト1,漫画,ああああ"}
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		if tmp != tt {
			t.Fatalf("%v!=%v", tmp, tt)
			t.FailNow()
		}
	}
	str = "test1.pdf"
	tt = PdftoZip{Name: "test1", InputFile: str, OutputFile: "[テスト1]テスト01.zip", Tag: "テスト,テスト1,漫画,ああああ"}
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		if tmp != tt {
			t.Fatalf("%v!=%v", tmp, tt)
			t.FailNow()
		}
	}
	str = "test01.pdf"
	tt = PdftoZip{Name: "test01", InputFile: str, OutputFile: "[テスト1]テスト01.zip", Tag: "テスト,テスト1,漫画,ああああ"}
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		if tmp != tt {
			t.Fatalf("%v!=%v", tmp, tt)
			t.FailNow()
		}
	}
	str = "aa.pdf"
	tt = PdftoZip{Name: "aa", InputFile: str, OutputFile: "aa.zip", Tag: ""}
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		if tmp != tt {
			t.Fatalf("%v!=%v", tmp, tt)
			t.FailNow()
		}
	}
	str = "aa1.pdf"
	tt = PdftoZip{Name: "aa1", InputFile: str, OutputFile: "aa1.zip", Tag: ""}
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		if tmp != tt {
			t.Fatalf("%v!=%v", tmp, tt)
			t.FailNow()
		}
	}
	t.Log("---------------Create Table End OK ---------------")
}

func TestAddFileTable(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	if err := Setup(cfg); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	}
	t.Log("---------------Add Table Start ---------------")
	str := "test.pdf"
	if tmp, err := CreatePdfToZip(str); err != nil {
		t.Fatalf("%v", err.Error())
		t.FailNow()
	} else {
		ptmp := &table.Filelists{Name: tmp.Name}
		if err := AddFileTable(ptmp); err != nil {
			t.Fatalf("%v", err.Error())
			t.FailNow()
		}
		if jj, _ := sql.ReadID(table.FILELIST, 1); jj != "[]" {
			if tmpr, ok := table.JsonToStruct(table.FILELIST, []byte(jj)).([]table.Filelists); ok {
				if tmpr[0] != *ptmp {
					sql.Delete(table.FILELIST, 1)
					t.Fatalf("%v!=%v", tmpr[0], *ptmp)
					t.FailNow()
				}
			}
		} else {
			t.Fatalf("test.db error")
			t.FailNow()
		}
		ptmp.Pdfpass = tmp.InputFile
		ptmp.Zippass = tmp.OutputFile
		ptmp.Tag = tmp.Tag
		if err := AddFileTable(ptmp); err != nil {
			sql.Delete(table.FILELIST, 1)
			t.Fatalf("%v", err.Error())
			t.FailNow()
		}
		if jj, _ := sql.ReadID(table.FILELIST, 1); jj != "[]" {
			if tmpr, ok := table.JsonToStruct(table.FILELIST, []byte(jj)).([]table.Filelists); ok {
				if tmpr[0] != *ptmp {
					sql.Delete(table.FILELIST, 1)
					t.Fatalf("%v!=%v", tmpr[0], *ptmp)
					t.FailNow()
				}
			}
		} else {
			t.Fatalf("test.db error")
			t.FailNow()
		}
		sql.Delete(table.FILELIST, 1)

	}
	t.Log("---------------Add Table End OK ---------------")
}
