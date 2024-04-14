package healthcheck

import "net/http"

func Init(mux *http.ServeMux) error {
	mux.HandleFunc("/health", Health)
	return nil
}
