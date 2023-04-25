package copycopyfile

import (
	"bookserver/config"
	"errors"
	"fmt"
	"io"
	"os"
)

var zipPass string    //コピー前のzipフォルダ
var pdfPass string    //コピー前のpdfフォルダ
var publicPass string //コピー先の公開フォルダ

// RemoveFile(filename) = error
//
// 指定ファイルを公開フォルダから削除
//
// filename string: 削除するファイル名
func RemoveFile(filename string) error {
	if Exists(publicPass + filename) {
		if err := os.Remove(publicPass + filename); err != nil {
			return err
		} else {
			fmt.Println("Pulic folder Delete", filename)
		}
	}
	return nil
}

// Copyfile(filename) = error
//
// zipとpdfフォルダ内にあるファイルを確認してpulicフォルダにコピーをする
//
// filename string: 各フォルダのパスと一致するファイルをコピーする
func CopyFile(filename string) error {
	var pass string
	if Exists(zipPass + filename) {
		pass = zipPass
	} else if Exists(pdfPass + filename) {
		pass = pdfPass
	} else {
		return errors.New("not file:" + filename)
	}
	file, err := os.Open(pass + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fp, err := os.Create(publicPass + filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, file); err != nil {
		return err
	}
	fmt.Println("Pulic folder Copy", filename)
	return nil
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
	if err := os.MkdirAll(publicPass, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(pdfPass, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(zipPass, 0777); err != nil {
		return err
	}
	return nil
}

// ファイルの存在確認
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
