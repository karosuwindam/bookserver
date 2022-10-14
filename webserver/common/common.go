package common

import (
	"bookserver/config"
	"bookserver/message"
	"bookserver/table"
	"fmt"
	"net/http"
	"strings"
)

type UserType int

const (
	ADMIN UserType = 1
	USER  UserType = 1 << 1
	GUEST UserType = 1 << 2
)

type WebServerConfig struct {
	Sql *table.SQLStatus
}

func Setup(cfg *config.Config) (*WebServerConfig, error) {
	out := &WebServerConfig{}
	if sql, err := table.Setup(cfg); err != nil {
		return nil, err
	} else {
		out.Sql = sql
	}
	return out, nil
}

// URLの解析
func UrlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

// SQLの読み取り結果の出力
func Sqlreadmessageback(out message.Result, joutdata string, w http.ResponseWriter) {
	jout := out.Output()
	w.WriteHeader(out.Code)
	if joutdata != "" {
		jout = strings.Replace(jout, "\"result\":\"\"", "\"result\":"+joutdata, -1)
		fmt.Fprintf(w, "%s", jout)

	} else {
		fmt.Fprintf(w, "%s", jout)

	}

}
