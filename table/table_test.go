package table

import (
	"bookserver/config"
	"os"
	"testing"
)

func TestTableSetup(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, err := Setup(cfg)
	if err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	}
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)

}

func TestTableWriteRead(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)

}
