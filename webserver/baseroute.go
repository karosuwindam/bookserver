package webserver

import (
	"bookserver/config"
	"bookserver/webserver/weblogin"
)

// 子供モジュールの設定処理
func (t *SetupServer) route(cfg *config.Config) {
	weblogin.Setup(cfg)
	t.Add("/login", weblogin.WebServerLogin)
	t.Add("/logout", weblogin.WebServerLogout)
}
