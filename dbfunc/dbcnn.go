package dbfunc

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB(dbpath string) error {
	var err error

	if os.IsExist(err) {
		os.Remove(dbpath)
	}

	DB, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	SyncDB()

	return nil
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("failed to close database")
	}
	sqlDB.Close()
}

func SyncDB() {
	DB.AutoMigrate(
		&CountryDB{},
		&MccMncDB{},
	)
}

func InsertCountry(c Country) error {
	var imsis []MccMncDB

	for _, imsi := range c.Imsi {
		imsis = append(imsis, MccMncDB{
			Iso:     c.Iso,
			Mcc:     imsi.Mcc,
			Mnc:     imsi.Mnc,
			CC:      imsi.CC,
			Network: imsi.Network,
		})
	}

	country := CountryDB{
		Name:   c.Name,
		Iso:    c.Iso,
		Flag:   c.Flag,
		Domain: c.Domain,
		Imsi:   imsis,
	}

	DB.Create(&country)

	return nil
}