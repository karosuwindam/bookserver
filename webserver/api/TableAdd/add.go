package tableadd

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// POSTで受け取ったJSONデータをベースにテーブルへデータ登録をする。
func PostAddTable(w http.ResponseWriter, r *http.Request) {
	table := r.PathValue("table")
	data, _ := io.ReadAll(r.Body)
	log.Println("info:", r.URL, r.Method, string(data))
	if checkTableData(table) != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if err := addTableSelect(table, data); err != nil {
		log.Println("error:", err)
		return
	}
	w.Write([]byte("{\"message\":\"ok\"}"))
}

func addTableSelect(table string, data []byte) error {
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
