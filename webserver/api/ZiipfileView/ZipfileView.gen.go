package ziipfileview

import (
	"net/http"
)

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url+"/{id}", GetZipFileList)
	return nil
}
