package ziptopdf

import (
	"bookserver/config"
	"fmt"
	"image"
	"os"

	"github.com/signintech/gopdf"
)

var tmpPass string //画像を一時保存するパス
var pdfPass string //pdfの参照フォルダ
var zipPass string //zipの参照フォルダ
var imgPass string //画像を保存するフォルダパス

// Setup(cfg) = error
//
// コマンド確認とフォルダ設定
func SetUp(cfg *config.Config) error {
	var err error
	tmpPass = cfg.Folder.Tmp
	if tmpPass[len(tmpPass)-1:] != "/" {
		tmpPass += "/"
	}
	pdfPass = cfg.Folder.Pdf
	if pdfPass[len(pdfPass)-1:] != "/" {
		pdfPass += "/"
	}
	zipPass = cfg.Folder.Zip
	if zipPass[len(zipPass)-1:] != "/" {
		zipPass += "/"
	}
	imgPass = cfg.Folder.Img
	if imgPass[len(imgPass)-1:] != "/" {
		imgPass += "/"
	}
	if err := os.MkdirAll(tmpPass, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(pdfPass, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(zipPass, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(imgPass, 0777); err != nil {
		return err
	}
	return err
}

// ZipToPdf(filename, outputFile) = error
//
// zipからPDFファイル作成
// ToDo
func ZipToPdf(filename, outputFIle string) error {
	//一時フォルダにzipファイルを解凍する
	//フォルダ内のファイルを取得する
	//フォルダ内のファイルを調べる
	//pdfファイルを作成する
	//一時フォルダ内のファイルを削除する
	return nil
}

// 一時フォルダにzipファイルを解凍する
// ToDo

// img2pdf(imgFiles pdfFile) = error
//
// image(png,jpg)ファイルからpdfファイルを作成する
//
// imgFiles : 画像ファイルパスの配列
// pdfFile : pdf出力先のファイルパス
func img2pdf(imgFiles []string, pdfFile string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for _, imgFile := range imgFiles {
		pdf.AddPage()
		//画像ファイルを読み込む
		f, err := os.Open(imgFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer f.Close()
		_, format, err := image.DecodeConfig(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//フォーマットチェック
		if format == "jpeg" || format == "jpg" || format == "png" {
			pdf.Image(imgFile, 0, 0, gopdf.PageSizeA4)
		}
	}
	return pdf.WritePdf(pdfFile)
}
