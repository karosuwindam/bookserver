package webserver

import (
	"bookserver/config"
	"bookserver/webserver/api"
	healthcheck "bookserver/webserver/healthCheck"
	"bookserver/webserver/indexpage"
	"bookserver/webserver/viewpage"
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	cfg.mux.HandleFunc("/", indexpage.Init("/"))
	return nil
}

func Start(ctx context.Context) error {
	var err error = nil
	var wg sync.WaitGroup
	if config.TraData.TracerUse {
		hander := otelhttp.NewHandler(cfg.mux, "http-server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		)

		// if config.NewRelic.App != nil {
		// 	_, hander = newrelic.WrapHandle(config.NewRelic.App, "http-server", hander)
		// }
		srv = &http.Server{
			Addr:         cfg.hostname + ":" + cfg.port,
			Handler:      hander,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 60 * time.Second,
		}
	} else {

		// if config.NewRelic.App != nil {
		// 	_, hander := newrelic.WrapHandle(config.NewRelic.App, "http-server", cfg.mux)

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
	}
	if config.NewRelic.App != nil {
		_, srv.Handler = newrelic.WrapHandle(config.NewRelic.App, "http-server", srv.Handler)
	}
	l, err := net.Listen(cfg.protocol, srv.Addr)
	if err != nil {
		return err
	}
	log.Println("info: Start Server", cfg.hostname+":"+cfg.port)
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
	log.Println("info: Server Stop")
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
