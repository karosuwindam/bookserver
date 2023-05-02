package tabledata

import (
	"bookserver/api/tabledata/add"
	"bookserver/api/tabledata/edit"
	"bookserver/api/tabledata/read"
	"bookserver/api/tabledata/search"
	"bookserver/config"
	"bookserver/table"
	"bookserver/webserver"
)

// routeのベースフォルダ
var route []webserver.WebConfig = []webserver.WebConfig{}

func sqlSetup(cfg *config.Config) (*table.SQLStatus, error) {
	var err error
	if sqlcfg, err := table.Setup(cfg); err == nil {
		return sqlcfg, err
	}
	return nil, err
}

// Setup() = []webserver.WebConfig
//
// セットアップして、HTMLのルートフォルダを用意する
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	var err error
	if sqlcfg, err := sqlSetup(cfg); err == nil {
		route = append(route, add.Setup(sqlcfg)...)
		route = append(route, edit.Setup(sqlcfg)...)
		route = append(route, read.Setup(sqlcfg)...)
		route = append(route, search.Setup(sqlcfg)...)
	}
	return route, err
}
