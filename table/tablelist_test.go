package table

import (
	"testing"
)

func TestTablelistSetup(t *testing.T) {
	tablelistsetup()
	if !ckType(&Booknames{}) {
		t.FailNow()
	}
	if !ckType(&Copyfile{}) {
		t.FailNow()
	}
	if !ckType(&Filelists{}) {
		t.FailNow()
	}
	for tname, _ := range tablelist {
		readBaseCreate(tname)
	}

}

func TestCreateSerchKeyword(t *testing.T) {
	tablelistsetup()
	booksrdata := createSerchText(BOOKNAME, "a")
	ckbookdata := []string{"booktype", "burand", "ext", "name", "title", "writer"}
	for i := 1; i < len(ckbookdata); i++ {
		if booksrdata[ckbookdata[i]] != booksrdata[ckbookdata[i-1]] {
			t.Fatalf("%v != %v", booksrdata[ckbookdata[i]], booksrdata[ckbookdata[i-1]])
			t.FailNow()
		}
	}
	t.Log(booksrdata)
	copysrdata := createSerchText(COPYFILE, "b")
	ckcopydata := []string{"zippass"}
	for i := 1; i < len(ckcopydata); i++ {
		if copysrdata[ckcopydata[i]] != copysrdata[ckcopydata[i-1]] {
			t.Fatalf("%v != %v", copysrdata[ckcopydata[i]], copysrdata[ckcopydata[i-1]])
			t.FailNow()
		}
	}
	t.Log(copysrdata)
	filesrdata := createSerchText(FILELIST, "c")
	ckfiledata := []string{"name", "pdfpass", "tag", "zippass"}
	for i := 1; i < len(ckfiledata); i++ {
		if filesrdata[ckfiledata[i]] != filesrdata[ckfiledata[i-1]] {
			t.Fatalf("%v != %v", filesrdata[ckfiledata[i]], filesrdata[ckfiledata[i-1]])
			t.FailNow()
		}
	}
	t.Log(filesrdata)

}

func TestJsonToStruct(t *testing.T) {
	jsondata := readfile("./jsonsample/booknames.json")
	if tmp := JsonToStruct(BOOKNAME, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Booknames)
		if !ok {
			t.FailNow()
		}
	}
	jsondata = readfile("./jsonsample/copyfile.json")
	if tmp := JsonToStruct(COPYFILE, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Copyfile)
		if !ok {
			t.FailNow()
		}
	}
	jsondata = readfile("./jsonsample/filelist.json")
	if tmp := JsonToStruct(FILELIST, []byte(jsondata)); tmp == nil {
		t.FailNow()
	} else {
		_, ok := tmp.(Filelists)
		if !ok {
			t.FailNow()
		}
	}
}
