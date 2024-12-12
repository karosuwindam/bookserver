package copyfilepublic

import (
	"bookserver/controller/copyfile"
	"bookserver/webserver/api/common"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type CopyFileRecv struct {
	Id   int  `json:"Id"`
	Flag bool `json:"Flag"`
}

func PostCopyFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, _ := io.ReadAll(r.Body)
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Body", string(data),
	)
	msg := common.Message("NG")
	//処理を後で各
	tmp := CopyFileRecv{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("PostCopyFile json.Unmarshal error"),
			"Error", err,
		)
		w.Write(msg.Json())
		return
	}
	if err := copyfile.Add(tmp.Id, tmp.Flag); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("PostCopyFile copyfile.Add id=%v error", tmp.Id),
			"id", tmp.Id,
			"flag", tmp.Flag,
			"Error", err,
		)
		w.Write(msg.Json())
		return
	}
	msg = common.Message("OK")
	w.Write(msg.Json())
}
