package add

import (
	"bookserver/message"
	"bookserver/table"
	"bookserver/webserver/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type sqlAdd struct {
	sql *table.SQLStatus
}

func (cfg *sqlAdd) sqladd(w http.ResponseWriter, r *http.Request) {
	out := message.Result{Name: "sql add", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
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
	if urlPoint > len(sUrl) {
		out.Option += "table not input"
		out.Code = http.StatusNotFound
		out.Result = "[]"
	} else {
		b, _ := io.ReadAll(r.Body)

		out.Option += "table=" + tName
		switch tName {
		case table.BOOKNAME:
			jout := table.Booknames{}
			if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
				log.Println(err.Error())
				out.Option += " NG"
			} else {
				if err := cfg.sql.Add(tName, &jout); err != nil {
					log.Println(err.Error())
					out.Option += " NG"
				} else {
					if jread, err := cfg.sql.ReadID(tName, jout.Id); err != nil {
						log.Println(err.Error())
						out.Option += " NG"
					} else {
						out.Option += " OK"
						joutdata = fmt.Sprintf("%s", jread)
						out.Result = ""
					}
				}
			}
		case table.FILELIST:
			jout := table.Filelists{}
			if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
				log.Println(err.Error())
				out.Option += " NG"
			} else {
				if err := cfg.sql.Add(tName, &jout); err != nil {
					log.Println(err.Error())
					out.Option += " NG"
				} else {
					if jread, err := cfg.sql.ReadID(tName, jout.Id); err != nil {
						log.Println(err.Error())
						out.Option += " NG"
					} else {
						out.Option += " OK"
						joutdata = fmt.Sprintf("%s", jread)
						out.Result = ""
					}
				}
			}
		case table.COPYFILE:
			jout := table.Copyfile{}
			if err := json.Unmarshal(b, &jout); err != nil {
				log.Println(err.Error())
				out.Option += " NG"
			} else {
				if err := cfg.sql.Add(tName, &jout); err != nil && jout.Zippass == "" {
					log.Println(err.Error())
					out.Option += " NG"
				} else {
					if jread, err := cfg.sql.ReadID(tName, jout.Id); err != nil {
						log.Println(err.Error())
						out.Option += " NG"
					} else {
						out.Option += " OK"
						joutdata = fmt.Sprintf("%s", jread)
						out.Result = ""
					}
				}
			}
		}
		// if jdata, err := cfg.sql.ReadAll(tName); err != nil {
		// 	out.Code = http.StatusNotFound
		// 	log.Println(err.Error())
		// } else {
		// 	joutdata = fmt.Sprintf("%s", jdata)
		// 	out.Result = ""
		// }

	}
	common.Sqlreadmessageback(out, joutdata, w)
}
func (cfg *sqlAdd) websqladd(w http.ResponseWriter, r *http.Request) {

	switch strings.ToUpper(r.Method) {
	case "POST":
		cfg.sqladd(w, r)
	default:
	}
}

func WebSQLRead(cfg interface{}, w http.ResponseWriter, r *http.Request) {
	switch cfg.(type) {
	case *table.SQLStatus:
		t := sqlAdd{sql: cfg.(*table.SQLStatus)}
		t.websqladd(w, r)
	}
}
