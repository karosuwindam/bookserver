package textroot

import (
	"bookserver/config"
	"bookserver/textroot/view"
	"bookserver/textroot/viewpage"
	"bookserver/webserverv2"
)

var Route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/edit/", Handler: viewpage.Viewhtml},
	{Pass: "/", Handler: viewpage.Viewhtml},
}

func Setup(cfg *config.Config) error {
	if route, err := view.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, route...)
	}
	return nil
}
