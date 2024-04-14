package booknames

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Booknames struct {
	Id       uint   `json:"Id" gorm:"primarykey"`
	Name     string `json:"Name" gorm:"unique"`
	Title    string `json:"Title"`
	Writer   string `json:"Writer"`
	Burand   string `json:"Burand"`
	Booktype string `json:"Booktype"`
	Ext      string `json:"Ext"`
}

type sqlBooknames struct {
	Booknames
	CreateAt time.Time
	UpdateAt time.Time
}

func (o sqlBooknames) TableName() string {
	return "booknames"
}

const (
	TODAY   time.Duration = 1
	TOWEEK  time.Duration = 7
	TOMONTH time.Duration = 30
)

var basedb *gorm.DB

func Init(db *gorm.DB) error {
	basedb = db
	return db.AutoMigrate(&sqlBooknames{})
}

func (t *Booknames) Add() error {
	if t.Name == "" {
		return errors.New("not input Name data")
	}
	nowTime := time.Now()
	tmp := sqlBooknames{
		Booknames: *t,
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

func GetAll() ([]Booknames, error) {
	out := []Booknames{}
	tmps := []sqlBooknames{}
	if results := basedb.Find(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		tmp := sbn.Booknames
		out = append(out, tmp)
	}
	return out, nil
}

func GetId(id int) (Booknames, error) {
	out := Booknames{}
	tmps := []sqlBooknames{}
	if results := basedb.Where("id = ?", id).First(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Booknames
		break
	}
	return out, nil
}

func GetName(name string) (Booknames, error) {
	out := Booknames{}
	tmps := []sqlBooknames{}
	if results := basedb.Where("name = ?", name).Last(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, sbn := range tmps {
		out = sbn.Booknames
		break
	}
	return out, nil
}

func (t *Booknames) Update() error {

	tmps := []sqlBooknames{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	if tmps[0].Id != t.Id {
		return errors.New("Not input id by" + strconv.Itoa(int(t.Id)))
	} else {
		tmps[0].Booknames = *t
		tmps[0].UpdateAt = time.Now()
	}
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	return nil
}

func Delete(id int) error {
	if results := basedb.Delete(&sqlBooknames{}, id); results.Error != nil {
		return results.Error
	}
	return nil
}

func Search(keyword string) ([]Booknames, error) {
	output := []Booknames{}
	tmps := []sqlBooknames{}
	basekey := []string{
		"name",
		"title",
		"writer",
		"burand",
		"ext",
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
		bns := tmp.Booknames
		output = append(output, bns)
	}

	return output, nil
}

func ReadByDataRangeDay(days time.Duration) ([]Booknames, error) {
	output := []Booknames{}
	tmpSQL := []sqlBooknames{}
	t := time.Now()
	if results := basedb.Where("update_at between ? and ?", t.Add(-days*24*time.Hour).Format("2006-01-02"), t.Add(24*time.Hour).Format("2006-01-02")).Find(&tmpSQL); results.Error != nil {
		return output, results.Error
	}
	for _, tmp := range tmpSQL {
		output = append(output, tmp.Booknames)
	}
	return output, nil
}

func ReadRandData(count int) ([]Booknames, error) {
	output := []Booknames{}
	tmpSQL := []sqlBooknames{}
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
		output = append(output, tmp.Booknames)
	}
	return output, nil
}
