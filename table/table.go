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

// (sql)CreateTable ()
//
// テーブルを作成する
func (sql *SQLStatus) CreateTable() {

	for name, typedata := range tablelist {
		sql.Cfg.CreateTable(name, typedata)
	}
}

// (sql)Close ()
//
// sqlを閉じる
func (sql *SQLStatus) Close() {
	sql.Cfg.Close()
}

// (sql)Add (tName, writedata) = error
//
// 対象のテーブルにデータを書き込む
//
// tName : 対象テーブル
// writedata : 書き込むデータの構造体
func (sql *SQLStatus) Add(tName string, writedata interface{}) error {
	if !ckType(writedata) {
		return errors.New("input write type error")
	}
	if err := sql.Cfg.Add(tName, writedata); err != nil {
		return err
	}
	return nil
}

// (sql)Edit (tName, writedata, id) = (string, error)
//
// 対象のテーブルでidを指定してデータを編集する
// 編集した結果をjson形式で返す
//
// tName : 対象テーブル
// writedata : 書き込むデータの構造体
// id : 対象のid
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

// (sql)ReadID (tName,id) = (string, error)
//
// 対象のテーブルからidを指定してデータを参照する
// その結果をjson形式で返す
//
// tName : 対象テーブル
// id : 対象のid
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

// (sql)ReadAll (tName) = (string, error)
//
// 対象のテーブルデータをすべて参照する
// その結果をjson形式で返す
//
// tName : 対象テーブル
func (sql *SQLStatus) ReadAll(tName string) (string, error) {
	readdata := readBaseCreate(tName)
	if readdata == nil {
		return "", errors.New("Not found Table")
	}

	if err := sql.Cfg.Read(tName, readdata); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}

// (sql)ReadWhileTime (tName, datetype) = (string, error)
//
// 特殊検索で、キーワードによって日付ごとの更新時刻データを取得
//
// tName : 対象テーブル
// datetype(string) : 対象の特殊キーワード today, toweek, tomonth
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

// (sql) ReadName(tName, keyword) = (string, error)
//
// 名前キーで一致するファイルを探す
//
// tName : 対象テーブル
// keyword : 名前キーで検索するキーワード
func (sql *SQLStatus) ReadName(tName, keyword string) (string, error) {
	if keyword == "" {
		return "", nil
	}
	readdata := readBaseCreate(tName)
	skeyword := baseNameMap(tName, keyword)
	if err := sql.Cfg.Read(tName, readdata, skeyword, sqlite.AND); err != nil {
		return "", err
	}
	bJSON, err := json.Marshal(readdata)
	if err != nil {
		return "", err
	}
	return string(bJSON), nil
}

// (sql) Search(tName, keyword) = (string, error)
//
// キーワードを検索する
//
// tName : 対象テーブル
// keyword : 検索するキーワード
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

// (sql) Delete(tName, id) = (string, error)
//
// idを指定して削除する。
//
// tName : 対象テーブル
// id : 対象ID
func (sql *SQLStatus) Delete(tName string, id int) (string, error) {
	if err := sql.Cfg.Delete(tName, id); err != nil {
		return "", err
	}
	return "", nil
}

// baseNameMap(tName, keyword) = map[string]string
//
// 検索用のデータで基本となる部分をテーブルから作成
//
// tName : 対象テーブル
// keyword : 検索のキーワード
func baseNameMap(tName, keyword string) map[string]string {
	output := map[string]string{}
	switch tName {
	case BOOKNAME:
		output["name"] = keyword
	case FILELIST:
		output["name"] = keyword
	case COPYFILE:
		output["zippass"] = keyword
	}
	return output
}
