package filelists

import (
	"bookserver/config"
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// 登録したファイル情報
type Filelists struct {
	Id      uint   `json:"Id" gorm:"primarykey"`
	Name    string `json:"Name" gorm:"unique"`
	Pdfpass string `json:"Pdfpass"`
	Zippass string `json:"Zippass"`
	Tag     string `json:"Tag"`
}

type sqlFilelists struct {
	Filelists
	CreateAt time.Time
	UpdateAt time.Time
}

const (
	TODAY   time.Duration = 1
	TOWEEK  time.Duration = 7
	TOMONTH time.Duration = 30
)

func (o sqlFilelists) TableName() string {
	return "filelists"
}

var basedb *gorm.DB

func Init(db *gorm.DB) error {
	basedb = db
	return db.AutoMigrate(&sqlFilelists{})
}

func (t *Filelists) Add() error {
	if t.Name == "" {
		return errors.New("not input Name data")
	}
	nowTime := time.Now()

	tmp := sqlFilelists{
		Filelists: *t,
		CreateAt:  nowTime,
		UpdateAt:  nowTime,
	}
	if tmp.Id != 0 {
		tmp.Id = 0
	}
	if results := basedb.Create(&tmp); results.Error != nil {
		return results.Error
	}
	return nil
}

func GetAll(ctx context.Context) ([]Filelists, error) {
	ctx, span := config.TracerS(ctx, "GetAll", "Get All Filelists")
	defer span.End()

	out := []Filelists{}
	tmps := []sqlFilelists{}
	if results := basedb.Find(&tmps); results.Error != nil {
		config.TracerError(span, results.Error)
		return out, results.Error
	}
	for _, sbn := range tmps {
		tmp := sbn.Filelists
		out = append(out, tmp)
	}
	return out, nil
}

func GetId(id int) (Filelists, error) {
	out := Filelists{}
	tmps := []sqlFilelists{}
	if results := basedb.Where("id = ?", id).First(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Filelists
		break
	}
	return out, nil
}

func GetName(name string) (Filelists, error) {
	out := Filelists{}
	tmps := []sqlFilelists{}
	if results := basedb.Where("name = ?", name).Last(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Filelists
		break
	}
	return out, nil
}

func (t *Filelists) Update() error {

	tmps := []sqlFilelists{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	if tmps[0].Id != t.Id {
		return errors.New("Not input id by" + strconv.Itoa(int(t.Id)))
	} else {
		tmps[0].Filelists = *t
		tmps[0].UpdateAt = time.Now()
	}
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	return nil
}

func Delete(id int) error {
	if results := basedb.Delete(&sqlFilelists{}, id); results.Error != nil {
		return results.Error
	}
	return nil
}

func Search(keyword string) ([]Filelists, error) {

	output := []Filelists{}
	tmps := []sqlFilelists{}
	basekey := []string{
		"name",
		"tag",
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
		bns := tmp.Filelists
		output = append(output, bns)
	}

	return output, nil
}

func ReadByDataRangeDay(days time.Duration) ([]Filelists, error) {
	output := []Filelists{}
	tmpSQL := []sqlFilelists{}
	t := time.Now()
	if results := basedb.Where("update_at between ? and ?", t.Add(-days*24*time.Hour).Format("2006-01-02"), t.Add(24*time.Hour).Format("2006-01-02")).Find(&tmpSQL); results.Error != nil {
		return output, results.Error
	}
	for _, tmp := range tmpSQL {
		output = append(output, tmp.Filelists)
	}
	return output, nil
}

func ReadRandData(count int) ([]Filelists, error) {
	output := []Filelists{}
	tmpSQL := []sqlFilelists{}
	type tmpid struct {
		Id uint
	}
	tmpId := []tmpid{}
	if results := basedb.Select("id").Find(&tmpSQL).Scan(&tmpId); results.Error != nil {
		return output, results.Error
	}
	if len(tmpId) <= count {
		if results := basedb.Find(&tmpSQL); results.Error != nil {
			return output, results.Error
		}
	} else {
		//ランダムでcountの数だけIDを取り出す
		randids := []any{}
		for i := 0; i < count; i++ {
			tmpCount := rand.Intn(len(tmpId) - 1)
			flag := true
			for _, randid := range randids {
				if randid == tmpId[tmpCount].Id {
					flag = false
					break
				}
			}
			if flag {
				randids = append(randids, tmpId[tmpCount].Id)
			} else {
				i--
			}
		}
		cmd := "id = ?"
		for i := 1; i < len(randids); i++ {
			cmd += " or id = ?"
		}
		if results := basedb.Where(cmd, randids...).Find(&tmpSQL); results.Error != nil {
			return output, results.Error
		}

	}
	for _, tmp := range tmpSQL {
		output = append(output, tmp.Filelists)
	}
	return output, nil
}
