package fileupload

import "bookserver/table"

func (t *UploadPass) ReadDb(str string) (string, error) {
	return t.Sql.ReadName(table.BOOKNAME, str)
}
