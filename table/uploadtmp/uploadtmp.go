package uploadtmp

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

type UploadTmp struct {
	Id         uint   `gorm:"primarykey"`
	Name       string //アップロードしたファイル名
	SavePdf    string //PDFの保存先
	SaveZip    string //Zipの保存先
	Count      int    //コンバートを実行した回数
	ThreadFlag bool   //定期処理フラグ
}

type sqlUploadTmp struct {
	UploadTmp
	CreateAt time.Time
	UpdateAt time.Time
}

func (o sqlUploadTmp) TableName() string {
	return "uploadtmps"
}

var basedb *gorm.DB

var mu sync.Mutex

func Init(db *gorm.DB) error {
	basedb = db
	return db.AutoMigrate(&sqlUploadTmp{})
}

func (t *UploadTmp) Add() error {
	if t.Name == "" {
		return errors.New("not input Name data")
	}
	nowTime := time.Now()
	tmp := sqlUploadTmp{
		UploadTmp: *t,
		CreateAt:  nowTime,
		UpdateAt:  nowTime,
	}
	if tmp.Id != 0 {
		tmp.Id = 0
	}
	if results := basedb.Create(&tmp); results.Error != nil {
		return results.Error
	}
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("name = ?", t.Name).Last(&tmps); results.Error != nil {
		return results.Error
	}
	t.Id = tmps[0].Id
	return nil
}

func GetAllbyFlag(flag bool) ([]UploadTmp, error) {
	out := []UploadTmp{}
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("thread_flag = ?", flag).Find(&tmps); results.Error != nil {
		return out, results.Error
	}
	for _, tmp := range tmps {
		out = append(out, tmp.UploadTmp)
	}
	return out, nil
}

func (t *UploadTmp) CheckName() error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("name = ? AND thread_flag = ?", t.Name, false).First(&tmps); results.Error != nil && results.Error != gorm.ErrRecordNotFound {
		return results.Error
	}
	for _, tmp := range tmps {
		t.Id = tmp.Id
	}
	return nil
}

func (t *UploadTmp) CountClear() error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	tmps[0].Count = 0
	tmps[0].UpdateAt = time.Now()
	mu.Lock()
	defer mu.Unlock()
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	t.ThreadFlag = true
	t.Count = tmps[0].Count
	return nil
}

func (t *UploadTmp) CountUp() error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	tmps[0].Count++
	tmps[0].UpdateAt = time.Now()
	mu.Lock()
	defer mu.Unlock()
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	t.ThreadFlag = true
	t.Count = tmps[0].Count
	return nil
}

func (t *UploadTmp) FlagOn() error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	tmps[0].ThreadFlag = true
	tmps[0].UpdateAt = time.Now()
	mu.Lock()
	defer mu.Unlock()
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	t.ThreadFlag = true
	return nil
}

func (t *UploadTmp) SetPdfPath(name string) error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	tmps[0].SavePdf = name
	tmps[0].UpdateAt = time.Now()
	mu.Lock()
	defer mu.Unlock()
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	t.SavePdf = tmps[0].SavePdf
	return nil

}

func (t *UploadTmp) SetZipPath(name string) error {
	tmps := []sqlUploadTmp{}
	if results := basedb.Where("id = ?", t.Id).First(&tmps); results.Error != nil {
		return results.Error
	}
	tmps[0].SaveZip = name
	tmps[0].UpdateAt = time.Now()
	mu.Lock()
	defer mu.Unlock()
	if results := basedb.Save(tmps); results.Error != nil {
		return results.Error
	}
	t.SaveZip = tmps[0].SaveZip
	return nil

}
