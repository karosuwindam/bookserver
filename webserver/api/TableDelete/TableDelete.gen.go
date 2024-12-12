package tabledelete

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "GET "+url+"/{table}/{id}", GetChackTalbeById)
	config.TraceHttpHandleFunc(mux, "DELETE "+url+"/{table}/{id}", DeleteTableById)
	return nil
}
