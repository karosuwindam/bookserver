package tableadd

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// POSTで受け取ったJSONデータをベースにテーブルへデータ登録をする。
func PostAddTable(w http.ResponseWriter, r *http.Request) {
	table := r.PathValue("table")
	data, _ := io.ReadAll(r.Body)

	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
		"Data", string(data),
	)
	if checkTableData(table) != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if err := addTableSelect(table, data); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("PostAddTable addTableSelect table=%v", table),
			"Table", table,
			"Data", string(data),
			"Error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.Write([]byte("{\"message\":\"ok\"}"))
}

func addTableSelect(table string, data []byte) error {
	ctx := context.TODO()
	slog.DebugContext(ctx, "addTableSelect Run", "Table", table, "Data", string(data))
	switch table {
	case BOOKNAMES:
		return addTableBooknames(data)
	case COPYFILES:
		return addTableCopyfiles(data)
	case FILELISTS:
		return addTableFilelists(data)
	}
	return nil
}

func addTableBooknames(data []byte) error {
	tmp := booknames.Booknames{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := tmp.Add(); err != nil {
		return err
	}
	return nil
}

func addTableCopyfiles(data []byte) error {
	tmp := copyfiles.Copyfile{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := tmp.Add(); err != nil {
		return err
	}
	return nil
}

func addTableFilelists(data []byte) error {
	tmp := filelists.Filelists{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := tmp.Add(); err != nil {
		return err
	}
	return nil
}
