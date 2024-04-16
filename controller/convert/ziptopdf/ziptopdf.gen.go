package ziptopdf

import (
	"bookserver/config"
	"os"

	"github.com/pkg/errors"
)

var tmpPass string //画像を一時保存するパス
var pdfPass string //pdfの参照フォルダ
var zipPass string //zipの参照フォルダ
var imgPass string //画像を保存するフォルダパス

func Init() error {
	//フォルダパスの設定
	if err := setupPath(); err != nil {
		return errors.Wrap(err, "Error Setup Folder Path")
	}
	//フォルダの作成
	if err := createPath(); err != nil {
		return errors.Wrap(err, "Error Create Folder Path")
	}
	return nil
}

// 各フォルダパスを設定する
func setupPath() error {
	tmpPass = config.BScfg.Tmp
	if tmpPass[len(tmpPass)-1:] != "/" {
		tmpPass += "/"
	}
	pdfPass = config.BScfg.Pdf
	if pdfPass[len(pdfPass)-1:] != "/" {
		pdfPass += "/"
	}
	zipPass = config.BScfg.Zip
	if zipPass[len(zipPass)-1:] != "/" {
		zipPass += "/"
	}
	imgPass = config.BScfg.Img
	if imgPass[len(imgPass)-1:] != "/" {
		imgPass += "/"
	}

	return nil
}

// 各フォルダパスに対してフォルダを作成する
func createPath() error {
	if err := os.MkdirAll(tmpPass, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(zipPass, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(pdfPass, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(imgPass, 0777); err != nil {
		return err
	}
	return nil
}
