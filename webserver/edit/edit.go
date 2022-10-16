package edit

import (
	"bookserver/message"
	"bookserver/table"
	"bookserver/webserver/common"
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

type sqlEdit struct {
	sql *table.SQLStatus
}

const (
	IDPOINT = 1
)

func (cfg *sqlEdit) sqleditwrite(tName string, b []byte, id int) (string, error) {
	switch tName {
	case table.BOOKNAME:
		jout := table.Booknames{}
		if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
		} else {
			if jread, err := cfg.sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	case table.FILELIST:
		jout := table.Filelists{}
		if err := json.Unmarshal(b, &jout); err != nil && jout.Name == "" {
		} else {
			if jread, err := cfg.sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	case table.COPYFILE:
		jout := table.Copyfile{}
		if err := json.Unmarshal(b, &jout); err != nil {
		} else {
			if jread, err := cfg.sql.Edit(tName, &jout, id); err != nil {
			} else {
				return fmt.Sprintf("%s", jread), nil
			}
		}
	}
	return "", errors.New("Data NG")
}

func (cfg *sqlEdit) sqledit(w http.ResponseWriter, r *http.Request) {

	out := message.Result{Name: "sql edit", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
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
		out.Result = []string{}
	} else if urlPoint+IDPOINT == len(sUrl) || sUrl[urlPoint+IDPOINT] == "" {
		out.Option += "table=" + tName + " id not input"
		out.Code = http.StatusNotFound
		out.Result = []string{}
	} else {
		id, err := strconv.Atoi(sUrl[urlPoint+IDPOINT])
		if err != nil {
			out.Option += "table=" + tName + " id input error"
			out.Code = http.StatusNotFound
			out.Result = []string{}
		} else {
			b, _ := io.ReadAll(r.Body)
			out.Option += "table=" + tName + " id=" + sUrl[urlPoint+IDPOINT]
			if jdata, err := cfg.sqleditwrite(tName, b, id); err != nil {
				out.Code = http.StatusNotFound
				log.Println(err.Error())
			} else {
				message.Println("Edit database for", tName, "id=", id, "data:", string(b))
				joutdata = fmt.Sprintf("%s", jdata)
				out.Result = ""
			}

		}

	}
	common.Sqlreadmessageback(out, joutdata, w)
}

func (cfg *sqlEdit) sqleditget(w http.ResponseWriter, r *http.Request) {
	out := message.Result{Name: "sql edit", Code: http.StatusOK, Date: time.Now(), Option: r.Method + ":"}
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
		out.Result = []string{}
	} else if urlPoint+IDPOINT == len(sUrl) || sUrl[urlPoint+IDPOINT] == "" {
		out.Option += "table=" + tName + " id not input"
		out.Code = http.StatusNotFound
		out.Result = []string{}
	} else {
		id, err := strconv.Atoi(sUrl[urlPoint+IDPOINT])
		if err != nil {
			out.Option += "table=" + tName + " id input error"
			out.Code = http.StatusNotFound
			out.Result = []string{}
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

func (cfg *sqlEdit) websqledit(w http.ResponseWriter, r *http.Request) {

	switch strings.ToUpper(r.Method) {
	case "POST":
		cfg.sqledit(w, r)
	default:
		cfg.sqleditget(w, r)
	}
}

func WebSQLEdit(cfg interface{}, w http.ResponseWriter, r *http.Request) {
	switch cfg.(type) {
	case *table.SQLStatus:
		t := sqlEdit{sql: cfg.(*table.SQLStatus)}
		t.websqledit(w, r)
	}
}
