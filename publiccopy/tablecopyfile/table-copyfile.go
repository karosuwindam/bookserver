package tablecopyfile

import (
	"bookserver/config"
	"bookserver/table"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type FileType string

const (
	PDF FileType = "pdf"
	ZIP FileType = "zip"
)

var zipPass string    //コピー前のzipフォルダ
var pdfPass string    //コピー前のpdfフォルダ
var publicPass string //公開用のフォルダ

var sql *table.SQLStatus

// 公開フォルダをチェック、filecopyテーブルをの有無を切り替える
func CkCopyFileAll() error {
	t := time.Now()
	fmt.Println("copyfile table chack start")
	var wg sync.WaitGroup
	if jdata, err := sql.ReadAll(table.COPYFILE); err != nil {
		return err
	} else {
		if jdata == "[]" {
			return nil
		} else if ary, ok := table.JsonToStruct(table.COPYFILE, []byte(jdata)).([]table.Copyfile); ok {
			for _, str := range ary {
				wg.Add(1)
				go func(str table.Copyfile) {
					defer wg.Done()
					tmp := str
					if count := FIleSize(publicPass + str.Zippass); count > 0 {
						str.Copyflag = 1
					} else {
						str.Copyflag = 0
					}
					if tmp != str {
						if _, err := sql.Edit(table.COPYFILE, &str, str.Id); err == nil {
							fmt.Println(str.Zippass, "change flag to", str.Copyflag)
						}
					}
				}(str)
			}
		}
	}
	wg.Wait()
	fmt.Println("copyfile table chack end", time.Now().Sub(t).Seconds())
	return nil
}

// AddCoppyFIle(name, OnOff)= error
//
// copyfilesテーブルに登録するデータ
// name string: 登録するファイル名の指定
// OnOff bool: Copyflagを設定する
func AddCopyFIle(name string, OnOff bool) error {
	flag := 0
	if OnOff {
		flag = 1
	}
	if jdata, err := sql.ReadName(table.COPYFILE, name); err != nil {
		return err
	} else {
		if jdata == "[]" {
			if count := getFileSize(name); count > 0 {
				data := table.Copyfile{
					Zippass:  name,
					Filesize: count,
					Copyflag: flag,
				}
				if err := sql.Add(table.COPYFILE, &data); err != nil {
					return err
				}
			} else {
				return errors.New("file not fond " + name)
			}
		} else {
			if t, ok := table.JsonToStruct(table.COPYFILE, []byte(jdata)).([]table.Copyfile); ok {
				t[0].Copyflag = flag

				if count := getFileSize(name); count > 0 {
					t[0].Filesize = count
					if _, err := sql.Edit(table.COPYFILE, &t[0], t[0].Id); err != nil {
						return err
					}
				} else {
					if _, err := sql.Delete(table.COPYFILE, t[0].Id); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// GetTableByName(id filetype) (string, error)
//
// # Filelistsテーブルからidと種類を指定してファイル名を引き出す
//
// id int : 取り出し用のid
// filetype FileType: 取り出すための種類
func GetTableByName(id int, filetype FileType) (string, error) {
	if jdata, err := sql.ReadID(table.FILELIST, id); err != nil {
		return "", err
	} else {
		if jdata == "[]" {
			return "", errors.New("input id not search for FileList by ID:" + strconv.Itoa(id))
		} else {
			if filearr, ok := table.JsonToStruct(table.FILELIST, []byte(jdata)).([]table.Filelists); ok {
				switch filetype {
				case PDF:
					return filearr[0].Pdfpass, nil
				case ZIP:
					return filearr[0].Zippass, nil
				default:
					return filearr[0].Zippass, nil
				}

			}
		}
	}
	return "", nil
}

// セットアップ
func Setup(cfg *config.Config) error {

	publicPass = cfg.Folder.Public
	if publicPass[len(publicPass)-1:] != "/" {
		publicPass += "/"
	}
	pdfPass = cfg.Folder.Pdf
	if pdfPass[len(pdfPass)-1:] != "/" {
		pdfPass += "/"
	}
	zipPass = cfg.Folder.Zip
	if zipPass[len(zipPass)-1:] != "/" {
		zipPass += "/"
	}
	if scfg, err := table.Setup(cfg); err != nil {
		return err
	} else {
		sql = scfg
	}
	return nil
}

// ファイルサイズの確認
func FIleSize(filename string) int {
	size := 0
	if fileinfo, err := os.Stat(filename); err != nil {
		return -1
	} else {
		size = int(fileinfo.Size())
	}
	return size
}

// getFileSize(name) = int
//
// zipとpdfフォルダから対象の名前のファイルのサイズを探す
// name string:対象のファイル名
func getFileSize(name string) int {
	if count := FIleSize(pdfPass + name); count > 0 {
		return count
	} else if count := FIleSize(zipPass + name); count > 0 {
		return count
	}
	return -1
}
