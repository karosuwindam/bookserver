package weblogin

import (
	"bookserver/message"
	"fmt"
	"net/http"
)

func getlogout(w http.ResponseWriter, r *http.Request) {
	msg := message.Message{Name: "login", Code: http.StatusOK}
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	nm, _ := ses.Values["name"].(string)
	if flg {
		msg.Status = "User:" + nm + " " + "LOGOUT OK"
	} else {
		msg.Status = "User:" + nm + " " + "LOGOUT NG"
	}
	ses.Values["login"] = nil
	ses.Values["name"] = nil
	if nm != "" {
		jwttmp[nm] = ""
	}
	ses.Save(r, w)
	fmt.Fprintf(w, "%v\n", msg.Output())
}

func WebServerLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		getlogout(w, r)
	}
}
