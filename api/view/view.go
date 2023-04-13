package view

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/table"
	"bookserver/webserverv2"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sql *table.SQLStatus

var apiname string = "view"       //api名
var apinameimage string = "image" //api名
var zippath string                //zipのフォルダパス

// webzipread(w, r)
//
// listの取得
func webzipreadlist(w http.ResponseWriter, r *http.Request) common.Result {
	msg := common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	sUrl := common.UrlAnalysis(r.URL.Path)
	tId := ""
	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+1 {
			tId = sUrl[i+1]
			break
		}
	}
	if id, err := strconv.Atoi(tId); err != nil {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		if json, err := sql.ReadID(table.FILELIST, id); err != nil {
			msg.Code = http.StatusNotFound
			msg.Result = []string{}
		} else {
			if data, ok := table.JsonToStruct(table.FILELIST, []byte(json)).([]table.Filelists); ok {
				zipname := data[0].Zippass
				f := openfile(zipname)
				if f.flag {
					msg.Result = f.convertjson()
				} else {
					msg.Code = http.StatusNotFound
					msg.Result = []string{}
				}
			}
		}
	}
	return msg
}

// webzipread(w, r)
//
// /view/:idの動作
func webzipread(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	switch strings.ToUpper(r.Method) {
	default:
		msg = webzipreadlist(w, r)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)
}

var route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/" + apiname + "/", Handler: webzipread},
	{Pass: "/" + apinameimage + "/", Handler: webzipreadimage},
}

// Setup
func Setup(cfg *config.Config) ([]webserverv2.WebConfig, error) {
	zippath = cfg.Folder.Zip
	if zippath[len(zippath)-1:] != "/" {
		zippath += "/"
	}
	if sqlcfg, err := table.Setup(cfg); err != nil {
		return nil, err
	} else {
		sql = sqlcfg
	}
	return route, nil
}
