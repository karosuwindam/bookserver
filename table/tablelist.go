package table

import "reflect"

type Booknames struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Title    string `json:"title" db:"title"`
	Writer   string `json:"writer" db:"writer"`
	Brand    string `json:"brand" db:"brand"`
	Booktype string `json:"booktype" db:"booktype"`
	Ext      string `json:"ext" db:"ext"`
}

type Copyfile struct {
	Id       int    `json:"id" db:"id"`
	Zippass  string `json:"zippass" db:"zippass"`
	Filesize int    `json:"filesize" db:"filesize"`
	Copyflag int    `json:"copyflag" db:"copyflag"`
}

type Filelists struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Pdfpass string `json:"pdfpass" db:"pdfpass"`
	Zippass string `json:"zippass" db:"zippass"`
	Tag     string `json:"tag" db:"tag"`
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
