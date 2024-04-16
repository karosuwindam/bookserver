package tableedit

import "net/http"

func Init(url string, mux *http.ServeMux) error {

	mux.HandleFunc("GET "+url+"/{table}/{id}", GetReadId)
	mux.HandleFunc("POST "+url+"/{table}/{id}", PostTableEditdId)
	mux.HandleFunc("DELETE "+url+"/{table}/{id}", DeleteTableById)
	return nil
}
