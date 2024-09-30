package table

import (
	"bookserver/config"
	"bookserver/dbconvert"
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"bookserver/table/historyviews"
	"bookserver/table/uploadtmp"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
)

var db *gorm.DB

func Init() error {
	var err error
	if config.DB.DBNAME == "sqlite3" {
		db, err = openSqlite3()
	} else {
		return errors.New("Not Db type")
	}
	if err != nil {
		return err
	}
	if err = tableInit(db); err != nil {
		return errors.Wrap(err, "tableInit")
	}
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return errors.Wrap(err, "tableTrace")
	}
	return nil
}

func openSqlite3() (*gorm.DB, error) {
	filepass := config.DB.DBROOTPASS
	if filepass[len(filepass)-1:] != "/" {
		filepass += "/"
	}
	if f, err := os.Stat(filepass); os.IsNotExist(err) || !f.IsDir() {
		os.MkdirAll(filepass, 0766)
	}
	filepass += config.DB.DBFILE
	if _, err := os.Stat(filepass); err != nil {
		fmt.Println("Create sqlite3 file: ", filepass)
	} else { //ファイルがある場合の処理
		tmpdb, err := gorm.Open(sqlite.Open(filepass), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		tableVersion := 0
		if d, err := ReadTableVersion(tmpdb); err != nil {
			//読み取りに失敗する場合はv0として判断
		} else {
			tableVersion = d.TableVersion
		}
		//バージョン情報を読み取り
		if VersionChack() != tableVersion {
			//バージョンが異なる場合
			tmpFIlePass := ""
			//ファイルのバックアップ
			if i := strings.LastIndex(filepass, "."); i > 0 {
				tmpFIlePass = filepass[:i]
				tmpFIlePass += "_v" + strconv.Itoa(tableVersion) + "_" + time.Now().Format("20060102")
				tmpFIlePass += filepass[i:]
			}
			if err := fileCopy(filepass, tmpFIlePass); err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("file Copy error %v to %v", filepass, tmpFIlePass))
			}
			tmpdb = nil
			os.Remove(filepass)
			tmpdb, _ = gorm.Open(sqlite.Open(filepass), &gorm.Config{})
			if err := tableInit(tmpdb); err != nil {
				os.Remove(filepass)
				fileCopy(tmpFIlePass, filepass)
				os.Remove(tmpFIlePass)
				return nil, errors.Wrap(err, fmt.Sprintf("taible Init err"))
			}
			//データのコンバート処理
			if err := dbconvert.DbConvertSQL3(tmpFIlePass, filepass, tableVersion, VersionChack()); err != nil {
				return nil, err
			}

		}

	}
	logger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	// return gorm.Open(sqlite.Open(filepass), &gorm.Config{})
	return gorm.Open(sqlite.Open(filepass), &gorm.Config{Logger: logger})

}

func tableInit(db *gorm.DB) error {
	if err := booknames.Init(db); err != nil {
		return errors.Wrap(err, "booknames table init error")
	}
	if err := filelists.Init(db); err != nil {
		return errors.Wrap(err, "filelists table init error")
	}
	if err := copyfiles.Init(db); err != nil {
		return errors.Wrap(err, "copyfiles table init error")
	}
	if err := uploadtmp.Init(db); err != nil {
		return errors.Wrap(err, "uploadtmp table init error")
	}
	if err := historyviews.Init(db); err != nil {
		return errors.Wrap(err, "historyviews table init error")
	}
	//テーブルバージョン初期化
	if err := InitVTable(db); err != nil {
		return errors.Wrap(err, "bookservers table init error")
	}
	return nil
}

func fileCopy(srcName, dstName string) error {

	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
