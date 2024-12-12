package historyview

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "GET "+url, GetHistoryRead)
	return nil
}
