package ziipfileview

import (
	readzipfile "bookserver/controller/readZipfile"
	"bookserver/webserver/api/common"
	"log/slog"
	"net/http"
	"strconv"
)

func GetZipFileList(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "", "URL", r.URL, "Method", r.Method)
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(r.Context(), "GetZipFileList", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	//idを指定してfilelistテーブルからzipファイル名を取得
	d, err := readzipfile.GetZiplist(id)
	if err != nil {
		slog.ErrorContext(r.Context(), "GetZipFileList", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}

	//jsonファイルに変換してメッセージを返す
	msg := common.Message(d)
	if d := msg.Json(); d == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(d)

	}

}
