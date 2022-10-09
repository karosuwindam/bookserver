package table

import "testing"

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
