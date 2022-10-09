package main

import (
	"bookserver/config"
	"bookserver/table"
	"encoding/json"
	"fmt"
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
	if rd, err := cfg.sql.ReadAll(table.BOOKNAME); err == nil {
		fmt.Println(rd)
	}
	writedata := table.Booknames{}
	jsond := "{\"Id\":\"1\",\"Name\":\"test\",\"Title\":\"tt\",\"Writer\":\"ttt\",\"Brand\":\"tttt\",\"Booktype\":\"aaaa\",\"Ext\":\"bbb\"}"

	json.Unmarshal([]byte(jsond), &writedata)
	cfg.sql.Add(table.BOOKNAME, &writedata)
	bJson, _ := json.Marshal(&writedata)
	fmt.Println(string(bJson))

	return nil
}

func main() {
	cfg := Setup()
	if cfg == nil {
		return
	}
	Run(cfg)
}
