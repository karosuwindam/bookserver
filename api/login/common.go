package login

import (
	"bookserver/config"
	"bookserver/webserverv2"

	"github.com/gorilla/sessions"
)

type UserType int

const (
	ADMIN UserType = 1
	USER  UserType = 1 << 1
	GUEST UserType = 1 << 2
)

var apiname string = "login" //api名

type tResurlt struct {
	status string `json:"status"`
	Token  string `json:"token"`
}

type JwtData struct {
	userName string   `json:"user_name"`
	usertype UserType `json:"user_type"`
	ext      int64    `json:"ext"`
}

var psckdata map[string]string = map[string]string{}
var jwttmp map[string]string = map[string]string{}
var jwtsecretkey string
var cs *sessions.CookieStore = sessions.NewCookieStore([]byte("secret-key-12345"))

// routeのベースフォルダ
var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{"/" + apiname, webServerLogin},
	{"/" + "logout", webServerLogout},
}

func Setup(cfg *config.Config) ([]webserverv2.WebConfig, error) {
	psckdata["admin"] = "admin1234"
	str, _ := createjwt("admin", "admin1234")
	Unpackjwt(str)
	jwtsecretkey = cfg.SeretKey.JwtKey

	return route, nil

}
