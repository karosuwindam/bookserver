package table

type booknames struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Title    string `json:"title" db:"title"`
	Writer   string `json:"writer" db:"writer"`
	Brand    string `json:"brand" db:"brand"`
	Booktype string `json:"booktype" db:"booktype"`
	Ext      string `json:"ext" db:"ext"`
}

type copyfile struct {
	Id       int    `json:"id" db:"id"`
	Zippass  string `json:"zippass" db:"zippass"`
	Filesize int    `json:"filesize" db:"filesize"`
	Copyflag int    `json:"copyflag" db:"copyflag"`
}

type filelists struct {
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

type readBase struct {
	pData interface{}
}

var tablelist map[string]interface{}

func tablelistsetup() {
	tablelist = map[string]interface{}{}
	tablelist[BOOKNAME] = booknames{}
	tablelist[COPYFILE] = copyfile{}
	tablelist[FILELIST] = filelists{}
	return
}

func jsonconvertstruct(table, json string) interface{} {

	return nil
}

func readBaseCreate(tname string) readBase {
	var out readBase
	switch tname {
	case BOOKNAME:
		tmp := []booknames{}
		out.pData = tmp
	case COPYFILE:
		tmp := []copyfile{}
		out.pData = tmp
	case FILELIST:
		tmp := []filelists{}
		out.pData = tmp
	}

	return out
}
