package view

import (
	"bookserver/api/common"
	"bookserver/config"
	"bookserver/textroot/textread"
	"bookserver/webserver"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var Apiname string = "view"

const (
	ROOTPATH = "./html"
)

// 静的HTMLのページを返す
func viewhtml(w http.ResponseWriter, r *http.Request) {
	textdata := []string{".html", ".htm", ".css", ".js"}

	sUrl := common.UrlAnalysis(r.URL.Path)
	var tId string
	var tPage string
	for i, url := range sUrl {
		if url == Apiname {
			if len(sUrl)-i > 1 {
				tId = sUrl[i+1]
				if _, err := strconv.Atoi(tId); err != nil {
					tId = ""
				}
			}
			if len(sUrl)-i > 2 && tId != "" {
				tPage = sUrl[i+2]
				if _, err := strconv.Atoi(tPage); err != nil {
					tPage = ""
				}
			} else if tId != "" {
				tPage = "1"
			}
			break
		}
	}
	if tId != "" && tPage != "" {
		tmp := map[string]string{}
		tmp["id"] = tId
		tmp["page"] = tPage
		upath := "/" + Apiname + "/" + "index.html"
		fmt.Fprint(w, textread.ConvertData(textread.ReadHtml(ROOTPATH+upath), tmp))
		return
	}
	upath := r.URL.Path
	if tId != "" {
		for i, url := range sUrl {
			if url == Apiname {
				upath = "/" + url
				for j := i + 2; j < len(sUrl); j++ {
					upath += "/" + sUrl[j]
				}
				break
			}
		}
	}
	tmp := map[string]string{}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	fmt.Println(r.Method + ":" + r.URL.Path)
	if upath[len(upath)-1:] == "/" {
		upath += "index.html"
	}
	if !textread.Exists(ROOTPATH + upath) {
		w.WriteHeader(404)
		log.Printf("ERROR request:%v\n", r.URL.Path)
		return
	} else {
		for _, data := range textdata {
			if len(upath) > len(data) {
				if upath[len(upath)-len(data):] == data {
					fmt.Fprint(w, textread.ConvertData(textread.ReadHtml(ROOTPATH+upath), tmp))
					return
				}
			}
		}
		buffer := textread.ReadOther(ROOTPATH + upath)
		// bodyに書き込み
		w.Write(buffer)
	}
	return
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + Apiname + "/", Handler: viewhtml},
}

func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
