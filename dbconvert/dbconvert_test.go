package dbconvert

import (
	v0struct "bookserver/dbconvert/v0Struct"
	v1struct "bookserver/dbconvert/v1Struct"
	"log"
	"testing"
)

func TestDBConvert(t *testing.T) {
	filepass := "./db/orijinal_test.db"
	if _, err := readDBFile(filepass); err != nil {
		t.Fatal(err)
	}

	if d, err := v0struct.ReadBooknames(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
	if d, err := v1struct.ReadBooknames(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
	if d, err := v0struct.ReadCopyfiles(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
	if d, err := v1struct.ReadCopyfiles(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
	if d, err := v0struct.ReadFilelists(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
	if d, err := v1struct.ReadFilelists(basedb); err != nil {
		t.Fatal(err)
	} else {
		log.Println(d[0])
	}
}
