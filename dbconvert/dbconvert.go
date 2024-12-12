package dbconvert

import (
	v0struct "bookserver/dbconvert/v0Struct"
	v1struct "bookserver/dbconvert/v1Struct"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var basedb *gorm.DB

func readDBFile(filepass string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepass), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	basedb = db
	return db, nil
}

func DbConvertSQL3(srcName, dstName string, srcV, dstV int) error {
	if srcV == 0 && dstV == 1 { //
		return covertSQL3_0_1(srcName, dstName)
	}
	return errors.New(fmt.Sprintf("error input data %v v=%v %v v=%v", srcName, dstName, srcV, dstV))
}

func covertSQL3_0_1(srcName, dstName string) error {
	v0db, err := gorm.Open(sqlite.Open(srcName), &gorm.Config{})
	if err != nil {
		return err
	}
	v1db, err := gorm.Open(sqlite.Open(dstName), &gorm.Config{})
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("start copy db %v to %v", srcName, dstName))

	return covertWrite0to1(v0db, v1db)
}

func covertWrite0to1(v0db, v1db *gorm.DB) error {
	timeNow := time.Now()
	//Copyfileのコピー
	slog.Info("copyfile table copy")

	if ary, err := v0struct.ReadCopyfiles(v0db); err != nil {
		return err
	} else {
		for _, d := range ary {
			tmpCopyfile := v1struct.Copyfile{
				Id:       d.Id,
				Zippass:  d.Zippass,
				Filesize: d.Filesize,
				Copyflag: d.Copyflag,
			}
			tmp := v1struct.Copyfile_sql{
				Copyfile: tmpCopyfile,
				CreateAt: nil,
				UpdateAt: nil,
			}
			var tmpCat time.Time
			if (d.CreatedAt != time.Time{}) {
				tmpCat = d.CreatedAt
				tmp.CreateAt = &tmpCat
			}
			var tmpUat time.Time
			if (d.UpdatedAt != time.Time{}) {
				tmpUat = d.UpdatedAt
				tmp.UpdateAt = &tmpUat
			}
			if err = tmp.Write(v1db); err != nil {
				return err
			}

		}
	}
	slog.Info(fmt.Sprintf("copyfile table copy end %s", time.Since(timeNow)))

	timeNow = time.Now()

	slog.Info("Booknames table copy")

	if ary, err := v0struct.ReadBooknames(v0db); err != nil {
		return err
	} else {
		for _, d := range ary {
			tmpBooknames := v1struct.Booknames{
				Id:       d.Id,
				Name:     d.Name,
				Title:    d.Title,
				Writer:   d.Writer,
				Burand:   d.Burand,
				Booktype: d.Booktype,
				Ext:      d.Ext,
			}
			tmp := v1struct.Booknames_sql{
				Booknames: tmpBooknames,
				CreateAt:  nil,
				UpdateAt:  nil,
			}
			var tmpCat time.Time
			if (d.CreatedAt != time.Time{}) {
				tmpCat = d.CreatedAt
				tmp.CreateAt = &tmpCat
			}
			var tmpUat time.Time
			if (d.UpdatedAt != time.Time{}) {
				tmpUat = d.UpdatedAt
				tmp.UpdateAt = &tmpUat
			}
			if err = tmp.Write(v1db); err != nil {
				return err
			}

		}
	}
	slog.Info(fmt.Sprintf("Booknames table copy end %s", time.Since(timeNow)))

	timeNow = time.Now()

	slog.Info("Filelists table copy")
	if ary, err := v0struct.ReadFilelists(v0db); err != nil {
		return err
	} else {
		for _, d := range ary {
			tmpFilelist := v1struct.Filelists{
				Id:      d.Id,
				Name:    d.Name,
				Pdfpass: d.Pdfpass,
				Zippass: d.Zippass,
				Tag:     d.Tag,
			}
			if tmpFilelist.Tag[len(tmpFilelist.Tag)-1:] == "," {
				tmpFilelist.Tag = tmpFilelist.Tag[:len(tmpFilelist.Tag)-1]
			}
			tmp := v1struct.Filelists_sql{
				Filelists: tmpFilelist,
				CreateAt:  nil,
				UpdateAt:  nil,
			}
			var tmpCat time.Time
			if (d.CreatedAt != time.Time{}) {
				tmpCat = d.CreatedAt
				tmp.CreateAt = &tmpCat
			}
			var tmpUat time.Time
			if (d.UpdatedAt != time.Time{}) {
				tmpUat = d.UpdatedAt
				tmp.UpdateAt = &tmpUat
			}
			if err = tmp.Write(v1db); err != nil {
				return err
			}

		}
	}
	slog.Info(fmt.Sprintf("Filelists table copy end %s", time.Since(timeNow)))

	return nil
}
