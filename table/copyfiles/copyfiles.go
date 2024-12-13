package copyfiles

import (
	"bookserver/config"
	"context"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// 共有フォルダに登録したファイル情報
type Copyfile struct {
	Id       uint   `json:"Id" gorm:"primarykey;"`
	Zippass  string `json:"Zippass"`
	Filesize int    `json:"Filesize"`
	Copyflag int    `json:"Copyflag"`
}

type sqlCopyfile struct {
	Copyfile
	CreateAt time.Time
	UpdateAt time.Time
}

func (o sqlCopyfile) TableName() string {
	return "copyfile"
}

const (
	OFF int = 0
	ON  int = 1
)

var basedb *gorm.DB

func Init(db *gorm.DB) error {
	basedb = db
	return db.AutoMigrate(&sqlCopyfile{})
}

func (t *Copyfile) Add() error {
	if t.Zippass == "" {
		return errors.New("not input Name data")
	}
	nowTime := time.Now()
	tmp := sqlCopyfile{
		Copyfile: *t,
		CreateAt: nowTime,
		UpdateAt: nowTime,
	}
	if tmp.Id != 0 {
		tmp.Id = 0
	}
	if results := basedb.Create(&tmp); results.Error != nil {
		return results.Error
	}
	return nil
}

func GetAll(ctx context.Context) ([]Copyfile, error) {
	ctx, span := config.TracerS(ctx, "GetAll", "Get All Copyfile")
	defer span.End()
	out := []Copyfile{}
	tmps := []sqlCopyfile{}
	if results := basedb.Find(&tmps); results.Error != nil {
		config.TracerError(span, results.Error)
		return out, results.Error
	}
	for _, sbn := range tmps {
		tmp := sbn.Copyfile
		out = append(out, tmp)
	}
	return out, nil
}

func GetId(id int) (Copyfile, error) {
	out := Copyfile{}
	tmps := []sqlCopyfile{}
	if results := basedb.Where("id = ?", id).First(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Copyfile
		break
	}
	return out, nil
}

func GetZipName(filename string) (Copyfile, error) {
	out := Copyfile{}
	tmps := []sqlCopyfile{}
	if results := basedb.Where("zippass = ?", filename).First(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Copyfile
		break
	}
	return out, nil
}

func (t *Copyfile) Update() error {

	tmps := []sqlCopyfile{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	if tmps[0].Id != t.Id {
		return errors.New("Not input id by" + strconv.Itoa(int(t.Id)))
	} else {
		tmps[0].Copyfile = *t
		tmps[0].UpdateAt = time.Now()
	}
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	return nil
}

// idにデータが入っているものを強制的にCopyFlagをONにする
func (t *Copyfile) ON() error {
	tmps := []sqlCopyfile{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	if tmps[0].Id != t.Id {
		return errors.New("Not input id by" + strconv.Itoa(int(t.Id)))
	} else {
		tmps[0].Copyfile.Copyflag = ON
		tmps[0].UpdateAt = time.Now()
	}
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	return nil
}

// idにデータが入っているものを強制的にCopyFlagをOFFにする
func (t *Copyfile) OFF() error {
	tmps := []sqlCopyfile{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	if tmps[0].Id != t.Id {
		return errors.New("Not input id by" + strconv.Itoa(int(t.Id)))
	} else {
		tmps[0].Copyfile.Copyflag = OFF
		tmps[0].UpdateAt = time.Now()
	}
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	return nil
}

func Delete(id int) error {
	if results := basedb.Delete(&sqlCopyfile{}, id); results.Error != nil {
		return results.Error
	}
	return nil
}

func Search(keyword string) ([]Copyfile, error) {
	output := []Copyfile{}
	tmps := []sqlCopyfile{}
	basekey := []string{
		"zippass",
	}
	cmd := ""
	if keyword == "" {
		return output, nil
	}
	for i, s := range basekey {
		if i == 0 {
			cmd = s + " LIKE " + "\"%" + keyword + "%\""
		} else {
			cmd += " OR " + s + " LIKE " + "\"%" + keyword + "%\""
		}
	}
	if results := basedb.Where(cmd).Find(&tmps); results.Error != nil {
		return output, results.Error
	}
	for _, tmp := range tmps {
		bns := tmp.Copyfile
		output = append(output, bns)
	}

	return output, nil
}

func OnOFFSearch(flag int) ([]Copyfile, error) {
	output := []Copyfile{}
	tmps := []sqlCopyfile{}
	if results := basedb.Where("copyflag = ?", flag).Find(&tmps); results.Error != nil {
		return output, results.Error
	}
	for _, tmp := range tmps {
		bns := tmp.Copyfile
		output = append(output, bns)
	}
	return output, nil
}
