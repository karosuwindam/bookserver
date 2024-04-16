package controller

import (
	"bookserver/controller/convert"
	"bookserver/controller/copyfile"
	readzipfile "bookserver/controller/readZipfile"
)

type HealthData struct {
	Convert  interface{} `json:"Convert"`
	Cash     interface{} `json:"Cash"`
	Copyfile interface{} `json:"Copyfile"`
}

// 処理状態を確認する処理
func HealthCheck() HealthData {
	output := HealthData{}
	//変換情報を取得
	output.Convert = convert.CheackHealth()
	output.Cash = readzipfile.HealthCheck()
	output.Copyfile = copyfile.HealthCheck()

	//:ToDo
	return output
}
