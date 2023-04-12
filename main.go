package main

import (
	"bookserver/api"
	"bookserver/config"
	"bookserver/textroot"
	"bookserver/transform"
	"bookserver/webserverv2"
	"context"
	"fmt"
	"log"
	"os"
)

func Config(cfg *config.Config) (*webserverv2.SetupServer, error) {
	api.Setup(cfg)
	transform.Setup(cfg)
	scfg, err := webserverv2.NewSetup(cfg)
	if err != nil {
		return nil, err
	}
	webserverv2.Config(scfg, api.Route, "/v1")
	webserverv2.Config(scfg, textroot.Route, "")
	return scfg, nil
}
func Run(ctx context.Context) error {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	cfg, err := config.EnvRead()
	if err != nil {
		return nil
	}
	if scfg, err := Config(cfg); err == nil {
		s, err := scfg.NewServer()
		if err == nil {

			go transform.Run(ctx)
			err = s.Run(ctx)
			cancel()
			transform.Wait()
			return err
		}
	}
	return err
}

func EndCK() {
}
func main() {
	log.SetFlags(log.Llongfile | log.Flags())
	ctx := context.Background()
	fmt.Println("start")
	if err := Run(ctx); err != nil {
		EndCK()
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("end")
	EndCK()

}
