package webserver

import (
	"bookserver/config"
	"bookserver/message"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type SetupServer struct {
	protocol string
	hostname string
	port     string

	routefunc      map[string]func(interface{}, http.ResponseWriter, *http.Request)
	routeinterface map[string]interface{}

	mux *http.ServeMux
}

type Server struct {
	Srv *http.Server
	L   net.Listener
}

type Status struct {
	Status string `json:status`
}

func NewSetup(data *config.Config) (*SetupServer, error) {
	cfg := &SetupServer{
		protocol: data.Server.Protocol,
		hostname: data.Server.Hostname,
		port:     data.Server.Port,
	}
	cfg.mux = http.NewServeMux()
	cfg.routefunc = map[string]func(interface{}, http.ResponseWriter, *http.Request){}
	cfg.routeinterface = map[string]interface{}{}
	cfg.mux.HandleFunc("/v1/", cfg.v1)
	return cfg, nil
}

func (t *SetupServer) NewServer() (*Server, error) {
	message.Println("Setupserver", t.protocol, t.hostname+":"+t.port)
	l, err := net.Listen(t.protocol, t.hostname+":"+t.port)
	if err != nil {
		return nil, err
	}
	return &Server{
		Srv: &http.Server{Handler: t.muxHandler()},
		L:   l,
	}, nil
}

//未チェック
func (t *SetupServer) AddV1(route string, sdata interface{}, funcdata func(interface{}, http.ResponseWriter, *http.Request)) error {
	var v1route string
	if route[0] != "/"[0] {
		v1route = "/v1" + "/" + route
	} else {
		v1route = "/v1" + route
	}

	if t.routefunc[v1route] != nil {
		return errors.New("")
	}
	t.routefunc[v1route] = funcdata
	t.routeinterface[v1route] = sdata
	return nil
}

func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) {
	t.mux.HandleFunc(route, handler)
}

func (t *SetupServer) muxHandler() http.Handler { return t.mux }

func v1OtherBuck(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	msg := message.Message{Name: "v1", Status: url, Code: http.StatusOK}
	rst := message.Result{Name: "v1", Code: http.StatusOK, Option: "", Date: time.Now(), Result: msg}
	fmt.Fprintf(w, "%v", rst.Output())
}

func (t *SetupServer) v1(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	flag := false
	var check string
	for turl, _ := range t.routeinterface {
		if len(url) >= len(turl) {
			if turl == url[:len(turl)] {
				check = turl
				flag = true
				break
			}
		}
	}
	if flag { //登録済み
		t.routefunc[check](t.routeinterface[check], w, r)
	} else { //日登録
		v1OtherBuck(w, r)
	}

}
