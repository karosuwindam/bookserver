package tableadd

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("POST "+url+"/{table}", PostAddTable)
	return nil
}
