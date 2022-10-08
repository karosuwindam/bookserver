package main

import (
	"bookserver/config"
	"bookserver/table"
)

type mainconfig struct {
	sql *table.SQLStatus
}

func Setup() *mainconfig {
	cfg, err := config.EnvRead()
	if err != nil {
		return nil
	}
	sql, err := table.Setup(cfg)
	if err != nil {
		return nil
	}
	return &mainconfig{sql: sql}
}

func Run(cfg *mainconfig) error {
	defer cfg.sql.Close()

	return nil
}

func main() {
	cfg := Setup()
	if cfg == nil {
		return
	}
	Run(cfg)
}
