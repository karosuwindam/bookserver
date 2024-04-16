package tablesearch

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("POST "+url, PostSerchTable)
	mux.HandleFunc("GET "+url+"/{table}/{keyword}", GetSearchTable)
	return nil
}
