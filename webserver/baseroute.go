package webserver

import (
	"bookserver/config"
	"bookserver/webserver/add"
	"bookserver/webserver/common"
	"bookserver/webserver/edit"
	"bookserver/webserver/read"
	"bookserver/webserver/search"
	"bookserver/webserver/weblogin"
)

// 子供モジュールの設定処理
func (t *SetupServer) route(cfg *config.Config) {
	weblogin.Setup(cfg)
	t.Add("/login", weblogin.WebServerLogin)
	t.Add("/logout", weblogin.WebServerLogout)
	if comcfg, err := common.Setup(cfg); err == nil {
		t.sql = comcfg.Sql
		t.AddV1(common.GUEST, "/read/", comcfg.Sql, read.WebSQLRead)
		t.AddV1(common.GUEST, "/search/", comcfg.Sql, search.WebSQLSearch)
		t.AddV1(common.ADMIN, "/add/", comcfg.Sql, add.WebSQLRead)
		t.AddV1(common.ADMIN, "/edit/", comcfg.Sql, edit.WebSQLEdit)
	}
}
