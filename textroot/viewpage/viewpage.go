package viewpage

import (
	"bookserver/config"
	"bookserver/textroot/textread"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	ROOTPATH = "./html"
)

// 静的HTMLのページを返す
func Viewhtml(w http.ResponseWriter, r *http.Request) {
	textdata := []string{".html", ".htm", ".css", ".js"}
	upath := r.URL.Path
	tmp := map[string]string{}
	tmp["version"] = version
	tmp["uploadsize"] = uploadSize
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

var version string = ""    //バージョン表示
var uploadSize string = "" //アップロードサイズ

// Setup(cfg) = error
//
// /配下の静的ページを返す設定
//
// cfg : 基本設定
func Setup(cfg *config.Config) error {
	version = cfg.Version
	uploadSize = cfg.Upload.MAX_MULTI_MEMORY
	return nil
}
