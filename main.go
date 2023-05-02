package main

import (
	"bookserver/api"
	"bookserver/api/listdata"
	"bookserver/api/view"
	"bookserver/config"
	"bookserver/health"
	"bookserver/proffdebug"
	"bookserver/publiccopy"
	"bookserver/pyroscopesetup"
	"bookserver/textroot"
	"bookserver/transform"
	"bookserver/webserver"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Confg(cfg) = (*webserver.SetupServer, error)
//
// メイン処理の設定
func Config(cfg *config.Config) (*webserver.SetupServer, error) {
	if err := publiccopy.Setup(cfg); err != nil {
		return nil, err
	}
	api.Setup(cfg)
	transform.Setup(cfg)
	proffdebug.Setup(cfg)
	textroot.Setup(cfg)
	scfg, err := webserver.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	webserver.Config(scfg, api.Route, "/v1")
	webserver.Config(scfg, proffdebug.Route, "/debug")
	if route, err := health.SetUp(cfg); err == nil {
		webserver.Config(scfg, route, "")
	}
	webserver.Config(scfg, textroot.Route, "")
	return scfg, nil
}

// Run(ctx) = error
//
// メイン処理
func Run(ctx context.Context) error {
	var err error
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	cfg, err := config.EnvRead()
	if err != nil {
		return nil
	}
	if scfg, err := Config(cfg); err == nil {
		s, err := scfg.NewServer()
		if err == nil {
			go listdata.Loop(ctx)
			go view.Loop(ctx)
			go transform.Run(ctx)
			go publiccopy.Loop(ctx)
			if err = s.Run(ctx); err != nil {
				log.Println(err)
			}
			stop()
			return EndCK()
		}
	}
	return err
}

// EndCK() = error
//
// 終了時の処理
func EndCK() error {
	var err error = nil
	if err1 := publiccopy.Wait(); err1 != nil {
		err = err1
	}
	if err1 := transform.Wait(); err1 != nil {
		err = err1
	}
	if err1 := listdata.Wait(); err1 != nil {
		err = err1
	}
	return err
}

// main()
//
// 　メイン処理
func main() {
	// flag.Parse() //コマンドラインオプションの有効
	log.SetFlags(log.Llongfile | log.Flags())
	pyro := pyroscopesetup.Setup()
	pyroscopesetup.Add("status", "debug")
	pyro.Run()
	ctx := context.Background()
	fmt.Println("start")
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")

}
