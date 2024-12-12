package main

import (
	"bookserver/config"
	"bookserver/controller"
	"bookserver/table"
	"bookserver/webserver"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := table.Init(); err != nil {
		panic(err)
	}
	if err := webserver.Init(); err != nil {
		panic(err)
	}
	if err := controller.Init(); err != nil {
		panic(err)
	}
}

func Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	config.TracerStart(config.TraData.GrpcURL, config.TraData.ServiceName, ctx)
	defer config.TracerStop(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	go func(ctx context.Context) { //コントローラの開始
		defer wg.Done()
		if err := controller.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)
	go func() { //Webサーバの開始処理
		defer wg.Done()
		if err := webserver.Start(ctx); err != nil {
			panic(err)
		}
	}()

	<-sigs
	Stop(context.Background())
	cancel()
	wg.Wait()
}

func Stop(ctx context.Context) {
	slog.InfoContext(ctx, "Main Stop Start")
	webserver.Stop(ctx)
	slog.DebugContext(ctx, "webserver Stop")
	if err := controller.Stop(); err != nil {
		panic(err)
	}
	slog.DebugContext(ctx, "controller Stop")
	slog.InfoContext(ctx, "Main Stop End")
}

func main() {

	Start()
}
