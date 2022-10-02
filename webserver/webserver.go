package webserver

import (
	"bookserver/config"
	"bookserver/message"
	"net"
	"net/http"
)

type SetupServer struct {
	protocol string
	hostname string
	port     string

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

func (t *SetupServer) Add(route string, handler func(http.ResponseWriter, *http.Request)) {
	t.mux.HandleFunc(route, handler)
}

func (t *SetupServer) muxHandler() http.Handler { return t.mux }
