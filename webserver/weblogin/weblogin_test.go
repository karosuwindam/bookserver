package weblogin

import (
	"bookserver/config"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {

}

func TestJwt(t *testing.T) {
	cfg, _ := config.EnvRead()
	_ = Setup(cfg)
	tokendata, err := createjwt("admin", "admin1234")
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	t.Logf("token: %v", tokendata)
	output, err := Unpackjwt(tokendata)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	if output.userName != "admin" {
		t.Fatalf("err name %v = %v", output.userName, "admin")
		t.FailNow()
	}
	t.Logf("name: %v,time: %v", output.userName, time.Unix(output.ext, 0))
	t.Log("----------------- Jwt check OK --------------------------")

}
