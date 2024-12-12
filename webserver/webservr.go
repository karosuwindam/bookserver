package webserver

import (
	"bookserver/config"
	"bookserver/webserver/api"
	healthcheck "bookserver/webserver/healthCheck"
	"bookserver/webserver/indexpage"
	"bookserver/webserver/viewpage"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// SetupServer
// サーバ動作の設定
type SetupServer struct {
	protocol string // Webサーバーのプロトコル
	hostname string //Webサーバのホスト名
	port     string //Webサーバの解放ポート

	mux *http.ServeMux //webサーバのmux
}

var srv *http.Server // Webサーバの管理関数

var cfg SetupServer

func HelloWeb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Web"))
}

func Init() error {
	cfg = SetupServer{
		protocol: config.Web.Protocol,
		hostname: config.Web.Hostname,
		port:     config.Web.Port,
		mux:      http.NewServeMux(),
	}
	api.Init(cfg.mux)
	if err := healthcheck.Init(cfg.mux); err != nil {
		return errors.Wrap(err, "healthcheck.Init()")
	}
	if err := viewpage.Init("/view", cfg.mux); err != nil {
		return errors.Wrap(err, "viewpage.Init()")
	}
	// cfg.mux.HandleFunc("/", indexpage.Init("/"))
	config.TraceHttpHandleFunc(cfg.mux, "/", indexpage.Init("/"))
	return nil
}

func Start(ctx context.Context) error {
	var err error = nil
	var wg sync.WaitGroup

	// if config.TraData.TracerUse {
	// 	hander := otelhttp.NewHandler(cfg.mux, "http-server",
	// 		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	// 	)
	// 	srv = &http.Server{
	// 		Addr:         cfg.hostname + ":" + cfg.port,
	// 		Handler:      hander,
	// 		ReadTimeout:  30 * time.Second,
	// 		WriteTimeout: 60 * time.Second,
	// 	}
	// } else {
	srv = &http.Server{
		Addr:         cfg.hostname + ":" + cfg.port,
		Handler:      cfg.mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	// }
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx,
		fmt.Sprintf("Start Server %v:%v", cfg.hostname, cfg.port),
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
			panic(err)
		} else {
			err = nil
		}
	}()
	wg.Wait()
	slog.InfoContext(ctx, "Server Stop")
	return err
}

func Stop(ctx context.Context) error {
	if srv == nil {
		return nil
	}
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
