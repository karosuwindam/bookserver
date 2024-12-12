package indexpage

import (
	"bookserver/config"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var baseurl string

func Init(url string) func(w http.ResponseWriter, r *http.Request) {
	baseurl = url

	return index
}

func index(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len(baseurl):]
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	pass := config.Web.StaticPage
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	tmpUrls := strings.Split(url, "/")
	if url == "" || url == "index.html" || url == "index.htm" {
		filepath := pass + "index.html"
		_, err := os.Stat(filepath)
		if err == nil {
			tmp := make(map[string]string)
			tmp["base_title"] = "新刊取得"
			tmp["maxfilesize"] = config.BScfg.MAX_MULTI_MEMORY
			title := os.Getenv("WEB_TITLE")
			if title != "" {
				tmp["base_title"] = title
			}
			tmp["version"] = config.Version
			tpl := template.Must(template.ParseFiles(filepath))
			tpl.Execute(w, tmp)
			return
		}

	} else if tmpUrls[len(tmpUrls)-1] == "" {
		filepath := pass + url + "index.html"
		_, err := os.Stat(filepath)
		if err == nil {
			tmp := make(map[string]string)
			tmp["base_title"] = "新刊取得"
			title := os.Getenv("WEB_TITLE")
			if title != "" {
				tmp["base_title"] = title
			}
			tpl := template.Must(template.ParseFiles(filepath))
			tpl.Execute(w, tmp)
			return
		}

	} else {
		filepath := pass + url
		_, err := os.Stat(filepath)
		if err == nil {
			if strings.Index(filepath, ".html") >= 0 {
				tmp := make(map[string]string)
				tpl := template.Must(template.ParseFiles(filepath))
				tpl.Execute(w, tmp)
				return

			} else {
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
		}
	}
	slog.WarnContext(ctx, "Notfond Page")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))

}
