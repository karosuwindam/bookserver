package pdftozip

import (
	"archive/zip"
	"bookserver/controller/convert/pnmtojpg"
	"bookserver/table/booknames"
	"bookserver/table/filelists"
	"io"
	"io/ioutil"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// pdfのファイルをfilelistへ登録する予定の情報を取得
// 末尾の数字を巻数として扱うための情報取得処理
func ConvertPdfToZipChack(filepath string) (filelists.Filelists, error) {
	var output = filelists.Filelists{}
	pass := ""
	count := 0

	if i := strings.Index(strings.ToLower(filepath), ".pdf"); i > 0 {
		output.Pdfpass = filepath
		output.Name = filepath[:i]
		// pass =
		tmp := ""
		for j := len(output.Name); j > 0; j-- {

			tmp = output.Name[j-1:]

			if z, err := strconv.Atoi(tmp); err != nil {
				pass = output.Name[:j]
				break
			} else {
				count = z
			}
		}
	} else {
		return output, errors.New("input File name not pdf")
	}
	//zipのファイル名を作るための情報をデータベースから検索
	d, err := booknames.GetName(output.Name)
	if err == nil {
		tmpAry := []string{}
		if d.Title != "" {
			output.Zippass = d.Title + ".zip"
			tmpAry = append(tmpAry, d.Title)
		}
		if d.Writer != "" {
			output.Zippass = "[" + d.Writer + "]" + output.Zippass
			tmpAry = append(tmpAry, d.Writer)
		}
		if d.Burand != "" {
			tmpAry = append(tmpAry, d.Burand)
		}
		if d.Booktype != "" {
			tmpAry = append(tmpAry, d.Booktype)
		}
		if d.Ext != "" {
			tmpAry = append(tmpAry, d.Ext)
		}
		for i, tt := range tmpAry {
			if i == 0 {
				output.Tag = tt
			} else {
				output.Tag += "," + tt
			}
		}
	} else if d, err = booknames.GetName(pass); err == nil {
		tmpAry := []string{}
		if d.Title != "" {
			output.Zippass = d.Title
			if tmp := strconv.Itoa(count); len(tmp) == 1 {
				output.Zippass = output.Zippass + "0" + tmp
			} else {
				output.Zippass = output.Zippass + tmp

			}
			tmpAry = append(tmpAry, output.Zippass)
			output.Zippass = output.Zippass + ".zip"
		}
		if d.Writer != "" {
			output.Zippass = "[" + d.Writer + "]" + output.Zippass
			tmpAry = append(tmpAry, d.Writer)
		}
		if d.Burand != "" {
			tmpAry = append(tmpAry, d.Burand)
		}
		if d.Booktype != "" {
			tmpAry = append(tmpAry, d.Booktype)
		}
		if d.Ext != "" {
			tmpAry = append(tmpAry, d.Ext)
		}
		for i, tt := range tmpAry {
			if i == 0 {
				output.Tag = tt
			} else {
				output.Tag += "," + tt
			}
		}

	} else {
		slog.Info("sql booknames not serch", "name", pass, "error", err.Error())
		output.Zippass = output.Name + ".zip"
	}
	return output, nil
}

// pdfのファイルを読み取ったときに画像を取り出してzipへ変換
func ConvertPdfToZip(filepath string) error {

	if i := strings.Index(strings.ToLower(filepath), ".pdf"); i <= 0 {
		return errors.New("input File name not pdf")
	}
	defer func() {

		//一時フォルダ内のファイルを削除
		if err := removeTmpFileFolder(filepath); err != nil {
			slog.Error("removeTmpFileFolder", "error", err.Error())
		}
	}()
	//pdfimageコマンドを使ってtmpフォルダへpdfファイルから画像と取り出す
	if err := outputImageByPdf(filepath); err != nil {
		return errors.Wrap(err, "pdf to image covert error ")
	}
	//画像ファイルでjpg以外の画像をjpgへ変換
	if err := imgToJpg(filepath); err != nil {
		return errors.Wrap(err, "Not Open file")
	}
	//一番初めに取り出した画像ファイルを特定フォルダへコピ
	if err := imageCopyToJpg(filepath); err != nil {
		return errors.Wrap(err, "imageCopyToJpg")
	}
	//zipのファイル名を作るための情報をデータベースから検索
	if b, err := ConvertPdfToZipChack(filepath); err != nil {
		return errors.Wrap(err, "Not Create table data")
	} else {
		rewriteFlag := false
		if _, err := os.Stat(zipPass + b.Zippass); err == nil {
			rewriteFlag = true
		}
		//zipファイルを作成する
		if err := createZipfile(filepath, b.Zippass); err != nil {
			return errors.Wrap(err, "Not Create Zip FIle")
		}
		//データベース正式登録
		if rewriteFlag { //既存ファイルがある場合はデータベースにも存在している可能性があるので
			if fileData, err := filelists.GetName(b.Name); err == nil && fileData.Name == b.Name {
				fileData.Pdfpass = b.Pdfpass
				fileData.Zippass = b.Zippass
				fileData.Tag = b.Tag
				if err1 := fileData.Update(); err1 != nil {
					return errors.Wrap(err1, "Not Chanage table data for filelists")
				}
				return nil
			}
		}
		if err := b.Add(); err != nil {
			return errors.Wrap(err, "Not Add table for filelists")
		}
	}

	return nil
}

func createZipfile(filename, outputName string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
		pass = filename[:i]
	}
	if pass[len(pass)-1:] != "/" {
		pass += "/"
	}
	slog.Info("Create Zip FIle", "zipPass", zipPass, "outputName", outputName)
	dest, err := os.Create(zipPass + outputName)
	if err != nil {
		return err
	}
	zipWrite := zip.NewWriter(dest)
	defer zipWrite.Close()
	files, err := ioutil.ReadDir(tmpPass + pass)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			addToZip(tmpPass+pass+f.Name(), zipWrite)
		}
	}
	return nil
}

// addToZip(filename zipWriter) = error
//
// zipへファイルを追加
//
// filename 追加するファイル名
// zipWriter 書き込み対象となるファイルポイント
func addToZip(filename string, zipWriter *zip.Writer) error {
	info, _ := os.Stat(filename)
	hdr, _ := zip.FileInfoHeader(info)
	hdr.Name = filename
	for _, s := range strings.Split(filename, "/") {
		hdr.Name = s
	}
	f, err := zipWriter.CreateHeader(hdr)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	f.Write(body)
	return nil
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
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
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
					slog.Warn("Pbm2jpg", "error", err.Error())
				} else {
					slog.Debug("Covert file", "filepass", filepass, "outputName", outputName)
					os.Remove(filepass)
				}
			} else if i := strings.Index(strings.ToLower(filepass), PPM); i > 0 {
				outputName := filepass[:i] + JPG
				if err := pnmtojpg.Ppm2jpg(filepass, outputName); err != nil {
					slog.Warn("Ppm2jpg", "error", err.Error())
				} else {
					slog.Debug("Covert file", "filepass", filepass, "outputName", outputName)
					os.Remove(filepass)
				}
			}
			<-ch
		}(tmpPass + pass + f.Name())
	}
	wg.Wait()
	return nil
}

// 番初めに取り出した画像ファイルを特定フォルダへコピ
func imageCopyToJpg(filename string) error {

	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
		pass = filename[:i]
	}
	outname := pass
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
			slog.Debug("copy img", "imgPass", imgPass+outname)
			return err
		}
	}
	return errors.New("not for jpg, png")
}

// pdfimageコマンドを使ってtmpフォルダへpdfファイルから画像を取り出す
func outputImageByPdf(filename string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
		pass = filename[:i]
	}
	cmdArry := []string{ //pdfimageコマンドのオプションを生成
		pdfPass + filename,
		tmpPass + pass + "/" + pass,
		"-j",
	}
	//一時保存用のフォルダを作成
	if err := os.MkdirAll(tmpPass+pass, 0777); err != nil {
		return err
	}
	//pdfimageコマンドの実行
	if err := exec.Command(pdfimages, cmdArry...).Run(); err != nil {
		if err1 := os.RemoveAll(tmpPass + pass); err1 != nil {
			return err1
		}
		return err
	}
	return nil
}

// 作成した一時フォルダを削除する
func removeTmpFileFolder(filename string) error {
	pass := filename
	if i := strings.Index(strings.ToLower(filename), ".pdf"); i > 0 {
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
