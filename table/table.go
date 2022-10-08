package table

import (
	"bookserver/config"
	"encoding/json"
	"strconv"

	"github.com/karosuwindam/sqlite"
)

const (
	DbPass = "./db/"
)

type SQLStatus struct {
	Sql  sqlite.SqliteConfig
	flag bool
}

func Setup(cfg *config.Config) (*SQLStatus, error) {
	output := &SQLStatus{}
	output.Sql = sqlite.Setup(cfg.Sql.DBROOTPASS + cfg.Sql.DBFILE)
	if err := output.Sql.Open(); err != nil {
		return nil, err
	}
	tablelistsetup()
	output.CreateTable()
	output.flag = true
	return output, nil
}

func (sql *SQLStatus) CreateTable() {

	for name, typedata := range tablelist {
		sql.Sql.CreateTable(name, typedata)
	}
}

func (sql *SQLStatus) Close() {
	sql.Sql.Close()
}

func (sql *SQLStatus) Add(tName, wJson string) error {
	//jsonｋらデータへ
	return nil
}

func (sql *SQLStatus) ReadId(tName string, id int) (string, error) {
	readdata := readBaseCreate(tName)
	key := map[string]string{"id": strconv.Itoa(id)}

	if err := sql.Sql.Read(tName, &readdata.pData, key); err != nil {
		return "", err
	}
	bJson, err := json.Marshal(readdata.pData)
	if err != nil {
		return "", err
	}
	return string(bJson), nil
}

func (sql *SQLStatus) ReadAll(tName string) (string, error) {
	readdata := readBaseCreate(tName)

	if err := sql.Sql.Read(tName, &readdata.pData); err != nil {
		return "", err
	}
	bJson, err := json.Marshal(readdata.pData)
	if err != nil {
		return "", err
	}
	return string(bJson), nil
}
