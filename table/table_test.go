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

func TestTableRead(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	t.Setenv("DB_FILE", "test-read.db")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	if json, err := sql.ReadAll(BOOKNAME); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(json)
	}
	if json, err := sql.ReadAll(COPYFILE); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(json)
	}
	if json, err := sql.ReadAll(FILELIST); err != nil {
		t.Fatalf(err.Error())
		t.FailNow()
	} else {
		t.Log(json)
	}

}
func TestTableWrite(t *testing.T) {
	t.Setenv("DB_ROOTPASS", "./")
	cfg, _ := config.EnvRead()
	sql, _ := Setup(cfg)
	defer sql.Close()
	defer os.Remove(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)

}
