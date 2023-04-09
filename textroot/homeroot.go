package textroot

import "bookserver/webserverv2"

var Route []webserverv2.WebConfig = []webserverv2.WebConfig{
	{Pass: "/edit/", Handler: viewhtml},
	{Pass: "/", Handler: viewhtml},
}
