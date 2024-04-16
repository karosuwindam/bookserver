package v0struct

import "gorm.io/gorm"

func ReadBooknames(db *gorm.DB) ([]Booknames, error) {
	tmps := []Booknames{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}

func ReadCopyfiles(db *gorm.DB) ([]Copyfile, error) {
	tmps := []Copyfile{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}

func ReadFilelists(db *gorm.DB) ([]Filelists, error) {
	tmps := []Filelists{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}
