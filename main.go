package main

import (
	"bookserver/config"
	"bookserver/fileupload"
	"bookserver/message"
	"bookserver/textread"
	"bookserver/webserver"
	"bookserver/webserver/common"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
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
		scfg.s.Add("/", textread.ViewHtml)
		scfg.s.AddV1(common.ADMIN, "/upload", upcfg, fileupload.FIleupload)

	}

	return scfg
}

func Run(cfg *mainconfig, ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	s, err := cfg.s.NewServer()
	if err != nil {
		return err
	}
	eg.Go(func() error {
		defer s.Sql.Close()
		errout := s.Srv.Serve(s.L)
		return errout
	})
	<-ctx.Done()
	message.Println("shutdown")
	if err := s.Srv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	return eg.Wait()
}

func main() {
	ctx := context.Background()
	cfg := Setup()

	if err := Run(cfg, ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
