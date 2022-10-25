package createzip

import (
	"archive/zip"
	"bookserver/dirread"
	"bookserver/message"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type PdftoZip struct {
	Pdfpass string
	Zippass string
	PdfName string
	ZipName string
	Tmppass string
}

//setup
func Setup(pdfpass, pdfname, zippass, tmppass string) (*PdftoZip, error) {
	t := &PdftoZip{Pdfpass: pdfpass, PdfName: pdfname, Zippass: zippass, Tmppass: tmppass}
	return t, nil
}

//Zipファイルの名前作成
func readPdftoZipName(pdfname string) string {
	output := ""
	return output
}

//PDFのフルパス作成
func (t *PdftoZip) pdffullpass() string {
	output := t.Pdfpass
	if output[len(output)-1] != "/"[0] {
		output += "/"
	}
	output += t.PdfName
	return output
}

//TMPフォルダ作成
func (t *PdftoZip) createtmpfolder() string {
	tmppass := t.Tmppass
	if tmppass[len(tmppass)-1] != "/"[0] {
		tmppass += "/"
	}
	tmppass += t.PdfName[:len(t.PdfName)-4]
	os.MkdirAll(tmppass, 0777)
	return tmppass
}

//ZIPのフルパス作成
func (t *PdftoZip) zipfullpass() string {
	output := t.Zippass
	if output[len(output)-1] != "/"[0] {
		output += "/"
	}
	output += t.ZipName
	return output

}

//フォルダ内のデータ削除
func (t *PdftoZip) removeimage() error {
	dirfolder, err := dirread.Setup(t.Tmppass)
	if err != nil {
		return err
	}
	for _, filedata := range dirfolder.Data {
		os.Remove(t.Tmppass + filedata.Name)
	}
	if err := os.Remove(t.Tmppass); err != nil {
		return err
	}
	return nil
}

//データのzip圧縮
func (t *PdftoZip) imagetoZip() error {
	zipfile := t.zipfullpass()
	dirfolder, _ := dirread.Setup(t.Tmppass)
	if err := dirfolder.Read("./"); err != nil {
		return err
	}
	dest, err := os.Create(zipfile)
	if err != nil {
		return err
	}
	zipWrite := zip.NewWriter(dest)
	defer zipWrite.Close()
	for _, file := range dirfolder.Data {
		if err := addToZip(t.Tmppass+file.Name, zipWrite); err != nil {
			return err
		}
	}
	message.Println("create zip file to", zipfile)
	return nil
}

//zipへファイルを追加
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

func (t *PdftoZip) Pdftoimage() {
	str := t.PdfName
	if strings.Index(str, "pdf") > 0 {
		filename := str[:len(t.PdfName)-4]
		cmdArry := []string{"pdfimages", t.pdffullpass(), t.createtmpfolder() + filename, "-j"}
		subcmd := ""
		for i, str := range cmdArry {
			if i != 0 {
				subcmd += " "
			}
			subcmd += str
		}
		err := exec.Command("sh", "-c", subcmd).Run()
		if err != nil { //作成失敗

		} else { //作成成功
			//zipファイルの作成
			if err := t.imagetoZip(); err != nil { //失敗時の処理

			}
			//pdfimagesでできたファイルの削除
			if err := t.removeimage(); err != nil {

			}
		}
	}
}
