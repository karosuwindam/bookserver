package createzip

import (
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

func (t *PdftoZip) pdftoimage() {
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

		}
	}
}
