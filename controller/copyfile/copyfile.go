package copyfile

import (
	"bookserver/table/copyfiles"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

type CopyFIleData struct {
	id   int
	flag bool
}

// idと状態を指定してその状態をテーブルに登録する
func (t *CopyFIleData) AddTable() error {
	tmp := copyfiles.Copyfile{}
	if d, err := ReadCopyFIleFlagById(t.id); err != nil {
		return err
	} else {
		tmp = d.Copyfile
	}
	//Zipファイルがあるか確認するか確認
	if i := checkZipFileSize(tmp.Zippass); i != 0 {
		tmp.Filesize = i
	} else {
		return errors.New(fmt.Sprintf("file is not %v", tmp.Zippass))
	}
	if t.flag {
		//ファイルのコピー処理
		if err := copyFileForZipToPublic(tmp.Zippass); err != nil {
			return errors.New(fmt.Sprintf("not Copy file %v", tmp.Zippass))
		}
		tmp.Copyflag = copyfiles.ON
	} else {
		if err := removeFileFromPublic(tmp.Zippass); err != nil {
			return errors.New(fmt.Sprintf("not Remove file %v", tmp.Zippass))
		}
		//ファイルの削除処理
		tmp.Copyflag = copyfiles.OFF
	}
	if tmp.Id != 0 {
		if err := tmp.Update(); err != nil {
			return err
		}
	} else {
		if err := tmp.Add(); err != nil {
			return err
		}
	}
	return nil
}

// ファイル名を指定してファイルサイズを確認
func checkZipFileSize(str string) int {
	psss := zippass + str
	if f, err := os.Stat(psss); err != nil {
		return 0
	} else {
		return int(f.Size())
	}
}

// ファイル名を指定して公開フォルダへコピを実施
func copyFileForZipToPublic(str string) error {
	zipFilePass := zippass + str
	publicFilePass := publicpass + str
	if _, err := os.Stat(zipFilePass); err != nil {
		return errors.Wrap(err, fmt.Sprintf("os.Stat(%v)", zipFilePass))
	}
	if _, err := os.Stat(publicFilePass); err == nil {
		return nil
	}
	f, err := os.Open(zipFilePass)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("os.Open(%v)", zipFilePass))
	}
	defer f.Close()
	fp, err := os.Create(publicFilePass)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("os.Create(%v)", publicFilePass))
	}
	defer fp.Close()

	if _, err := io.Copy(fp, f); err != nil {
		return errors.Wrap(err, fmt.Sprintf("io.Copy(%v,%v)", publicFilePass, zipFilePass))
	}
	log.Println("info:", "Copyfile", str)
	return nil

}

// ファイル名を指定して公開フォルダから削除
func removeFileFromPublic(str string) error {
	publicFilePass := publicpass + str
	if _, err := os.Stat(publicFilePass); err != nil {
		return nil
	}
	if err := os.Remove(publicFilePass); err != nil {
		return errors.Wrap(err, fmt.Sprintf("os.Remove(%v)", publicFilePass))
	}
	log.Println("info:", "Remove file", str)

	return nil
}

// IDに紐付いたファイルが公開フォルダにあるか確認を行いファイルがない場合は削除
func ChackCopyFileTableDataAll() error {
	//テーブル内で有効になているものをすべて取得
	if list, err := copyfiles.OnOFFSearch(copyfiles.ON); err != nil {
		return err
	} else {
		for _, d := range list {
			//公開フォルダにファイルの存在を確認
			publicFilePass := publicpass + d.Zippass
			if _, err := os.Stat(publicFilePass); err == nil {
				continue
			}
			log.Println("info:", fmt.Sprintln("Not file %v"), publicFilePass)
			//ファイルが存在しないものは無効にする
			if err := d.OFF(); err != nil {
				log.Println("error:", err)
			}
		}
	}
	return nil
}
