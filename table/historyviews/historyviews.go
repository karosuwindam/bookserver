package historyviews

import (
	"bookserver/config"
	"bookserver/table/filelists"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// アクセスファイル登録情報
type HistoryViews struct {
	Id     uint   `json:"Id" gorm:"primarykey"` //ID
	FileId int    `json:"FileId"`               //アクセスしたファイルのID
	Ip     string `json:"Ip"`                   //アクセス元のIPアドレス
	User   string `json:"User"`                 //アクセスしたユーザアカウントID 未実装
}

type sqlHistoryViews struct {
	HistoryViews
	CreateAt time.Time
	UpdateAt time.Time
}

func (o sqlHistoryViews) TableName() string {
	return "historyviews"
}

var basedb *gorm.DB

func Init(db *gorm.DB) error {
	basedb = db
	return db.AutoMigrate(&sqlHistoryViews{})
}

func (t *HistoryViews) Add() error {
	if _, err := filelists.GetId(t.FileId); err != nil {
		return errors.Wrap(err, fmt.Sprintf("filelists.GetId(%v)", t.FileId))
	}
	nowTime := time.Now()

	tmp := sqlHistoryViews{
		HistoryViews: *t,
		CreateAt:     nowTime,
		UpdateAt:     nowTime,
	}
	if tmp.Id != 0 {
		tmp.Id = 0
	}
	if results := basedb.Create(&tmp); results.Error != nil {
		return results.Error
	}
	return nil
}

func GetHistory(n int) ([]filelists.Filelists, error) {
	var out []filelists.Filelists
	var tmps, historys []sqlHistoryViews
	ctx := context.TODO()
	if results := basedb.Order("create_at DESC").Limit(config.BScfg.HistoryMax).Find(&historys); results.Error != nil {
		return out, results.Error
	}
	for _, history := range historys {
		flag := false
		for _, tmp := range tmps {
			if tmp.FileId == history.FileId {
				flag = true
				break
			}
		}
		if !flag {
			tmps = append(tmps, history)
		}
		if len(tmps) >= n {
			break
		}
	}
	for _, tmp := range tmps {
		if b, err := filelists.GetId(int(tmp.FileId)); err == nil {
			out = append(out, b)
		} else {
			slog.ErrorContext(ctx,
				fmt.Sprintf("GetHistory filelists.GetId(%v) error", tmp.FileId),
				"id", tmp.FileId,
				"Error", err,
			)
		}
	}
	return out, nil
}
