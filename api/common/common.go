package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// URLの解析
func UrlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

type Result struct {
	Name   string      `json:"Name"`
	Url    string      `json:"Url"`
	Code   int         `json:"Code"`
	Option string      `json:"Option"`
	Date   time.Time   `json:"Date"`
	Result interface{} `json:"Result"`
}

// 共通の出力設定
func CommonBack(msg Result, w http.ResponseWriter) {
	jsondata, _ := json.Marshal(msg)
	w.WriteHeader(msg.Code)
	fmt.Fprintf(w, "%v\n", string(jsondata))
}

// SQLの読み取り結果の出力
func Sqlreadmessageback(msg Result, w http.ResponseWriter) {
	jout, ok := msg.Result.(string)
	if !ok {
		jout = ""
	}
	msg.Result = ""
	jsondata, _ := json.Marshal(msg)
	w.WriteHeader(msg.Code)
	if jout != "" {
		out := strings.Replace(string(jsondata), "\"Result\":\"\"", "\"Result\":"+jout, -1)
		fmt.Fprintf(w, "%s", out)
	} else {
		fmt.Fprintf(w, "%v\n", string(jsondata))
	}

}

func Setup() error {
	return nil
}

// USERの状態で権限の確認
// ToDo
func CkLogin(msg *Result, w http.ResponseWriter, r *http.Request) bool {
	if true {
		return true
	} else {
		msg.Code = http.StatusUnauthorized
		msg.Result = "Not Login"
	}
	return false
}

// ファイルの存在確認
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
