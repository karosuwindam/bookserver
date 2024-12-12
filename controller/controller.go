package controller

import (
	"bookserver/controller/convert"
	"bookserver/controller/copyfile"
	readzipfile "bookserver/controller/readZipfile"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var shutdown chan bool
var shutdown_back chan bool
var run chan bool

// 定期的に処理するプログラムについて
func Init() error {
	shutdown = make(chan bool, 1)
	shutdown_back = make(chan bool, 1)
	run = make(chan bool, 1)
	if err := convert.Init(); err != nil {
		return err
	}
	if err := readzipfile.Init(); err != nil {
		return err
	}
	if err := copyfile.Init(); err != nil {
		return err
	}
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	slog.InfoContext(ctx, "controller Start")

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		readzipfile.Run(ctx)
	}(ctx)
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := copyfile.Run(ctx); err != nil {
			slog.ErrorContext(ctx,
				"copyfile.Run error",
				"Error", err,
			)
		}
	}(ctx)
	wg.Add(1)
	go func() { //定期処理
		defer wg.Done()
	loop:
		for {
			select {
			case <-shutdown:
				shutdown_back <- true
				break loop
			case <-run:
				ctxconvert, _ := context.WithTimeout(ctx, 1*time.Hour)
				if err := convert.CheckCovertData(); err != nil {
					slog.ErrorContext(ctxconvert,
						"CheckCovertData error",
						"Error", err,
					)
				}
			case <-time.After(5 * time.Second):
				ctxconvert, _ := context.WithTimeout(ctx, 5*time.Second)
				if err := convert.CheckCovertData(); err != nil {
					slog.ErrorContext(ctxconvert,
						"CheckCovertData error",
						"Error", err,
					)
				}
			}
		}
	}()
	wg.Wait()
	slog.InfoContext(ctx, "controller Stop")
	return nil
}

func Stop() error {
	ctx := context.TODO()
	slog.DebugContext(ctx, "controller Stop Start")
	shutdown <- true
	if err := readzipfile.Stop(); err != nil {
		return err
	}
	if err := copyfile.Stop(); err != nil {
		return err
	}
	select {
	case <-shutdown_back:
		slog.DebugContext(ctx, "controller Stop End")
		break
	case <-time.After(5 * time.Second):
		return errors.New("Shutdown error")
	}
	return nil
}

func Start() {
	if len(run) > 0 {
		run <- true
	}
}
