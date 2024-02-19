package dbfunc

import (
	"gorm.io/gorm"
)

type CountryDB struct {
	gorm.Model
	Id 		int 		`gorm:"primaryKey, autoIncrement"` 
	Name 	string 		`gorm:"type:varchar(255)"`
	Iso   	string  	`gorm:"type:varchar(2)"`
	Flag   	string  	`gorm:"type:varchar(20)"`
	Domain 	string   	`gorm:"type:varchar(4)"`
	Imsi   	[]MccMncDB 	`gorm:"foreignKey:iso"`
}

type MccMncDB struct {
	gorm.Model
	Id 		int `gorm:"primaryKey, autoIncrement"`
	Iso 	string	`gorm:"type:varchar(2)"`
	Mcc 	string 	`gorm:"type:varchar(3)"`
	Mnc 	string 	`gorm:"type:varchar(3)"`
	CC  	string 	`gorm:"type:varchar(3)"`
	Network string 	`gorm:"type:varchar(255)"`
}

type Country struct {
	Name 	string 		`json:"name"`
	Iso   	string  	`json:"iso"`
	Flag   	string  	`json:"flag"`
	Domain 	string   	`json:"domain"`
	Imsi   	[]MccMnc 	`json:"imsi"`
}

type MccMnc struct {
	Iso 	string	`json:"iso"`
	Mcc 	string 	`json:"mcc"`
	Mnc 	string 	`json:"mnc"`
	CC  	string 	`json:"cc"`
	Network string 	`json:"network"`
}

type Flags struct {
	Iso     string `json:"iso"`
	Emoji   string `json:"emoji"`
	Unicode string `json:"unicode"`
	Name    string `json:"name"`
}

type Domains struct {
	Country string `json:"country"`
	Domain  string `json:"domain"`
}

var AllFlags []Flags
var AllDomains []Domains
var AllImsis []MccMnc
var AllCountries []Country

