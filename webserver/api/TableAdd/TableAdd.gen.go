package tableadd

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	config.TraceHttpHandleFunc(mux, "POST "+url+"/{table}", PostAddTable)
	return nil
}
