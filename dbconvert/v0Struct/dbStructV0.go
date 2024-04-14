package v0struct

import "time"

type Booknames struct {
	Id        uint   `json:"Id" gorm:"primarykey"`
	Name      string `json:"Name`
	Title     string `json:"Title"`
	Writer    string `json:"Writer"`
	Burand    string `json:"Burand"`
	Booktype  string `json:"Booktype"`
	Ext       string `json:"Ext"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Copyfile struct {
	Id        uint   `json:"Id" gorm:"primarykey"`
	Zippass   string `json:"Zippass"`
	Filesize  int    `json:"Filesize"`
	Copyflag  int    `json:"Copyflag"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (o Copyfile) TableName() string {
	return "copyfile"
}

type Filelists struct {
	Id        uint   `json:"Id" gorm:"primarykey"`
	Name      string `json:"Name"`
	Pdfpass   string `json:"Pdfpass"`
	Zippass   string `json:"Zippass"`
	Tag       string `json:"Tag"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
