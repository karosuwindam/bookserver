package main

import (
	"bookserver/api"
	"bookserver/config"
	"bookserver/proffdebug"
	"bookserver/pyroscopesetup"
	"bookserver/textroot"
	"bookserver/transform"
	"bookserver/webserverv2"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Config(cfg *config.Config) (*webserverv2.SetupServer, error) {
	api.Setup(cfg)
	transform.Setup(cfg)
	proffdebug.Setup(cfg)
	textroot.Setup(cfg)
	scfg, err := webserverv2.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	webserverv2.Config(scfg, api.Route, "/v1")
	webserverv2.Config(scfg, proffdebug.Route, "/debug")
	webserverv2.Config(scfg, textroot.Route, "")
	return scfg, nil
}
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

			go transform.Run(ctx)
			err = s.Run(ctx)
			stop()
			return err
		}
	}
	return err
}

func EndCK() {
	transform.Wait()
}
func main() {
	// flag.Parse() //コマンドラインオプションの有効
	log.SetFlags(log.Llongfile | log.Flags())
	pyro := pyroscopesetup.Setup()
	pyroscopesetup.Add("status", "debug")
	pyro.Run()
	ctx := context.Background()
	fmt.Println("start")
	if err := Run(ctx); err != nil {
		EndCK()
		log.Println(err)
		os.Exit(1)
	}
	EndCK()
	fmt.Println("end")

}
