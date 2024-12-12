package tableread

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"bookserver/webserver/api/common"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

// urlでtableとidをしていて読み取る
func GetReadId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	table := r.PathValue("table")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(ctx, "GetReadId strconv.Atoi error", "error", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if tmp := readTablebyId(id, table); tmp == "" {
		//検索したが失敗したときの処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	} else {
		//検索して成功したときの処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tmp))
	}
	return
}

// urlでtableを指定して読み取る
func GetReadAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	tables := r.PathValue("table")
	if tmp := readTableAll(tables); tmp != "" {
		//検索して成功したときの処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tmp))
		return
	}

	//検索したが失敗したときの処理
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[]"))
	return
}

func readTablebyId(id int, table string) string {
	var output string
	ctx := context.TODO()
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetId(id); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTablebyId booknames.GetId Id=%v error", id),
				"id", id,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTablebyId booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.GetId(id); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTablebyId filelists.GetId Id=%v error", id),
				"id", id,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTablebyId filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.GetId(id); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTablebyId copyfiles.GetId Id=%v error", id),
				"id", id,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTablebyId copyfiles json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("readTablebyId table:%v id:%v", table, id),
		"table", table,
		"id", id,
		"output", output,
	)
	return output
}

func readTableAll(table string) string {
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetAll(); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll booknames.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTableAll booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.GetAll(); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll filelists.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTableAll filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.GetAll(); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll copyfiles.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					"readTableAll copyfiles json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("readTableAll table:%v", table),
		"table", table,
		"output", output,
	)
	return output
}
