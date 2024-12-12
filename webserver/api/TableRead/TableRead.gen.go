package tableread

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "GET "+url+"/{table}", GetReadAll)
	config.TraceHttpHandleFunc(mux, "GET "+url+"/{table}/{id}", GetReadId)
	return nil
}
