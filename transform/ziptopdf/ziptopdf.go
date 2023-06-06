package ziptopdf

import (
	"archive/zip"
	"bookserver/config"
	"bookserver/transform/writetable"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/signintech/gopdf"
)

var tmpPass string //画像を一時保存するパス
var pdfPass string //pdfの参照フォルダ
var zipPass string //zipの参照フォルダ
var imgPass string //画像を保存するフォルダパス

// zipフォルダの管理
type zipFolder struct {
	zipFileName string
	outFolder   string
	ziplist     []string
}

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

// ZipToPdf(t) = error
//
// zipからPDFファイル作成
//
// t(writetable.ZipToPdf) : 作成テーブル
func ZipToPdf(t writetable.ZipToPdf) error {
	//一時フォルダにzipファイルを解凍する
	z, err := unzip(t.InputFile)
	if err != nil {
		return err
	}
	//一時フォルダ内のファイルを削除する
	defer z.removeFolder()
	//イメージファイルから表紙を作成
	f, err := z.imgCopy(t.Name)
	if err != nil {
		return err
	}
	//pdfファイルを作成する
	if err := z.img2pdf(t.OutputFile); err != nil {
		//ファイル作成に失敗したら表紙を削除
		os.Remove(imgPass + f)
		return err
	}
	return nil
}

// unzip(zipFile) = zipFolder ,error
//
// 一時フォルダにzipファイルを解凍する
//
// zipFile(string) : zipフォルダ内のあるzipファイル名
func unzip(zipFile string) (zipFolder, error) {
	var tmpZipFolder string
	output := zipFolder{zipFileName: zipFile, ziplist: []string{}}
	// zipファイル名から一時フォルダにフォルダを作成する
	if len(zipFile) > 4 {
		tmpZipFolder = tmpPass + zipFile[:len(zipFile)-4] + "/"
		output.outFolder = tmpZipFolder
	} else {
		return output, fmt.Errorf("zipファイル名が不正です。")
	}

	//出力フォルダの存在チェック
	if _, err := os.Stat(tmpZipFolder); err != nil {
		// 出力先フォルダを作成する
		if err := os.Mkdir(tmpZipFolder, 0777); err != nil {
			return output, err
		}
	}
	// zipファイルを開く
	r, err := zip.OpenReader(zipPass + zipFile)
	if err != nil {
		return output, err
	}
	defer r.Close()

	// zipファイル内のファイルを展開する
	for _, f := range r.File {
		// 出力先ファイルを作成する
		outFile, err := os.Create(tmpZipFolder + f.Name)
		output.ziplist = append(output.ziplist, f.Name)
		if err != nil {
			return output, err
		}
		defer outFile.Close()
		// zipファイル内のファイルを読み込む
		rc, err := f.Open()
		if err != nil {
			return output, err
		}
		defer rc.Close()
		// 出力先ファイルに書き込む
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

// (t *zipFolder) removeFolder() = error
//
// zipが作成した一時フォルダを削除する
func (t *zipFolder) removeFolder() error {
	err := os.RemoveAll(t.outFolder)
	t.outFolder = ""
	t.zipFileName = ""
	t.ziplist = []string{}
	return err
}

// (t *zipFolder) imgCopy(imgName) = string, error
//
// 1page目のjpgファイルをイメージフォルダにコピーする
//
// imgName(string) : 拡張子がないファイル名
func (t *zipFolder) imgCopy(imgName string) (string, error) {
	for _, imgFile := range t.ziplist {
		//画像ファイルを読み込む
		f, err := os.Open(t.outFolder + imgFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, format, err := image.DecodeConfig(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		f.Close()
		//フォーマットチェック
		fileName := ""
		if format == "jpeg" || format == "jpg" {
			fileName = imgPass + imgName + ".jpg"
		} else if format == "png" {
			fileName = imgPass + imgName + ".png"
		}
		if fileName != "" {

			f, _ := os.Open(t.outFolder + imgFile)
			defer f.Close()
			out, err := os.Create(fileName)
			if err != nil {
				return imgFile, err
			}
			defer out.Close()
			_, err = io.Copy(out, f)

			return imgFile, err
		}
	}
	return "", fmt.Errorf("not image file")
}

// img2pdf(imgFiles pdfFile) = error
//
// image(png,jpg)ファイルからpdfファイルを作成する
//
// pdfFile : pdf出力先のファイルパス
// ToDo: サイズが固定なのでこの部分を後で修正する
func (t *zipFolder) img2pdf(pdfFile string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for _, imgFile := range t.ziplist {
		pdf.AddPage()
		//画像ファイルを読み込む
		f, err := os.Open(t.outFolder + imgFile)
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
			pdf.Image(t.outFolder+imgFile, 0, 0, gopdf.PageSizeA4)
		}
	}
	return pdf.WritePdf(pdfPass + pdfFile)
}
