package dirread

import "testing"

func TestDirpassCheck(t *testing.T) {
	if !dirpasscheck("./") {
		t.Fatalf("Don't pass")
		t.FailNow()
	}
	if dirpasscheck("./dirreead.go") {
		t.Fatalf("pass check miss")
		t.FailNow()

	}
}

func TestDirread(t *testing.T) {
	readdate, _ := Setup("./")
	if err := readdate.Read("./"); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	for _, datafile := range readdate.Data {
		t.Log(datafile)

	}
}
