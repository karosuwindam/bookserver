package table

import (
	"bookserver/config"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/karosuwindam/sqlite"
)

// SQLStatus SQLの状態を渡すデータ
type SQLStatus struct {
	//SQLの読み取り設定
	Cfg sqlite.SqliteConfig
	//
	flag bool
}

// Setup (*config.Config) = *SQLStatus, error
//
// セットアップ情報
//
// cfg(*config.Config) : 設定情報
func Setup(cfg *config.Config) (*SQLStatus, error) {
	output := &SQLStatus{}
	output.Cfg = sqlite.Setup(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)
	if err := output.Cfg.Open(); err != nil {
		return nil, err
	}
	tablelistsetup()
	output.CreateTable()
	output.flag = true
	return output, nil
}

// CreateTable ()
func (sql *SQLStatus) CreateTable() {

	for name, typedata := range tablelist {
		sql.Cfg.CreateTable(name, typedata)
	}
}

// Close ()
func (sql *SQLStatus) Close() {
	sql.Cfg.Close()
}

// Add (tName, writedata)
func (sql *SQLStatus) Add(tName string, writedata interface{}) error {
	if !ckType(writedata) {
		return errors.New("input write type error")
	}
	if err := sql.Cfg.Add(tName, writedata); err != nil {
		return err
	}
	return nil
}

// ReadID (tName,id)
func (sql *SQLStatus) ReadID(tName string, id int) (string, error) {
	readdata := readBaseCreate(tName)
	key := map[string]string{"id": strconv.Itoa(id)}

	if err := sql.Cfg.Read(tName, readdata, key); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}

// ReadAll (tName)
func (sql *SQLStatus) ReadAll(tName string) (string, error) {
	readdata := readBaseCreate(tName)

	if err := sql.Cfg.Read(tName, readdata); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}
