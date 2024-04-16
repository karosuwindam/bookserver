package tabledelete

import "errors"

const (
	BOOKNAMES string = "booknames"
	COPYFILES string = "copyfiles"
	FILELISTS string = "filelists"
)

var errorsTableData error = errors.New("table not fond")

var tables []string = []string{
	BOOKNAMES,
	COPYFILES,
	FILELISTS,
}

// 有効なテーブルであることを確認
func checkTableData(table string) error {
	for _, t := range tables {
		if table == t {
			return nil
		}
	}
	return errorsTableData
}
