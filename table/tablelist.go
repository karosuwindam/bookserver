package table

import (
	"encoding/json"
	"log"
	"reflect"
)

type Booknames struct {
	Id       int    `json:"Id" db:"id"`
	Name     string `json:"Name" db:"name"`
	Title    string `json:"Title" db:"title"`
	Writer   string `json:"Writer" db:"writer"`
	Burand   string `json:"Burand" db:"burand"`
	Booktype string `json:"Booktype" db:"booktype"`
	Ext      string `json:"Ext" db:"ext"`
}

type Copyfile struct {
	Id       int    `json:"Id" db:"id"`
	Zippass  string `json:"Zippass" db:"zippass"`
	Filesize int    `json:"Filesize" db:"filesize"`
	Copyflag int    `json:"Copyflag" db:"copyflag"`
}

type Filelists struct {
	Id      int    `json:"Id" db:"id"`
	Name    string `json:"Name" db:"name"`
	Pdfpass string `json:"Pdfpass" db:"pdfpass"`
	Zippass string `json:"Zippass" db:"zippass"`
	Tag     string `json:"Tag" db:"tag"`
}

const (
	BOOKNAME = "booknames"
	COPYFILE = "copyfile"
	FILELIST = "filelists"
)

// テーブル内の型情報を格納
var tablelist map[string]interface{}

// tablelistsetup()
//
// 初期化用関数、テーブルリストを作成する。
func tablelistsetup() {
	tablelist = map[string]interface{}{}
	tablelist[BOOKNAME] = Booknames{}
	tablelist[COPYFILE] = Copyfile{}
	tablelist[FILELIST] = Filelists{}
	return
}

// readBaseCreate(string) = interface{}
//
// SQL読み取り用の型を作成
func readBaseCreate(tname string) interface{} {
	var out interface{}
	switch tname {
	case BOOKNAME:
		out = &[]Booknames{}
	case COPYFILE:
		out = &[]Copyfile{}
	case FILELIST:
		out = &[]Filelists{}
	default:
		return nil
	}

	return out
}

// CkList (string) = bool
//
// 名前が含まれているか確認する関数
func CkList(tName string) bool {
	if tablelist[tName] != nil {
		return true
	}
	return false
}

// ckType(interface{}) = bool
//
// 変数の型の確認
//
// a(interface{}) : 型を代入
func ckType(a interface{}) bool {
	switch a.(type) {
	case *Booknames, *Filelists, *Copyfile:
		return true
	case *[]Booknames, *[]Filelists, *[]Copyfile:
		return true

	}
	return false
}

// createSerchText (tname, keyword) = map[string]string
//
// 検索用のmap配列を作る
// 対象の構造体からstringを探して構造体に挿入
func createSerchText(tname, keyword string) map[string]string {
	output := map[string]string{}
	if tablelist[tname] == nil {
		return output
	}
	st := reflect.TypeOf(tablelist[tname])
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if f.Type.Kind() == reflect.String && f.Tag.Get("db") != "" {
			output[f.Tag.Get("db")] = keyword
		}
	}
	return output
}

// JsonToStruct(tName data) = interface{}
//
// Jsonデータを指定した型に変換実施
func JsonToStruct(tName string, data []byte) interface{} {
	var out interface{}
	switch tName {
	case BOOKNAME:
		if string(data)[:1] == "[" {
			tmp := []Booknames{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp
		} else {
			tmp := Booknames{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp
		}
	case COPYFILE:
		if string(data)[:1] == "[" {
			tmp := []Copyfile{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp
		} else {
			tmp := Copyfile{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp
		}
	case FILELIST:
		if string(data)[:1] == "[" {

			tmp := []Filelists{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp
		} else {
			tmp := Filelists{}
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Println(err)
				return nil
			}
			out = tmp

		}
	default:
		return nil

	}
	return out
}
