package v1struct

import "gorm.io/gorm"

func ReadCopyfiles(db *gorm.DB) ([]Copyfile_sql, error) {
	tmps := []Copyfile_sql{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}

func (t *Copyfile_sql) Write(db *gorm.DB) error {
	if results := db.Create(t); results.Error != nil {
		return results.Error
	}
	return nil
}

func ReadBooknames(db *gorm.DB) ([]Booknames_sql, error) {
	tmps := []Booknames_sql{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}

func (t *Booknames_sql) Write(db *gorm.DB) error {
	if results := db.Create(t); results.Error != nil {
		return results.Error
	}
	return nil
}
func ReadFilelists(db *gorm.DB) ([]Filelists_sql, error) {
	tmps := []Filelists_sql{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}

func (t *Filelists_sql) Write(db *gorm.DB) error {
	if results := db.Create(t); results.Error != nil {
		return results.Error
	}
	return nil
}
func ReadUploadTmp(db *gorm.DB) ([]UploadTmp_sql, error) {
	tmps := []UploadTmp_sql{}

	if results := db.Find(&tmps); results.Error != nil {
		return tmps, results.Error
	}
	return tmps, nil
}
