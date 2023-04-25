package publiccopy

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	input = make(chan interface{}, 1)
	if err := Add(""); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}
	run(<-input)
	if err := Add(""); err != nil {
		t.Fatalf("%v+", err)
		t.FailNow()
	}

	if err := Add(""); err == nil {
		t.FailNow()
	} else {
		t.Logf("%v", err)

	}

}

func TestEnd(t *testing.T) {
	input = make(chan interface{}, 1)
	endch = make(chan error, 1)
	Add("")
	if err := End(); err == nil {
		t.FailNow()
	} else {
		fmt.Println(err)
	}
	endch <- nil
	if err := End(); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
	run(<-input)
	if err := End(); err != nil {
		t.Fatalf("%v", err)
		t.FailNow()
	}
}
