package ziipfileview

import (
	readzipfile "bookserver/controller/readZipfile"
	"bookserver/webserver/api/common"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func GetZipFileList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(ctx, "GetZipFileList Atoi error",
			"tmpid", tmpid,
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	//idを指定してfilelistテーブルからzipファイル名を取得
	d, err := readzipfile.GetZiplist(readzipfile.ContextWriteZipId(ctx, id))
	if err != nil {
		slog.ErrorContext(ctx, "GetZipFileList readzipfile.GetZiplist error",
			"id", id,
			"error", err,
		)
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
