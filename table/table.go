package table

import (
	"bookserver/config"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

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

// Edit (tName, writedata, id)
func (sql *SQLStatus) Edit(tName string, writedata interface{}, id int) (string, error) {
	readdata := readBaseCreate(tName)
	key := map[string]string{"id": strconv.Itoa(id)}

	if err := sql.Cfg.Read(tName, readdata, key); err != nil {
		return "", err
	}
	// wv := reflect.ValueOf(writedata)
	if err := sql.Cfg.Update(tName, writedata); err != nil {
		return "", err
	}
	readdata = readBaseCreate(tName)
	if err := sql.Cfg.Read(tName, readdata, key); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
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

// ReadWhileTime (tName, datetype)
// datetype(string) : today, toweek, tomonth
func (sql *SQLStatus) ReadWhileTime(tName, datetype string) (string, error) {
	readdata := readBaseCreate(tName)

	switch strings.ToLower(datetype) {
	case "today":
		if err := sql.Cfg.ReadToday(tName, readdata); err != nil {
			return "", err
		}
	case "toweek":
		if err := sql.Cfg.ReadToWeek(tName, readdata); err != nil {
			return "", err
		}
	case "tomonth":
		if err := sql.Cfg.ReadToMonth(tName, readdata); err != nil {
			return "", err
		}
	default:
		return "", errors.New("datetype err :" + datetype)
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil

}

// 一致するファイルを探す
func (sql *SQLStatus) ReadName(tName, keyword string) (string, error) {
	if keyword == "" {
		return "", nil
	}
	readdata := readBaseCreate(tName)
	skeyword := map[string]string{"name": keyword}
	if err := sql.Cfg.Read(tName, readdata, skeyword, sqlite.AND); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}

// キーワードを検索する
func (sql *SQLStatus) Search(tName, keyword string) (string, error) {
	if keyword == "" {
		return "", nil
	}
	readdata := readBaseCreate(tName)
	skeyword := createSerchText(tName, keyword)
	if err := sql.Cfg.Read(tName, readdata, skeyword, sqlite.ORLike); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}

// idを指定して削除する。
func (sql *SQLStatus) Delete(tName string, id int) (string, error) {
	if err := sql.Cfg.Delete(tName, id); err != nil {
		return "", err
	}
	return "", nil
}
