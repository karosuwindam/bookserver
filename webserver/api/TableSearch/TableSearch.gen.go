package tablesearch

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "POST "+url, PostSerchTable)
	config.TraceHttpHandleFunc(mux, "GET "+url+"/{table}/{keyword}", GetSearchTable)
	return nil
}
