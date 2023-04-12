package search

import (
	"bookserver/api/common"
	"bookserver/table"
	"bookserver/webserverv2"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var sql *table.SQLStatus

var key map[string]bool //特殊キーワードの取得リスト
var keynames []string = []string{
	"today", "toweek", "tomonth",
} //特殊キーワードリスト

var apiname string = "search" //api名

// websqlsearchget(w, r) = common.Result
//
// /read/:tablename/:keyword Keywordがあるデータを対象のテーブルから取得する
//
// 特殊キーワード today toweek tomonth
func websqlsearchget(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""
	tKeyword := ""

	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+2 {
			tName = sUrl[i+1]
			tKeyword = sUrl[i+2]
			break
		}
	}
	if tName == "" || tKeyword == "" {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		if key[tKeyword] {
			if jdata, err := sql.ReadWhileTime(tName, tKeyword); err != nil {
				msg.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				msg.Result = fmt.Sprintf("%s", jdata)
			}

		} else {
			if jdata, err := sql.Search(tName, tKeyword); err != nil {
				msg.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				msg.Result = fmt.Sprintf("%s", jdata)
			}

		}

	}

	return msg
}

// websqlread(w, r)
//
// /search/の動作
func websqlsearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	switch strings.ToUpper(r.Method) {
	default:
		msg = websqlsearchget(w, r)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/" + apiname + "/", Handler: websqlsearch},
}

// route動作について
func Setup(cfg *table.SQLStatus) []webserverv2.WebConfig {
	key = map[string]bool{}
	for _, keyname := range keynames {
		key[keyname] = true
	}
	sql = cfg
	return route
}
