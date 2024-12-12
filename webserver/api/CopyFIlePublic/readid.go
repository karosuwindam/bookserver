package copyfilepublic

import (
	"bookserver/controller/copyfile"
	"bookserver/webserver/api/common"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func GetFileDataFlagById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetFileDataFlagById strconv.Atoi id=%v", tmpid),
			"Id", tmpid,
			"Error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	// idに紐づいたファイルが存在するか確認
	if d, err := copyfile.ReadCopyFIleFlagById(id); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetFileDataFlagById ReadCopyFIleFlagById id=%v", id),
			"Id", id,
			"Error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		// 確認した結果JSON形式でを返す。
		msg := common.Message(d)
		w.Write(msg.Json())
	}

}
