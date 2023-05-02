package proffdebug

import (
	"bookserver/config"
	"bookserver/webserver"
	"fmt"
	"net/http"
	"net/http/pprof"
)

var Route []webserver.WebConfig = []webserver.WebConfig{}

func Setup(cfg *config.Config) error {
	tmpRoute := []webserver.WebConfig{
		{Pass: "/", Handler: basepageview},
		{Pass: "/pprof/", Handler: pprof.Index},
		{Pass: "/pprof/cmdline", Handler: pprof.Cmdline},
		{Pass: "/pprof/profile", Handler: pprof.Profile},
		{Pass: "/pprof/symbol", Handler: pprof.Symbol},
		{Pass: "/pprof/trace", Handler: pprof.Trace},
	}
	Route = append(Route, tmpRoute...)
	return nil
}

func basepageview(w http.ResponseWriter, r *http.Request) {
	str := ""
	str = "<html><head><title>debug</title></head>"
	str += "<body>"
	str += "<div><a href='pprof'>pprof</a></div>"
	str += "</body></html>"
	fmt.Fprintln(w, str)
}
