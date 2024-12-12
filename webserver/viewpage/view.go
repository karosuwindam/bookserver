package viewpage

import (
	"bookserver/table/historyviews"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

// /view/:idで呼び出される
// もしidの値が数列に変換できない場合は、静的ページのviewフォルダから対象ファイル名を読み取る
func GetIdView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.URL, r.Method),
		"url", r.URL,
		"Method", r.Method,
	)
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		htmlPageView(w, r)
	} else {
		slog.DebugContext(ctx,
			fmt.Sprintf("addhistory id=%v", id),
			"id", id,
		)
		addhistory(id, r)
		filepath := htmlpass + baseurl + "/index.html"
		tmp := make(map[string]string)
		tmp["id"] = tmpid
		tmp["page"] = "1"
		tpl := template.Must(template.ParseFiles(filepath))
		tpl.Execute(w, tmp)

	}

}

func addhistory(id int, r *http.Request) {
	tmp := historyviews.HistoryViews{
		FileId: id,
		Ip:     r.RemoteAddr,
		User:   "guest",
	}
	if err := tmp.Add(); err != nil {
		slog.Warn("addhistory error",
			"error", err,
		)
	}
}

func htmlPageView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "htmlPageView run")

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
	slog.WarnContext(ctx,
		fmt.Sprintf("error opne file name=%v", filepath),
		"error", err,
	)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))

}
