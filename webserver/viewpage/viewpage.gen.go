package viewpage

import (
	"bookserver/config"
	"net/http"
)

var baseurl string
var htmlpass string

func Init(url string, mux *http.ServeMux) error {
	baseurl = url
	if baseurl[len(baseurl)-1:] == "/" {
		baseurl = baseurl[:len(baseurl)-1]
	}

	htmlpass = config.Web.StaticPage
	if htmlpass[len(htmlpass)-1:] == "/" {
		htmlpass = htmlpass[:len(htmlpass)-1]
	}
	mux.HandleFunc("GET "+baseurl+"/{id}", GetIdView)
	return nil
}
