package search

import (
	"bookserver/api/common"
	"bookserver/table"
	"bookserver/webserverv2"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var sql *table.SQLStatus

var key map[string]bool //特殊キーワードの取得リスト
var keynames []string = []string{
	"today", "toweek", "tomonth", "rand",
} //特殊キーワードリスト

var apiname string = "search" //api名

type SearchKey struct {
	Table   string `json:"Table"`
	Keyword string `json:"Keyword"`
}

// websqlSearchPost(w,r) = common.Result
//
// /searchでSearchKey形式でPostしたデータを対象のテーブルから検索して取得する
func websqlSearchPost(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now()}
	b, _ := io.ReadAll(r.Body)
	fmt.Println("", string(b))
	msg.Option = r.Method + "," + string(b)
	jout := SearchKey{}
	if err := json.Unmarshal(b, &jout); err != nil || jout.Table == "" {
		msg.Code = http.StatusNotFound
		if err != nil {
			msg.Result = err.Error()
		}

	} else {
		if key[jout.Keyword] {
			if jout.Keyword == "rand" {
				if jdata, err := sql.ReadAll(jout.Table); err != nil {
					msg.Code = http.StatusNotFound
					log.Println(err.Error())
				} else if jdata == "[]" {
					msg.Result = fmt.Sprintf("%s", jdata)
				} else {
					msg.Result = fmt.Sprintf("%s", table.RandGenerate(table.JsonToStruct(jout.Table, []byte(jdata))))
				}
			} else {
				if jdata, err := sql.ReadWhileTime(jout.Table, jout.Keyword); err != nil {
					msg.Code = http.StatusNotFound
					log.Println(err.Error())
				} else {
					msg.Result = fmt.Sprintf("%s", jdata)
				}
			}

		} else {
			if jdata, err := sql.Search(jout.Table, jout.Keyword); err != nil {
				msg.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				msg.Result = fmt.Sprintf("%s", jdata)
			}

		}
	}
	return msg

}

// websqlsearchget(w, r) = common.Result
//
// /search/:tablename/:keyword Keywordがあるデータを対象のテーブルから取得する
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
	fmt.Printf("%v %v", r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	switch strings.ToUpper(r.Method) {
	case "POST":
		msg = websqlSearchPost(w, r)
	default:
		msg = websqlsearchget(w, r)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/" + apiname, Handler: websqlsearch},
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
