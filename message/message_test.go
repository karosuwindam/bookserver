package message

import (
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	msg := Message{Name: "data", Status: "status", Code: 200}
	ck := "{\"name\":\"data\",\"status\":\"status\",\"code\":200}"

	if ck != msg.Output() {
		t.Error(msg.Output())
		t.FailNow()
	}
	t.Log(msg.Output())

}

func TestReult(t *testing.T) {
	msg := Message{Name: "data", Status: "status", Code: 200}
	rst := Result{Name: "data", Code: 200, Option: "", Date: time.Now(), Result: msg}
	t.Log(rst.Output())
}
