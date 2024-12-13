package tableread

import (
	"bookserver/config"
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
	ctx, span := config.TracerS(ctx, "GetReadId", "Get Read Id")
	defer span.End()

	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	ctx, err := readRequestTableId(r)
	if err != nil {
		slog.ErrorContext(ctx, "GetReadId readRequestTableId error", "error", err)
	}
	if tmp := readTablebyId(ctx); tmp == "" {
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

func readRequestTableId(r *http.Request) (context.Context, error) {
	ctx := r.Context()
	table := r.PathValue("table")
	tmpid := r.PathValue("id")
	slog.DebugContext(ctx,
		fmt.Sprintf("readRequestTableId table:%v id:%v", table, tmpid),
		"table", table,
		"id", tmpid,
	)

	id, err := strconv.Atoi(tmpid)
	if err != nil {
		return r.Context(), err
	}
	ctx = contextWriteTableIdName(ctx, table, id)
	return ctx, nil
}

// urlでtableを指定して読み取る
func GetReadAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := config.TracerS(ctx, "GetReadId", "Get Read Id")
	defer span.End()

	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	ctx = readRequestTable(r)
	if tmp := readTableAll(ctx); tmp != "" {
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

func readRequestTable(r *http.Request) context.Context {
	ctx := r.Context()
	ctx = contextWriteTableName(ctx, r.PathValue("table"))
	return ctx
}

func readTablebyId(ctx context.Context) string {
	ctx, span := config.TracerS(ctx, "readTablebyId", "Read Table by Id")
	defer span.End()
	var output string
	v, ok := contextReadTableIdName(ctx)
	if !ok {
		msg := "readTablebyId context data not"
		slog.WarnContext(ctx,
			fmt.Sprintf(msg),
		)
		config.TracerError(span, fmt.Errorf(msg))
	}
	table := v.tableName
	id := v.id

	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetId(id); err != nil {
			config.TracerError(span, err)
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
			config.TracerError(span, err)
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
			config.TracerError(span, err)
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

func readTableAll(ctx context.Context) string {
	ctx, span := config.TracerS(ctx, "readTableAll", "Read Table All")
	defer span.End()

	var output string
	table, ok := contextReadTableName(ctx)
	if !ok {
		msg := "readTableAll context data not"
		slog.WarnContext(ctx,
			fmt.Sprintf(msg),
			"table", table,
		)
		config.TracerError(span, fmt.Errorf(msg))
		return output
	}
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetAll(ctx); err != nil {
			config.TracerError(span, err)
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll booknames.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				config.TracerError(span, errj)
				slog.WarnContext(ctx,
					"readTableAll booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.GetAll(ctx); err != nil {
			config.TracerError(span, err)
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll filelists.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				config.TracerError(span, errj)
				slog.WarnContext(ctx,
					"readTableAll filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.GetAll(ctx); err != nil {
			config.TracerError(span, err)
			slog.WarnContext(ctx,
				fmt.Sprintf("readTableAll copyfiles.GetAll error"),
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				config.TracerError(span, errj)
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
