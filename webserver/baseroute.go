package webserver

import (
	"bookserver/webserver/common"
	"bookserver/webserver/weblogin"
)

// 子供モジュールの設定処理
func (t *SetupServer) route() {
	wln := weblogin.Setup()
	t.AddV1(common.GUEST, "login", wln, weblogin.WebServerLogin)
	t.AddV1(common.GUEST, "logout", wln, weblogin.WebServerLogout)
}
