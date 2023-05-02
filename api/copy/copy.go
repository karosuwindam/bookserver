package copy

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/publiccopy"
	"bookserver/table"
	"bookserver/webserver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var apiname string = "copy" //api名

var sql *table.SQLStatus

// webcopyPost(w,r)
//
// /copyのPOST処理
func webcopyPost(w http.ResponseWriter, r *http.Request, msg *common.Result) {
	b, _ := io.ReadAll(r.Body)
	msg.Option += ":" + string(b)
	fmt.Println(r.Method, r.URL.Path, "data:", string(b))
	tmp := publiccopy.CopyFils{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		msg.Code = http.StatusBadRequest
		msg.Result = "NG"
		log.Println(err)
		return
	}
	if err := publiccopy.Add(tmp); err != nil {
		msg.Code = http.StatusRequestTimeout
		msg.Result = "Time out"
		log.Println(err)
		return
	}
	msg.Result = "OK"
}

// /copy/:idのGET処理
func webcopyGet(w http.ResponseWriter, r *http.Request, msg *common.Result) {
	sUrl := common.UrlAnalysis(r.URL.Path)
	tId := ""

	for i, url := range sUrl {
		if url == apiname && len(sUrl) > i+1 {
			tId = sUrl[i+1]
			break
		}
	}
	if tId == "" {
		msg.Code = http.StatusNotFound
		msg.Result = []string{}
	} else {
		if id, err := strconv.Atoi(tId); err != nil {
			log.Println(err)
			msg.Code = http.StatusNotFound
			msg.Result = []string{}
		} else {
			if jdata, err := sql.ReadID(table.FILELIST, id); err == nil && jdata != "[]" {
				if tmp, ok := table.JsonToStruct(table.FILELIST, []byte(jdata)).([]table.Filelists); ok {
					if output, err := sql.ReadName(table.COPYFILE, tmp[0].Zippass); err == nil {
						if output == "[]" {
							msg.Code = http.StatusNotFound
						}
						msg.Result = output
						return
					}
				}
			}
			msg.Code = http.StatusNotFound
			msg.Result = []string{}
		}
	}
}

// webcopy(w,r)
//
// /copyの動作
func webcopy(w http.ResponseWriter, r *http.Request) {
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	msg.Name = apiname
	msg.Url = r.URL.Path

	switch strings.ToUpper(r.Method) {
	case "POST":
		webcopyPost(w, r, &msg)
		common.CommonBack(msg, w)
	default:
		fmt.Println(r.Method, r.URL.Path)
		webcopyGet(w, r, &msg)
		common.Sqlreadmessageback(msg, w)
	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: webcopy},
	{Pass: "/" + apiname + "/", Handler: webcopy},
}

// SetupSetup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	if scfg, err := table.Setup(cfg); err != nil {
		return route, err
	} else {
		sql = scfg
	}
	return route, nil
}
