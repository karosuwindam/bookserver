package read

import (
	"bookserver/message"
	"bookserver/table"
	"bookserver/webserver/common"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type sqlRead struct {
	sql *table.SQLStatus
}

const (
	IDPOINT = 1
)

func (cfg *sqlRead) sqlreadlist(w http.ResponseWriter, r *http.Request) {
	out := message.Result{Name: "sql read", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
	sUrl := common.UrlAnalysis(r.URL.Path)
	urlPoint := 0
	tName := ""
	joutdata := ""
	for ; urlPoint < len(sUrl); urlPoint++ {
		if table.CkList(sUrl[urlPoint]) {
			tName = sUrl[urlPoint]
			break
		}
	}
	if urlPoint+IDPOINT > len(sUrl) {
		out.Option += "table not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else {
		out.Option += "table=" + tName
		if jdata, err := cfg.sql.ReadAll(tName); err != nil {
			out.Code = http.StatusNotFound
			log.Println(err.Error())
		} else {
			joutdata = fmt.Sprintf("%s", jdata)
			out.Result = ""
		}

	}
	common.Sqlreadmessageback(out, joutdata, w)
}
func (cfg *sqlRead) sqlreadget(w http.ResponseWriter, r *http.Request) {
	out := message.Result{Name: "sql read", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
	sUrl := common.UrlAnalysis(r.URL.Path)
	joutdata := ""
	urlPoint := 0
	tName := ""
	for ; urlPoint < len(sUrl); urlPoint++ {
		if table.CkList(sUrl[urlPoint]) {
			tName = sUrl[urlPoint]
			break
		}
	}
	if urlPoint+IDPOINT > len(sUrl) {
		out.Option += "table not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else if urlPoint+IDPOINT == len(sUrl) || sUrl[urlPoint+IDPOINT] == "" {
		out.Option += "table=" + tName + " id not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else {
		id, err := strconv.Atoi(sUrl[urlPoint+IDPOINT])
		if err != nil {
			out.Option += "table=" + tName + " id input error"
			out.Code = http.StatusNotFound
			out.Result = "[]"
		} else {
			out.Option += "table=" + tName + " id=" + sUrl[urlPoint+IDPOINT]
			if jdata, err := cfg.sql.ReadID(tName, id); err != nil {
				out.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				joutdata = fmt.Sprintf("%s", jdata)
				out.Result = ""
			}

		}

	}
	common.Sqlreadmessageback(out, joutdata, w)
}

func (cfg *sqlRead) websqlread(w http.ResponseWriter, r *http.Request) {

	switch strings.ToUpper(r.Method) {
	case "LIST":
		cfg.sqlreadlist(w, r)
	default:
		cfg.sqlreadget(w, r)
	}
}

func WebSQLRead(cfg interface{}, w http.ResponseWriter, r *http.Request) {
	switch cfg.(type) {
	case *table.SQLStatus:
		t := sqlRead{sql: cfg.(*table.SQLStatus)}
		t.websqlread(w, r)
	}
}
