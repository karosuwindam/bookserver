package listdata

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/table"
	"bookserver/webserver"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var apiname string = "listdata" //api名
var zipfolder string
var pdffolder string

// var imgfolder string
var sql *table.SQLStatus

// webZipRead
//
// /list/の動作
func webZipRead(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}
	switch strings.ToUpper(r.Method) {
	default:
		tmp, _ := json.Marshal(readListData())
		msg.Result = string(tmp)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path

	common.Sqlreadmessageback(msg, w)

}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: webZipRead},
	{Pass: "/" + apiname + "/", Handler: webZipRead},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	zipfolder = cfg.Folder.Zip
	pdffolder = cfg.Folder.Pdf
	// imgfolder = cfg.Folder.Img
	listData = []ListData{}
	if sqlcfg, err := table.Setup(cfg); err != nil {
		return nil, err
	} else {
		sql = sqlcfg
	}
	return route, nil
}
