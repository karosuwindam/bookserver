package controller

import (
	"bookserver/controller/convert"
	"bookserver/controller/copyfile"
	readzipfile "bookserver/controller/readZipfile"
	"context"
	"errors"
	"log"
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
	log.Println("info:", "controller Start")

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		readzipfile.Run(ctx)
	}(ctx)
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := copyfile.Run(ctx); err != nil {
			log.Println("error", err)
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
				if err := convert.CheckCovertData(); err != nil {
					log.Println("error:", err)
				}
			case <-time.After(5 * time.Second):
				if err := convert.CheckCovertData(); err != nil {
					log.Println("error:", err)
				}
			}
		}
	}()
	wg.Wait()
	log.Println("info:", "controller ShutDown")
	return nil
}

func Stop() error {
	shutdown <- true
	if err := readzipfile.Stop(); err != nil {
		return err
	}
	if err := copyfile.Stop(); err != nil {
		return err
	}
	select {
	case <-shutdown_back:
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
