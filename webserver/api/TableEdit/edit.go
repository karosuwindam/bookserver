package tableedit

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func PostTableEditdId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := r.PathValue("table")
	data, _ := io.ReadAll(r.Body)
	tmpid := r.PathValue("id")
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
		"Data", string(data),
	)
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if err := editTableSelect(table, id, data); err != nil {
		slog.ErrorContext(ctx,
			fmt.Sprintf("PostTableEditdId editTableSelect table=%v id=%v", table, id),
			"Table", table,
			"Id", id,
			"data", string(data),
			"Error", err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.Write([]byte("{\"message\":\"ok\"}"))

}

func editTableSelect(table string, id int, data []byte) error {
	switch table {
	case BOOKNAMES:
		return editTableBooknames(id, data)
	case COPYFILES:
		return editTableCopyfiles(id, data)
	case FILELISTS:
		return editTableFilelists(id, data)
	}
	return nil
}

func editTableBooknames(id int, data []byte) error {
	tmp := booknames.Booknames{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	tmp.Id = uint(id)
	if err := tmp.Update(); err != nil {
		return err
	}
	return nil
}

func editTableCopyfiles(id int, data []byte) error {
	tmp := copyfiles.Copyfile{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	tmp.Id = uint(id)
	if err := tmp.Update(); err != nil {
		return err
	}
	return nil
}

func editTableFilelists(id int, data []byte) error {
	tmp := filelists.Filelists{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	tmp.Id = uint(id)
	if err := tmp.Update(); err != nil {
		return err
	}
	return nil
}
