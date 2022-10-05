package login

import (
	"bookserver/webserver"
	"net/http"
)

type loginbaseconging struct {
}

func setup() *loginbaseconging {
	return &loginbaseconging{}
}

func (t *loginbaseconging) webServerLogin(w http.ResponseWriter, r *http.Request) {

}
func webServerLogins(data interface{}, w http.ResponseWriter, r *http.Request) {
	switch data.(type) {
	case *loginbaseconging:
		data.(*loginbaseconging).webServerLogin(w, r)
	default:
	}
}

func Add(s *webserver.SetupServer) {
	t := setup()
	s.AddV1(webserver.GUEST|webserver.ADMIN|webserver.USER, "login", t, webServerLogins)
}
