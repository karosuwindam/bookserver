package weblogin

import (
	"bookserver/message"
	"bookserver/webserver/common"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type loginbaseconging struct {
}

var psckdata map[string]string

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte("secret-key-12345"))

func Setup() *loginbaseconging {
	psckdata = map[string]string{}
	psckdata["admin"] = "admin1234"
	return &loginbaseconging{}
}

func login(w http.ResponseWriter, r *http.Request) {
	ses, _ := cs.Get(r, "hello-session")
	ses.Values["login"] = nil
	ses.Values["name"] = nil
	nm := r.PostFormValue("name")
	pw := r.PostFormValue("pass")
	if psckdata[nm] == pw {
		ses.Values["login"] = true
		ses.Values["name"] = nm
	}
	ses.Save(r, w)
	getlogin(w, r)
}

func getlogin(w http.ResponseWriter, r *http.Request) {
	msg := message.Message{Name: "login", Code: http.StatusOK}
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	nm, _ := ses.Values["name"].(string)
	if flg {
		msg.Status = "User:" + nm + " " + "LOGIN OK"
	} else {
		msg.Status = "User:" + nm + " " + "LOGIN NG"
	}
	fmt.Fprintf(w, "%v\n", msg.Output())
}

func CkUserlogin(name string, r *http.Request) common.UserType {
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	if flg {
		return common.ADMIN | common.GUEST | common.USER
	}
	return common.GUEST
}

func (t *loginbaseconging) webServerLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		login(w, r)
	default:
		getlogin(w, r)
	}
}
func WebServerLogin(data interface{}, w http.ResponseWriter, r *http.Request) {
	switch data.(type) {
	case *loginbaseconging:
		data.(*loginbaseconging).webServerLogin(w, r)
	default:
	}
}
