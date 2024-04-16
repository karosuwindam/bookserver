package tableread

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url+"/{table}", GetReadAll)
	mux.HandleFunc("GET "+url+"/{table}/{id}", GetReadId)
	return nil
}
