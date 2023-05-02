package add

import (
	"bookserver/api/common"
	"bookserver/table"
	"bookserver/webserver"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var sql *table.SQLStatus

var apiname string = "add" //api名

// sqladd(w, r) = common.Result
//
// /add/:tablename/ 読み込んだデータを特定テーブルに登録する
func sqladd(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""

	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+1 {
			tName = sUrl[i+1]
			break
		}
	}

	if tName == "" {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		b, _ := io.ReadAll(r.Body)
		jout := table.JsonToStruct(tName, b)
		id := -1
		switch tName {
		case table.BOOKNAME:
			if tmp, ok := jout.(table.Booknames); ok && tmp.Name != "" {
				if err := sql.Add(tName, &tmp); err != nil {
					msg.Code = http.StatusBadRequest
					msg.Option = "NG"
					msg.Result = err.Error()
				} else {
					fmt.Println("Add database for", tName, "data:", tmp)
					id = tmp.Id
				}
			}
		case table.COPYFILE:
			if tmp, ok := jout.(table.Copyfile); ok && tmp.Zippass != "" {
				if err := sql.Add(tName, &tmp); err != nil {
					msg.Code = http.StatusBadRequest
					msg.Option = "NG"
					msg.Result = err.Error()
				} else {
					fmt.Println("Add database for", tName, "data:", tmp)
					id = tmp.Id
				}
			}
		case table.FILELIST:
			if tmp, ok := jout.(table.Filelists); ok && tmp.Name != "" {
				if err := sql.Add(tName, &tmp); err != nil {
					msg.Code = http.StatusBadRequest
					msg.Option = "NG"
					msg.Result = err.Error()
				} else {
					fmt.Println("Add database for", tName, "data:", tmp)
					id = tmp.Id
				}
			}
		default:
			msg.Code = http.StatusBadRequest
			msg.Option = "NG"
		}
		if id != -1 {
			if jdata, err := sql.ReadID(tName, id); err != nil {
				msg.Code = http.StatusNotFound
				msg.Result = err.Error()
			} else {
				msg.Result = fmt.Sprintf("%s", jdata)
			}
		} else {
			msg.Code = http.StatusBadRequest
			msg.Option = "NG"
			msg.Result = []string{}
		}
	}

	return msg
}

// websqladd(w, r)
//
// /add/の動作
func websqladd(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	if common.CkLogin(&msg, w, r) {
		switch strings.ToUpper(r.Method) {
		case "POST":
			msg = sqladd(w, r)
		default:
		}
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname + "/", Handler: websqladd},
}

// route動作について
func Setup(cfg *table.SQLStatus) []webserver.WebConfig {
	sql = cfg
	return route
}
