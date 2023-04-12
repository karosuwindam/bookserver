package edit

import (
	"bookserver/api/common"
	"bookserver/table"
	"bookserver/webserverv2"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sql *table.SQLStatus

var apiname string = "edit" //api名

func sqleditwrite(tName string, b []byte, id int) (string, error) {
	switch tName {
	case table.BOOKNAME:
		jout := table.Booknames{}
		if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
		} else {
			if jread, err := sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	case table.FILELIST:
		jout := table.Filelists{}
		if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
		} else {
			if jread, err := sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	case table.COPYFILE:
		jout := table.Copyfile{}
		if err := json.Unmarshal(b, &jout); err != nil {
		} else {
			if jread, err := sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	}
	return "", errors.New("Data NG")
}

// sqldelete(w, r) = common.Result
//
// /edit/:tablename/:id IDが一致するデータを削除する
func sqldelete(w http.ResponseWriter, r *http.Request) common.Result {
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
		if id, err := strconv.Atoi(tId); err == nil {
			if _, err := sql.Delete(tName, id); err != nil {
				msg.Code = http.StatusNotFound
				msg.Result = err.Error()
			} else {
				fmt.Printf("delete %v OK", id)
				msg.Result = fmt.Sprintf("delete %v OK", id)
			}
		} else {
			msg.Code = http.StatusBadRequest
			msg.Result = err.Error()
		}
	}
	return msg

}

func sqledit(w http.ResponseWriter, r *http.Request) common.Result {
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
			b, _ := io.ReadAll(r.Body)
			if jdata, err := sqleditwrite(tName, b, id); err != nil {
				msg.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				if jdata != "[]" {
					fmt.Println("Edit database for", tName, "id=", id, "data:", string(b))
					msg.Result = fmt.Sprintf("%s", jdata)
				} else {
					msg.Code = http.StatusBadRequest

				}
			}

		}

	}
	return msg

}

// sqleditget(w, r) = common.Result
//
// /edit/:tablename/:id IDが一致するデータを取得する
func sqleditget(w http.ResponseWriter, r *http.Request) common.Result {
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

// websqledit(w, r)
//
// /edit/の動作
func websqledit(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	if common.CkLogin(&msg, w, r) {
		switch strings.ToUpper(r.Method) {
		case "POST":
			msg = sqledit(w, r)
		case "DELETE":
			msg = sqldelete(w, r)
		default:
			msg = sqleditget(w, r)
		}
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/" + apiname + "/", Handler: websqledit},
}

// route動作について
func Setup(cfg *table.SQLStatus) []webserverv2.WebConfig {
	sql = cfg
	return route
}
