package v1struct

import "time"

type Copyfile struct {
	Id       uint   `json:"Id" gorm:"primarykey"`
	Zippass  string `json:"Zippass"`
	Filesize int    `json:"Filesize"`
	Copyflag int    `json:"Copyflag"`
}

type Copyfile_sql struct {
	Copyfile
	CreateAt *time.Time
	UpdateAt *time.Time
}

func (o Copyfile_sql) TableName() string {
	return "copyfile"
}

type Filelists struct {
	Id      uint   `json:"Id" gorm:"primarykey"`
	Name    string `json:"Name" gorm:"unique"`
	Pdfpass string `json:"Pdfpass"`
	Zippass string `json:"Zippass"`
	Tag     string `json:"Tag"`
}

type Filelists_sql struct {
	Filelists
	CreateAt *time.Time
	UpdateAt *time.Time
}

func (o Filelists_sql) TableName() string {
	return "filelists"
}

type Booknames struct {
	Id       uint   `json:"Id" gorm:"primarykey"`
	Name     string `json:"Name" gorm:"unique"`
	Title    string `json:"Title"`
	Writer   string `json:"Writer"`
	Burand   string `json:"Burand"`
	Booktype string `json:"Booktype"`
	Ext      string `json:"Ext"`
}

type Booknames_sql struct {
	Booknames
	CreateAt *time.Time
	UpdateAt *time.Time
}

func (o Booknames_sql) TableName() string {
	return "booknames"
}

type UploadTmp struct {
	Id         uint   `gorm:"primarykey"`
	Name       string //アップロードしたファイル名
	SavePdf    string //PDFの保存先
	SaveZip    string //Zipの保存先
	Count      int    //コンバートを実行した回数
	ThreadFlag bool   //定期処理フラグ
}

type UploadTmp_sql struct {
	UploadTmp
	CreateAt time.Time
	UpdateAt time.Time
}

func (o UploadTmp_sql) TableName() string {
	return "uploadtmps"
}

// アクセスファイル登録情報
type HistoryViews struct {
	Id     uint   `json:"Id" gorm:"primarykey"` //ID
	FileId int    `json:"FileId"`               //アクセスしたファイルのID
	Ip     string `json:"Ip"`                   //アクセス元のIPアドレス
	User   string `json:"User"`                 //アクセスしたユーザアカウントID 未実装
}

type HistoryViews_sql struct {
	HistoryViews
	CreateAt time.Time
	UpdateAt time.Time
}

func (o HistoryViews_sql) TableName() string {
	return "historyviews"
}
