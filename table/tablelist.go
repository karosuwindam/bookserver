package table

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

// ckType(interface{}) = bool
//
// 変数の型の確認
//
// a(interface{}) : 型を代入
func ckType(a interface{}) bool {
	switch a.(type) {
	case *Booknames, *Filelists, *Copyfile:
		return true

	}
	return false
}
