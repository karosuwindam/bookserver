package copyfile

import (
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Copyfile struct {
	copyfiles.Copyfile
}

// idを指定してそのファイル名が別テーブルにあるか確認
func ReadCopyFIleFlagById(id int) (Copyfile, error) {
	output := Copyfile{}
	if fileData, err := filelists.GetId(id); err != nil {
		return output, errors.Wrap(err, fmt.Sprintf("filelists.GetId(%v)", id))
	} else {
		if d, err := copyfiles.GetZipName(fileData.Zippass); err == gorm.ErrRecordNotFound {
			output.Zippass = fileData.Zippass
			output.Copyflag = 0
		} else if err != nil {
			return output, errors.Wrap(err, fmt.Sprintf("copyfiles.GetZipName(%v)", fileData.Zippass))
		} else {
			output.Copyfile = d
		}
	}
	return output, nil
}
