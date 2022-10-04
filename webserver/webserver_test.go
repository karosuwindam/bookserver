package webserver

import (
	"bookserver/config"
	"bookserver/message"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

type hellov1s struct {
	str string
}

func hellov1(t interface{}, w http.ResponseWriter, r *http.Request) {
	switch t.(type) {
	case *hellov1s:
		out := t.(*hellov1s)
		fmt.Fprintf(w, "hello %v", out.str)
	default:
		hello(w, r)
	}
}

func TestWebServerSetup(t *testing.T) {
	t.Log("----------------- Server Setup --------------------------")
	cfg, _ := config.EnvRead()
	s, err := NewSetup(cfg)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	t.Log(s.protocol, s.hostname, s.port)
	t.Log("----------------- Server Setup OK --------------------------")
}

func TestWebServerNew(t *testing.T) {
	t.Log("----------------- Server New --------------------------")
	cfg, _ := config.EnvRead()
	ss, _ := NewSetup(cfg)
	ss.Add("/", hello)
	h := &hellov1s{str: "test"}
	if err := ss.AddV1("test", h, hellov1); err != nil {
		t.Error(err)
		t.FailNow()
	}
	s, err := ss.NewServer()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	ctx, cancel := context.WithCancel(context.Background())
	eq, ctx := errgroup.WithContext(ctx)
	eq.Go(func() error {
		return runServer(ctx, s)
	})
	rsp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Errorf("%+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll((rsp.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got) != "hello" {
		t.Errorf("message : %v", string(got))
	}
	rsp1, err := http.Get("http://localhost:8080/v1/test")
	if err != nil {
		t.Errorf("%+v", err)
	}
	defer rsp1.Body.Close()
	got1, err := io.ReadAll((rsp1.Body))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if string(got1) != "hello "+h.str {
		t.Errorf("message : %v", string(got1))
	}
	cancel()
	if err := eq.Wait(); err != nil {
		t.Fatal(err)
	}
	t.Log("----------------- Server New OK --------------------------")

}

func runServer(ctx context.Context, s *Server) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Srv.Serve(s.L); err != nil &&
			err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	<-ctx.Done()
	message.Println("server stop")
	s.Srv.Shutdown(context.Background())
	return eg.Wait()
}
