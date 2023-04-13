package webserverv2

import (
	"bookserver/config"
	"context"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServerRun(t *testing.T) {
	t.Log("----------------- Server Setup --------------------------")
	cfg, _ := config.EnvRead()
	s, _ := NewSetup(cfg)
	s.Add("/", hello)
	server, _ := s.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx)
	})
	rsp1, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("%+v", err)
	}

	got1, err := io.ReadAll((rsp1.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got1) != sampleword {
		t.Errorf("message : %v", string(got1))
	}
	defer rsp1.Body.Close()
	cancel()
	t.Log("----------------- Server Setup OK --------------------------")

}

func TestServerSetupRun(t *testing.T) {
	t.Log("----------------- Server Setup --------------------------")
	cfg, _ := config.EnvRead()
	s, _ := NewSetup(cfg)
	c := []WebConfig{
		{Pass: "/", Handler: hello},
	}
	if err := Config(s, c, "/"); err != nil {
		t.Errorf("%+v", err)
	}
	if err := Config(s, c, "/v1"); err != nil {
		t.Errorf("%+v", err)
	}
	server, _ := s.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx)
	})
	t.Log("----------------- Get / data --------------------------")
	rsp1, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("%+v", err)
	}
	got1, err := io.ReadAll((rsp1.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got1) != sampleword {
		t.Errorf("message : %v", string(got1))
	}
	defer rsp1.Body.Close()
	t.Log("----------------- Get / data OK --------------------------")
	t.Log("----------------- Get /v1/ data --------------------------")
	rsp2, err := http.Get("http://localhost:8080/v1/")
	if err != nil {
		t.Errorf("%+v", err)
	}
	got2, err := io.ReadAll((rsp2.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got2) != sampleword {
		t.Errorf("message : %v", string(got2))
	}
	defer rsp2.Body.Close()
	t.Log("----------------- Get / data OK --------------------------")
	cancel()
	t.Log("----------------- Server Setup OK --------------------------")

}
