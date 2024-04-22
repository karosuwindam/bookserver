package historyview

import (
	"bookserver/config"
	"bookserver/table/historyviews"
	"bookserver/webserver/api/common"
	"log"
	"net/http"
	"strconv"
)

func GetHistoryRead(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	var n int
	queryParams := r.URL.Query()
	tmp := queryParams.Get("n")
	if c, err := strconv.Atoi(tmp); err != nil {
		n = 10
	} else {
		if c >= config.BScfg.HistoryMax {
			n = 10
		}
	}

	ary, err := historyviews.GetHistory(n)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return

	}
	msg := common.Message(ary)
	w.WriteHeader(http.StatusOK)
	w.Write(msg.Json())
}
