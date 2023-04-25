package api

import (
	"bookserver/api/common"
	"bookserver/api/copy"
	"bookserver/api/download"
	"bookserver/api/login"
	"bookserver/api/tabledata"
	"bookserver/api/upload"
	"bookserver/api/view"
	"bookserver/config"
	"bookserver/webserverv2"
)

var Route []webserverv2.WebConfig = []webserverv2.WebConfig{}

func Setup(cfg *config.Config) error {
	//common
	if err := common.Setup(); err != nil {
		return err
	}
	//login
	if tmp, err := login.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}

	//upload
	if tmp, err := upload.Setup(); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}

	//table
	if tmp, err := tabledata.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//view
	if tmp, err := view.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//download
	if tmp, err := download.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	//copy
	if tmp, err := copy.Setup(cfg); err != nil {
		return err
	} else {
		Route = append(Route, tmp...)
	}
	return nil
}
