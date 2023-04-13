package login

import (
	"bookserver/api/common"
	"net/http"
)

func getlogout(w http.ResponseWriter, r *http.Request) {
	msg := common.Result{Name: "logout", Code: http.StatusOK}
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	nm, _ := ses.Values["name"].(string)
	if flg {
		msg.Option = "User:" + nm + " " + "LOGOUT OK"
	} else {
		msg.Option = "User:" + nm + " " + "LOGOUT NG"
	}
	ses.Values["login"] = nil
	ses.Values["name"] = nil
	if nm != "" {
		jwttmp[nm] = ""
	}
	ses.Save(r, w)
	common.CommonBack(msg, w)
}

func webServerLogout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		getlogout(w, r)
	}
}
