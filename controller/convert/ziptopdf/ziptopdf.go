package ziptopdf

import (
	"archive/zip"
	"bookserver/controller/convert/pnmtojpg"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

func ConvertZipToPdf(filepath string) error {
	// outputname := ""
	pass := ""
	// tag := ""

	if i := strings.Index(strings.ToLower(filepath), ".zip"); i > 0 {
		pass = filepath[:i]
		// outputname = pass + ".pdf"
	} else {
		return errors.New("input File name not pdf")
	}
	defer func() {

		//一時フォルダ内のファイルを削除
		if err := removeTmpFileFolder(filepath); err != nil {
			log.Println("error:", err)
		}
	}()
	//一時フォルダにzipファイルを解凍する
	if err := unzip(filepath); err != nil {
		return errors.Wrap(err, "unzip error")
	}
	//画像ファイルでjpg以外の画像をjpgへ変換
	if err := imgToJpg(filepath); err != nil {
		return errors.Wrap(err, "img To Jpg")
	}
	outputname := pass //表示画像として使用するファイル名を指定
	//一時ファイルの際の画像をコピー
	if err := imageCopyToJpg(filepath, outputname); err != nil {
		return errors.Wrap(err, "CopyImageName")
	}
	pdfName := outputname + ".pdf"
	//pdfファイルを作成
	if err := imgToPdf(filepath, pdfName); err != nil {
		return errors.Wrap(err, "Img To Pdf")
	}
	//データベースへ正式登録
	//:ToDo

	return nil
}

// 一時フォルダ内の画像ファイルをpdfへ変換
func imgToPdf(filename, pdfname string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".zip"); i > 0 {
		pass = filename[:i]
	}
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	files, err := ioutil.ReadDir(tmpPass + pass)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := os.Open(tmpPass + pass + file.Name())
		if err != nil {
			log.Println("error:", err)
			continue
		}
		defer f.Close()
		img, format, err := image.DecodeConfig(f)
		if err != nil {
			log.Println("error:", err)
		}
		if format == "jpeg" || format == "jpg" || format == "png" {
			//画像ファイルのサイズを取得する
			width := float64(img.Width)
			height := float64(img.Height)
			//サイズに応じてページサイズを変更する
			rect := &gopdf.Rect{H: height, W: width}
			rect = rect.UnitsToPoints(gopdf.UnitPT)
			pdf.AddPageWithOption(gopdf.PageOption{PageSize: rect})
			pdf.Image(tmpPass+pass+file.Name(), 0, 0, rect)

		}

	}
	return pdf.WritePdf(pdfPass + pdfname)
}

const (
	JPG string = ".jpg"
	PNG string = ".png"
	PBM string = ".pbm"
	PPM string = ".ppm"
)

// pbmやppm形式のファイルをjpgへ変換する
func imgToJpg(filename string) error {
	var ch chan bool = make(chan bool, 10)
	var wg sync.WaitGroup

	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".zip"); i > 0 {
		pass = filename[:i]
	}
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}

	files, err := ioutil.ReadDir(tmpPass + pass)
	if err != nil {
		return err
	}
	for _, f := range files {
		wg.Add(1)
		go func(filepass string) {
			defer wg.Done()
			ch <- true
			if i := strings.Index(strings.ToLower(filepass), PBM); i > 0 {
				outputName := filepass[:i] + JPG
				if err := pnmtojpg.Pbm2jpg(filepass, outputName); err != nil {
					log.Println("error:", filepass, outputName)
				} else {
					log.Panicln("debug:", "Covert file", filepass, "to", outputName)
					os.Remove(filepass)
				}
			} else if i := strings.Index(strings.ToLower(filepass), PPM); i > 0 {
				outputName := filepass[:i] + JPG
				if err := pnmtojpg.Ppm2jpg(filepass, outputName); err != nil {
					log.Println("error:", filepass, outputName)
				} else {
					log.Panicln("debug:", "Covert file", filepass, "to", outputName)
					os.Remove(filepass)
				}
			}
			<-ch
		}(tmpPass + pass + f.Name())
	}
	wg.Wait()
	return nil
}

// unzip(filename) = error
//
// 一時フォルダにzipファイルを解凍する
//
// filename(string) : zipフォルダ内のあるzipファイル名
func unzip(filename string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".zip"); i > 0 {
		pass = filename[:i]
	}
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	r, err := zip.OpenReader(zipPass + filename)
	if err != nil {
		return errors.Wrap(err, "zip file open error")
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		outFile, err := os.Create(tmpPass + pass + f.Name)
		defer outFile.Close()
		rc, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "open file error")
		}
		defer rc.Close()
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return errors.Wrap(err, "copy file error")
		}

	}

	return nil
}

// 作成した一時フォルダを削除する
func removeTmpFileFolder(filename string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".zip"); i > 0 {
		pass = filename[:i]
	}
	if f, err := os.Stat(tmpPass + pass); os.IsExist(err) || !f.IsDir() {
		return nil
	}
	if err := os.RemoveAll(tmpPass + pass); err != nil {
		return err
	}
	return nil
}

// 番初めに取り出した画像ファイルを特定フォルダへコピ
func imageCopyToJpg(filename, outputname string) error {

	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
		pass = filename[:i]
	}
	outname := outputname
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}

	files, err := ioutil.ReadDir(tmpPass + pass)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {

			if strings.Index(strings.ToLower(f.Name()), JPG) > 0 {
				outname += JPG
			} else if strings.Index(strings.ToLower(f.Name()), PNG) > 0 {
				outname += PNG
			} else {
				continue
			}
			src, err := os.Open(tmpPass + pass + f.Name())
			if err != nil {
				return err
			}
			defer src.Close()
			dst, err := os.Create(imgPass + outname)
			if err != nil {
				return err
			}
			defer dst.Close()
			_, err = io.Copy(dst, src)
			log.Println("debug:", "copy img", imgPass+outname)
			return err
		}
	}
	return errors.New("not for jpg, png")
}
