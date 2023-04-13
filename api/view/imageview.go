package view

import (
	"bookserver/api/common"
	"bookserver/table"
	"fmt"
	"net/http"
	"strconv"
)

// webzipread(w, r)
//
// /image/:id/:nameの動作
func webzipreadimage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	sUrl := common.UrlAnalysis(r.URL.Path)
	tName := ""
	tId := ""
	for i, url := range sUrl {
		if url == apinameimage && len(sUrl) > i+2 {
			tId = sUrl[i+1]
			tName = sUrl[i+2]
			break
		}
	}
	if tId == "" || tName == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		if id, err := strconv.Atoi(tId); err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			if json, err := sql.ReadID(table.FILELIST, id); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				if data, ok := table.JsonToStruct(table.FILELIST, []byte(json)).([]table.Filelists); ok {
					zipname := data[0].Zippass
					f := openfile(zipname)
					if buf, err := f.openZipRead(tName); err != nil {
						w.WriteHeader(http.StatusNotFound)

					} else {
						w.WriteHeader(http.StatusOK)
						fmt.Fprintf(w, "%s", buf)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		}
	}
}
