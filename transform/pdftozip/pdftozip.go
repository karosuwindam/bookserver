package pdftozip

import (
	"archive/zip"
	"bookserver/config"
	"bookserver/dirread"
	"bookserver/transform/pnmtojpg"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	JPG string = ".jpg"
	PNG string = ".png"
	PBM string = ".pbm"
	PPM string = ".ppm"
)

var pdfimages string = "pdfimages" //pdfからイメージを取り出すコマンド

func cmdck() error {
	if _, err := exec.Command(pdfimages, "-h").Output(); err != nil {
		return err
	}
	return nil
}

var tmpPass string //画像を一時保存するパス
var pdfPass string //pdfの参照フォルダ
var zipPass string //zipの参照フォルダ
var imgPass string //画像を保存するフォルダパス

// Setup(cfg) = error
//
// コマンド確認とフォルダ設定
//
// cfg : 設定
func SetUp(cfg *config.Config) error {
	var err error
	if err = cmdck(); err != nil {
		return err
	}
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

// removeimage(foldername) = error
//
// フォルダ内のデータ削除
//
// foldername : フォルダ名
func removeimage(foldername string) error {
	dirfolder, err := dirread.Setup(tmpPass + foldername + "/")
	if err != nil {
		return err
	}
	if err := dirfolder.Read("./"); err != nil {
		return err
	}
	var wp sync.WaitGroup
	for _, filedata := range dirfolder.Data {
		wp.Add(1)
		go func(pass string) {
			defer wp.Done()
			os.Remove(pass)
		}(filedata.RootPath + filedata.Name)
	}
	wp.Wait()
	if err := os.Remove(tmpPass + foldername); err != nil {
		return err
	}
	return nil
}

// imageCopy(filename, inputpass) = error
//
// 表紙に使用する画像ファイルを対象のフォルダにコピーする
//
// filename : コピー先のファイル名
// inputpass : コピーもとのファイル名
func imageCopy(filename, inputpass string) error {
	outname := filename
	if strings.Index(strings.ToLower(inputpass), JPG) > 0 {
		outname += JPG
	} else if strings.Index(strings.ToLower(inputpass), PNG) > 0 {
		outname += PNG
	} else {
		return errors.New("not for jpg, png")
	}
	src, err := os.Open(inputpass)
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
	return err
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

// renameData(dirfolder) = error
//
// 4桁時のリネーム処理
//
// dirfolder : 対象のフォルダのデータ
func renameData(dirfolder *dirread.Dirtype) error {
	if len(dirfolder.Data) > 1000 {
		//リネーム処理
		for _, data := range dirfolder.Data {
			newName := data.Name
			if j, err := strconv.Atoi(newName[len(newName)-8 : len(newName)-4]); err != nil || j <= 0 {
				i, _ := strconv.Atoi(newName[len(newName)-7 : len(newName)-4])
				newName = newName[:len(newName)-7] + fmt.Sprintf("%04d", i) + newName[len(newName)-4:]
			}
			if err := os.Rename(data.RootPath+data.Name, data.RootPath+newName); err != nil {
				return err
			}
		}
	}
	return nil
}

// imgToJpg(dirfolder) = error
//
// pbmやppm形式のファイルをjpgへ変換する
//
// dirfolder : 対象のフォルダのデータ
func imgToJpg(dirfolder *dirread.Dirtype) error {
	var ch chan bool = make(chan bool, 10)
	var wg sync.WaitGroup
	for _, data := range dirfolder.Data {
		go func(inputName string) {
			ch <- true
			wg.Add(1)
			if i := strings.Index(strings.ToLower(inputName), PBM); i > 0 {
				outputName := inputName[:i] + JPG
				if err := pnmtojpg.Pbm2jpg(inputName, outputName); err != nil {
					log.Println("err", inputName, outputName)
				} else {
					os.Remove(inputName)
				}
			} else if i := strings.Index(strings.ToLower(inputName), PPM); i > 0 {
				outputName := inputName[:i] + JPG
				if err := pnmtojpg.Ppm2jpg(inputName, outputName); err != nil {
					log.Println("err", inputName, outputName)
				} else {
					os.Remove(inputName)
				}
			}
			<-ch
			wg.Done()
		}(data.RootPath + data.Name)
	}
	wg.Wait()
	return nil
}

// imagetoZip(filename, outputFile) = error
//
// データのzip圧縮
//
// filename 入力ファイルから拡張子を除いたファイル名
// outputFile 出力先のファイル名
func imagetoZip(filename, outputFile string) error {
	if strings.Index(strings.ToLower(outputFile), ".zip") <= 0 || outputFile == "" {
		fmt.Println("file is not zip pass: " + outputFile)
		return nil
	}
	zipfile := zipPass + outputFile
	dirfolder, _ := dirread.Setup(tmpPass + filename + "/")
	if err := dirfolder.Read("./"); err != nil {
		return err
	}
	//4桁を超える場合はリネーム処理
	if err := renameData(dirfolder); err != nil {
		return err
	}
	dirfolder, _ = dirread.Setup(tmpPass + filename + "/")
	if err := dirfolder.Read("./"); err != nil {
		return err
	}
	//pbmファイルを全てjpgへ変換する
	imgToJpg(dirfolder)
	dirfolder, _ = dirread.Setup(tmpPass + filename + "/")
	if err := dirfolder.Read("./"); err != nil {
		return err
	}
	dest, err := os.Create(zipfile)
	if err != nil {
		return err
	}
	zipWrite := zip.NewWriter(dest)
	defer zipWrite.Close()
	imgflag := false
	for _, file := range dirfolder.Data {
		if file.Folder {
			continue
		}
		if !imgflag {
			if err := imageCopy(filename, file.RootPath+file.Name); err != nil {
				log.Println(err)
			} else {
				imgflag = true
			}
		}
		if err := addToZip(file.RootPath+file.Name, zipWrite); err != nil {
			return err
		}
	}
	fmt.Println("create zip file to", zipfile)
	return nil
}

// Pdftoimage(inputFile, outputFIle) = error
//
// pdfファイルから画像ファイルを取り出して、zipファイルを作る
//
// imputFile : 入力ファイル名
// outputFile : 出力ファイル名
func Pdftoimage(inputFile, outputFile string) error {
	if i := strings.Index(strings.ToLower(inputFile), ".pdf"); i > 0 {
		filename := inputFile[:i]
		cmdArry := []string{ //pdfimageコマンドのオプションを生成
			pdfPass + inputFile,
			tmpPass + filename + "/" + filename,
			"-j",
		}
		// 一時保存用のフォルダを作成
		if err := os.MkdirAll(tmpPass+filename, 0777); err != nil {
			return err
		}
		//pdfimageコマンドの実行
		if err := exec.Command(pdfimages, cmdArry...).Run(); err != nil {
			//失敗
			return err
		} else {
			//成功
			//zipファイルの作成
			if err := imagetoZip(filename, outputFile); err != nil { //失敗時の処理
				removeimage((filename))
				return err
			}
			//pdfimagesのファイル削除
			if err = removeimage(filename); err != nil {
				return err
			}
		}

	}
	return nil
}
