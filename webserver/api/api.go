package api

import (
	copyfilepublic "bookserver/webserver/api/CopyFIlePublic"
	fileupload "bookserver/webserver/api/FileUpload"
	tableadd "bookserver/webserver/api/TableAdd"
	tabledelete "bookserver/webserver/api/TableDelete"
	tableedit "bookserver/webserver/api/TableEdit"
	tableread "bookserver/webserver/api/TableRead"
	tablesearch "bookserver/webserver/api/TableSearch"
	ziipfileview "bookserver/webserver/api/ZiipfileView"
	zipfileimageview "bookserver/webserver/api/ZipfileImageView"
	"net/http"
)

type api struct {
	Router string
	Func   func(string, *http.ServeMux) error
}

var v1apis = []api{
	{"/search", tablesearch.Init},
	{"/add", tableadd.Init},
	{"/read", tableread.Init},
	{"/delete", tabledelete.Init},
	{"/upload", fileupload.Init},
	{"/list", ziipfileview.Init},
	{"/edit", tableedit.Init},
	{"/image", zipfileimageview.Init},
	{"/copy", copyfilepublic.Init},
}

func Init(mux *http.ServeMux) error {
	if err := v1apisetup(mux, "/v1"); err != nil {
		return err
	}
	return nil
}

func v1apisetup(mux *http.ServeMux, router string) error {
	if router == "/" {
		router = ""
	}
	if router[len(router)-1] == '/' {
		router = router[:len(router)-1]

	}
	for _, v := range v1apis {
		if err := v.Func(router+v.Router, mux); err != nil {
			return err
		}
	}
	return nil
}
