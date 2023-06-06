package textroot

import (
	"bookserver/config"
	"bookserver/textroot/view"
	"bookserver/textroot/viewpage"
	"bookserver/webserver"
)

var Route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/edit/", Handler: viewpage.Viewhtml},
	{Pass: "/", Handler: viewpage.Viewhtml},
}

// Setup(cfg) = error
//
// セットアップ設定
//
// cfg : 基本設定
func Setup(cfg *config.Config) error {
	if route, err := view.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, route...)
	}
	if err := viewpage.Setup(cfg); err != nil {
		return err
	}
	return nil
}
