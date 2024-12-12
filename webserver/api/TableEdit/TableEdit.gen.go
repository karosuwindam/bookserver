package tableedit

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {

	config.TraceHttpHandleFunc(mux, "GET "+url+"/{table}/{id}", GetReadId)
	config.TraceHttpHandleFunc(mux, "POST "+url+"/{table}/{id}", PostTableEditdId)
	config.TraceHttpHandleFunc(mux, "DELETE "+url+"/{table}/{id}", DeleteTableById)

	return nil
}
