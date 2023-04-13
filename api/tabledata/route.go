package tabledata

import (
	"bookserver/api/tabledata/add"
	"bookserver/api/tabledata/edit"
	"bookserver/api/tabledata/read"
	"bookserver/api/tabledata/search"
	"bookserver/config"
	"bookserver/table"
	"bookserver/webserverv2"
)

// routeのベースフォルダ
var route []webserverv2.WebConfig = []webserverv2.WebConfig{}

func sqlSetup(cfg *config.Config) (*table.SQLStatus, error) {
	var err error
	if sqlcfg, err := table.Setup(cfg); err == nil {
		return sqlcfg, err
	}
	return nil, err
}

// Setup() = []webserverv2.WebConfig
//
// セットアップして、HTMLのルートフォルダを用意する
func Setup(cfg *config.Config) ([]webserverv2.WebConfig, error) {
	var err error
	if sqlcfg, err := sqlSetup(cfg); err == nil {
		route = append(route, add.Setup(sqlcfg)...)
		route = append(route, edit.Setup(sqlcfg)...)
		route = append(route, read.Setup(sqlcfg)...)
		route = append(route, search.Setup(sqlcfg)...)
	}
	return route, err
}
