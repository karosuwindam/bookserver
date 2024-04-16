package copyfilepublic

import "net/http"

func Init(url string, mux *http.ServeMux) error {
	mux.HandleFunc("GET "+url+"/{id}", GetFileDataFlagById)
	mux.HandleFunc("POST "+url, PostCopyFile)
	return nil
}
