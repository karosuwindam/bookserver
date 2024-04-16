package viewpage

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

// /view/:idで呼び出される
// もしidの値が数列に変換できない場合は、静的ページのviewフォルダから対象ファイル名を読み取る

func GetIdView(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	tmpid := r.PathValue("id")
	_, err := strconv.Atoi(tmpid)
	if err != nil {
		htmlPageView(w, r)
	} else {
		filepath := htmlpass + baseurl + "/index.html"
		tmp := make(map[string]string)
		tmp["id"] = tmpid
		tmp["page"] = "1"
		tpl := template.Must(template.ParseFiles(filepath))
		tpl.Execute(w, tmp)

	}

}

func htmlPageView(w http.ResponseWriter, r *http.Request) {

	pass := htmlpass
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	url := r.URL.Path

	filepath := pass + url
	_, err := os.Stat(filepath)
	if err == nil {
		fp, _ := os.Open(filepath)
		defer fp.Close()
		buf := make([]byte, 1024)
		var buffer []byte
		for {
			n, err := fp.Read(buf)
			if err != nil {
				break
			}
			if n == 0 {
				break
			}
			buffer = append(buffer, buf[:n]...)
		}
		w.Write(buffer)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))

}
