package copyfilepublic

import (
	"bookserver/config"
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {

	config.TraceHttpHandleFunc(mux, "GET "+url+"/{id}", GetFileDataFlagById)
	config.TraceHttpHandleFunc(mux, "POST "+url, PostCopyFile)
	return nil
}
