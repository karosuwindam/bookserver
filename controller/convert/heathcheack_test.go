package convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	DataStoreInit()
	statusData.On()
	if statusData.Status != true {
		t.Fatal(fmt.Sprintln("error change data to true =", statusData.Status))
	}
	statusData.Change("aa")
	if statusData.Startfile["aa"] != "" || statusData.Endfile["aa"] != "aa" {
		t.Fatal(fmt.Sprintln("error change data a:", statusData.Startfile["aa"], "b:", statusData.Endfile["aa"]))

	}
	statusData.Add("aa")
	if statusData.Startfile["aa"] != "aa" || statusData.Endfile["aa"] != "aa" {
		t.Fatal(fmt.Sprintln("error change data a:", statusData.Startfile["aa"], "b:", statusData.Endfile["aa"]))

	}
	statusData.Change("aa")
	if statusData.Startfile["aa"] != "" || statusData.Endfile["aa"] != "aa" {
		t.Fatal(fmt.Sprintln("error change data a:", statusData.Startfile["aa"], "b:", statusData.Endfile["aa"]))
	}
	statusData.Off()
	if statusData.Status != false {
		t.Fatal(fmt.Sprintln("error change data to true =", statusData.Status))
	}
	var s interface{}
	s = CheackHealth()
	if b, err := json.Marshal(&s); err != nil {
		t.Fatal(err)
	} else {
		if string(b) != "{\"Status\":false,\"Startfile\":[],\"Endfile\":[\"aa\"]}" {
			t.Fatal("errordata", string(b))
		}
	}
	statusData.Clear()
	statusData.On()
	s = CheackHealth()
	if b, err := json.Marshal(&s); err != nil {
		t.Fatal(err)
	} else {
		if string(b) != "{\"Status\":true,\"Startfile\":[],\"Endfile\":[\"aa\"]}" {
			t.Fatal("errordata", string(b))
		}
	}
	statusData.StatusTIme["aa"] = time.Now().Add(-1 * time.Minute)
	statusData.Clear()
	s = CheackHealth()
	if b, err := json.Marshal(&s); err != nil {
		t.Fatal(err)
	} else {
		if string(b) != "{\"Status\":true,\"Startfile\":[],\"Endfile\":[]}" {
			t.Fatal("errordata", string(b))
		}
	}
}
