package tableread

import (
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"bookserver/webserver/api/common"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// urlでtableとidをしていて読み取る
func GetReadId(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	table := r.PathValue("table")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		log.Println("error:", err)
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
	log.Println("info:", r.URL, r.Method)
	tables := r.PathValue("table")
	if tmp := readTableAll(tables); tmp == "" {
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

func readTablebyId(id int, table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetId(id); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.GetId(id); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.GetId(id); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

func readTableAll(table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.GetAll(); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.GetAll(); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.GetAll(); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}
