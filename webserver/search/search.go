package search

import (
	"bookserver/message"
	"bookserver/table"
	"bookserver/webserver/common"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type sqlSearch struct {
	sql *table.SQLStatus
}

const (
	IDKEYWRODPOINT = 1
)

func (cfg *sqlSearch) websqlsearchget(w http.ResponseWriter, r *http.Request) {
	out := message.Result{Name: "sql Search", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
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
	if urlPoint+IDKEYWRODPOINT > len(sUrl) {
		out.Option += "table not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else if urlPoint+IDKEYWRODPOINT == len(sUrl) || sUrl[urlPoint+IDKEYWRODPOINT] == "" {
		out.Option += "table=" + tName + " keyword not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else {
		keyword := sUrl[urlPoint+IDKEYWRODPOINT]
		out.Option += "table=" + tName + " keyword=" + keyword
		if jdata, err := cfg.sql.Search(tName, keyword); err != nil {
			out.Code = http.StatusNotFound
			log.Println(err.Error())
		} else {
			joutdata = fmt.Sprintf("%s", jdata)
			out.Result = ""
		}

	}
	common.Sqlreadmessageback(out, joutdata, w)

}

func (cfg *sqlSearch) websqlsearch(w http.ResponseWriter, r *http.Request) {

	switch strings.ToUpper(r.Method) {
	default:
		cfg.websqlsearchget(w, r)
	}
}

func WebSQLSearch(cfg interface{}, w http.ResponseWriter, r *http.Request) {
	switch cfg.(type) {
	case *table.SQLStatus:
		t := sqlSearch{sql: cfg.(*table.SQLStatus)}
		t.websqlsearch(w, r)
	}
}
