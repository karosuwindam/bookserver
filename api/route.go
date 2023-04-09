package api

import (
	"bookserver/api/common"
	"bookserver/api/upload"
	"bookserver/webserverv2"
)

var Route []webserverv2.WebConfig = []webserverv2.WebConfig{}

func Setup() error {
	if err := common.Setup(); err != nil {
		return err
	}

	if tmp, err := upload.Setup(); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	return nil
}
