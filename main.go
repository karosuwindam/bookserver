package main

import (
	"bookserver/config"
	"bookserver/fileupload"
	"bookserver/webserver"
	"bookserver/webserver/common"
)

type mainconfig struct {
	s *webserver.SetupServer
}

func Setup() *mainconfig {
	scfg := &mainconfig{}
	cfg, err := config.EnvRead()
	if err != nil {
		return nil
	}
	if ss, err := webserver.NewSetup(cfg); err == nil {
		scfg.s = ss
	} else {
		return nil
	}

	if upcfg, err := fileupload.Setup(); err == nil {
		scfg.s.AddV1(common.ADMIN, "/upload", upcfg, fileupload.FIleupload)

	}

	return scfg
}

func Run(cfg *mainconfig) error {
	s, err := cfg.s.NewServer()
	if err != nil {
		return err
	}
	defer s.Sql.Close()
	return s.Srv.Serve(s.L)
}

func main() {
	cfg := Setup()
	if cfg == nil {
		return
	}
	Run(cfg)
}
