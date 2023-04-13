package view

import (
	"archive/zip"
	"bookserver/api/common"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type ZipFile struct {
	path      string   //zipファイルパス
	flag      bool     //有効の有無
	DataName  []string `json:"Name"`  //zip内のデータファイル名
	FileCount int      `json:"Count"` //zip内のファイル数
}

// zipファイルを開く準備
func openfile(name string) ZipFile {
	output := ZipFile{}
	if strings.Index(strings.ToLower(name), ".zip") <= 0 {
		return output
	}
	if !common.Exists(zippath + name) {
		return output
	}
	output.path = zippath + name
	if err := output.createfilelist(); err != nil {
		return output
	}
	output.flag = true

	return output
}

// zipファイルからファイルリストを用意する
func (f *ZipFile) createfilelist() error {
	if f.path == "" {
		return errors.New("path is not input")
	}
	if r, err := zip.OpenReader(f.path); err != nil {
		return err
	} else {
		defer r.Close()
		i := 0
		tmp := []string{}
		for _, f := range r.File {
			tmp = append(tmp, f.Name)
			i++
		}
		f.FileCount = i
		f.DataName = tmp
	}

	return nil
}

// zipファイル内からファイル名で開く
func (f *ZipFile) openZipRead(name string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if !f.flag {
		return buf, errors.New("not chack file")
	}
	if r, err := zip.OpenReader(f.path); err != nil {
		return buf, err
	} else {
		defer r.Close()
		if rc, err := r.Open(name); err != nil {
			return buf, err
		} else {
			_, err := io.Copy(buf, rc)
			if err != nil {
				return buf, err
			}
			rc.Close()
		}
	}
	return buf, nil
}

func (f *ZipFile) convertjson() string {
	json, _ := json.Marshal(f)
	return string(json)
}
