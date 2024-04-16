package tabledelete

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url+"/{table}/{id}", GetChackTalbeById)
	mux.HandleFunc("DELETE "+url+"/{table}/{id}", DeleteTableById)
	return nil
}
