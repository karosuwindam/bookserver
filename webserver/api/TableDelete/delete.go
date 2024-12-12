package tabledelete

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func GetChackTalbeById(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if err := checkTablebyId(id, table); err != nil {
		slog.WarnContext(ctx,
			fmt.Sprintf("GetChackTalbeById checkTablebyId table=%v id=%v", table, id),
			"Table", table,
			"Id", id,
			"Error", err,
		)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\":\"ng\"}"))
		return
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("GetChackTalbeById checkTablebyId table=%v id=%v ok", table, id),
		"Table", table,
		"Id", id,
	)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\":\"ok\"}"))
}

// DELETEで受け取ったurlをもとに削除する
func DeleteTableById(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if err := deleteTablebyId(id, table); err != nil {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\":\"ng\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\":\"ok\"}"))
}

func deleteTablebyId(id int, table string) error {
	switch table {
	case BOOKNAMES:
		return booknames.Delete(id)
	case FILELISTS:
		return filelists.Delete(id)
	case COPYFILES:
		return copyfiles.Delete(id)
	}
	return nil
}

func checkTablebyId(id int, table string) error {
	switch table {
	case BOOKNAMES:
		_, err := booknames.GetId(id)
		return err
	case FILELISTS:
		_, err := filelists.GetId(id)
		return err
	case COPYFILES:
		_, err := copyfiles.GetId(id)
		return err
	}
	return errors.New("not table")
}
