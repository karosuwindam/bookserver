package healthmessage

import (
	"testing"
)

func TestCreateMessage(t *testing.T) {
	m := Create("test")
	if m.Name != "test" {
		t.Fatalf("err create message")
		t.FailNow()
	}
	m.ChangeMessage(true, "OK")
	if m.Flag != true {
		t.Fatalf("not change flag")
		t.FailNow()
	}
	if m.Message != "OK" {
		t.Fatalf("not change message")
		t.FailNow()
	}
	if "{\"Name\":\"test\",\"Flag\":true,\"Message\":\"OK\"}" != m.JsonOut() {
		t.Fatalf("error json")
		t.FailNow()
	}
}
