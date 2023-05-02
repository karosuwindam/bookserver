package read

import (
	"bookserver/api/common"
	"bookserver/table"
	"bookserver/webserver"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sql *table.SQLStatus

var apiname string = "read" //api名

// sqlreadlist(w, r) = common.Resule
//
// /read/:table名 テーブル内の全てのデータを取得する
func sqlreadlist(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""
	for i, url := range sUrl {
		if url == apiname {
			tName = sUrl[i+1]
			break
		}
	}
	if tName == "" {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		if jdata, err := sql.ReadAll(tName); err != nil {
			msg.Code = http.StatusNotFound
			msg.Result = err.Error()
		} else {
			msg.Result = fmt.Sprintf("%s", jdata)
		}
	}
	return msg
}

// sqlreadget(w, r) = common.Result
//
// /read/:tablename/:id IDが一致するデータを取得する
func sqlreadget(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""
	tId := ""
	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+2 {
			tName = sUrl[i+1]
			tId = sUrl[i+2]
			break
		}
	}
	if tName == "" || tId == "" {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		id, err := strconv.Atoi(tId)
		if err != nil {
			msg.Code = http.StatusNotFound
			msg.Result = []string{}
		} else {
			if jdata, err := sql.ReadID(tName, id); err != nil {
				msg.Code = http.StatusNotFound
				msg.Result = err.Error()
			} else {
				if jdata != "[]" {
					msg.Result = fmt.Sprintf("%s", jdata)
				} else {
					msg.Code = http.StatusBadRequest
				}
			}
		}

	}

	return msg
}

// websqlread(w, r)
//
// /read/の動作
func websqlread(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	switch strings.ToUpper(r.Method) {
	case "LIST":
		msg = sqlreadlist(w, r)
	default:
		msg = sqlreadget(w, r)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname + "/", Handler: websqlread},
}

// route動作について
func Setup(cfg *table.SQLStatus) []webserver.WebConfig {
	sql = cfg
	return route
}
