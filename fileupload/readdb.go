package fileupload

import "bookserver/table"

func ReadDb(str string, sql *table.SQLStatus) {
	sql.Search(table.BOOKNAME, str)
}
