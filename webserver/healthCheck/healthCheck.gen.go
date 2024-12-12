package healthcheck

import (
	"bookserver/config"
	"net/http"
)

func Init(mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "/health", Health)
	return nil
}
