package table

import (
	"bookserver/config"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type BookServer struct {
	TableVersion int
}

const (
	V0Table int = 0 // Version 0.9 以下なら
	V1Table int = 1 // Version 0.10 以上
)

var ErrorTableVersion error = errors.New("err table version")

func ReadTableVersion(db *gorm.DB) (BookServer, error) {
	output := []BookServer{}

	if results := db.First(&output); results.Error != nil {
		return BookServer{}, results.Error
	}

	return output[0], nil
}

func InitVTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&BookServer{}); err != nil {
		return err
	}
	if d, err := ReadTableVersion(db); err != nil {
		return intputNowVersion(db)
	} else {
		if d.TableVersion != VersionChack() {
			//テーブルのバージョンが異なるときの処理
			return ErrorTableVersion
		}
	}
	return nil
}

func intputNowVersion(db *gorm.DB) error {
	input := []BookServer{
		{VersionChack()},
	}
	if results := db.Save(input); results.Error != nil {
		return results.Error
	}
	return nil

}

func VersionChack() int {
	version := V0Table
	data := 0
	tmp := strings.Split(config.Version, ".")
	if len(tmp) > 2 {
		i, _ := strconv.Atoi(tmp[0])
		data = i * 100
		i, _ = strconv.Atoi(tmp[1])
		data = i
	}
	if data < 9 {
		version = V0Table
	} else {
		version = V1Table
	}
	return version
}
