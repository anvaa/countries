package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"dbfunc"
)

func main() {
	var err error
	// Read in CSV file
	dbfunc.AllFlags, err = ReadInFlags("./flags.csv")
	if err != nil {
		fmt.Println(err)
	}

	dbfunc.AllDomains, err = ReadInDomains("./domains.csv")
	if err != nil {
		fmt.Println(err)
	}

	dbfunc.AllImsis, err = ReadInImsis("./mccmnc.csv")
	if err != nil {
		fmt.Println(err)
	}

	data, err := ReadInCounties("./coutries.csv")
	if err != nil {
		fmt.Println(err)
	}

	c, err := AddFlagToCountry(data)
	if err != nil {
		fmt.Println(err)
	}

	c, err = AddDomainToCountry(c)
	if err != nil {
		fmt.Println(err)
	}

	c, err = AddImsiToCountry(c)
	if err != nil {
		fmt.Println(err)
	}

	// Print the data
	// for _, country := range c {
	// 	fmt.Printf("%s, %s, %s, %s, %s\n",
	// 		country.Name, country.Iso, country.Flag, country.Domain, country.Imsi)
	// }

	err = SaveCountyToJson(c)
	if err != nil {
		fmt.Println(err)
	}

	err = SaveToDb(c)
	if err != nil {
		fmt.Println(err)
	}	

}

func SaveToDb(c []dbfunc.Country) error {

	if err := dbfunc.ConnectToDB("./data/countries.db"); err != nil {
		return err
	}

	for _, country := range c {
		if err := dbfunc.InsertCountry(country); err != nil {
			return err
		}
	}

	return nil
}

func SaveCountyToJson(c []dbfunc.Country) error {
	file, err := os.Create("./data/country.json")
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	file.Write(jsonData)
	return nil
}

func ReadInFlags(file string) ([]dbfunc.Flags, error) {
	var flags []dbfunc.Flags
	csvFile, err := os.Open(file)
	if err != nil {
		return flags, err
	}
	defer csvFile.Close()

	byteValue, _ := io.ReadAll(csvFile)
	lines := strings.Split(string(byteValue), "\n")
	for _, line := range lines {
		flag := strings.Split(line, ",")
		flags = append(flags, dbfunc.Flags{
			Iso:     flag[0],
			Emoji:   flag[1],
			Unicode: flag[2],
			Name:    flag[3],
		})
	}

	return flags, nil
}

func ReadInCounties(file string) ([]dbfunc.Country, error) {
	var country []dbfunc.Country
	csvFile, err := os.Open(file)
	if err != nil {
		return country, err
	}
	defer csvFile.Close()

	byteValue, _ := io.ReadAll(csvFile)
	lines := strings.Split(string(byteValue), "\n")
	for _, line := range lines {
		c := strings.Split(line, ",")
		
		co := strings.Replace(c[0], "\"", "", -1)
		iso := strings.Replace(c[1], "\r", "", -1)
		iso = strings.Replace(iso, "\"", "", -1)

		country = append(country, dbfunc.Country{
			Name: co,
			Iso:  iso,
		})
	}

	return country, nil
}

func ReadInDomains(file string) ([]dbfunc.Domains, error) {
	var domains []dbfunc.Domains
	csvFile, err := os.Open(file)
	if err != nil {
		return domains, err
	}
	defer csvFile.Close()

	byteValue, _ := io.ReadAll(csvFile)
	lines := strings.Split(string(byteValue), "\n")
	for _, line := range lines {
		domain := strings.Split(line, ",")

		dom := strings.Replace(domain[1], " ", "", -1)
		// dom = strings.Replace(dom, " ", "", -1)

		domains = append(domains, dbfunc.Domains{
			Country: strings.Replace(domain[0], "\"", "", -1),
			Domain:  dom,
		})
	}

	return domains, nil
}

func ReadInImsis(file string) ([]dbfunc.MccMnc, error) {
	var imsis []dbfunc.MccMnc
	csvFile, err := os.Open(file)
	if err != nil {
		return imsis, err
	}
	defer csvFile.Close()

	byteValue, _ := io.ReadAll(csvFile)
	lines := strings.Split(string(byteValue), "\n")
	for _, line := range lines {
		imsi := strings.Split(line, ",")
		imsis = append(imsis, dbfunc.MccMnc{
			Mcc: imsi[0],
			Mnc: imsi[1],
			Iso: imsi[2],
			CC:  imsi[4],
			Network: strings.Replace(imsi[5], "\r", "", -1),
		})
		// fmt.Println(imsi[0],imsi[1],imsi[2],imsi[3],imsi[4],imsi[5])
	}

	return imsis, nil
}

func AddFlagToCountry(country []dbfunc.Country) ([]dbfunc.Country, error) {
	for i, c := range country {
		if c.Flag == "" {
			country[i].Flag = findFlagToCountry(country[i])
		}
	}

	return country, nil
}

func findFlagToCountry(c dbfunc.Country) string {

	for _, flag := range dbfunc.AllFlags {
		if c.Iso == flag.Iso {
			// check if flag is valid
			if utf8.ValidString(flag.Unicode) {
				// Remove the U+ from the unicode string
				unicodeHex := strings.Replace(flag.Unicode, "U+", "", -1)
				unicode := strings.Split(unicodeHex, " ")

				// Convert the hex to a code point
				cp1, _ := strconv.ParseInt(unicode[0], 16, 32)
				cp2, _ := strconv.ParseInt(unicode[1], 16, 32)

				// Convert the code point to a rune, then to a string
				emoji := string(rune(cp1))
				emoji = emoji + string(rune(cp2))

				return emoji
			}
		}
	}

	return "-"
}

func AddDomainToCountry(country []dbfunc.Country) ([]dbfunc.Country, error) {
	for i, c := range country {
		if c.Domain == "" {
			country[i].Domain = findDomainToCountry(country[i])
		}
	}

	return country, nil
}

func findDomainToCountry(c dbfunc.Country) string {

	for _, d := range dbfunc.AllDomains {
		if c.Name == d.Country {
			return d.Domain
		}
	}

	return strings.ToLower(c.Iso)
}

func AddImsiToCountry(country []dbfunc.Country) ([]dbfunc.Country, error) {
	for i := range country {
		country[i].Imsi = findImsiToCountry(country[i])
	}
	return country, nil
}

func findImsiToCountry(c dbfunc.Country) []dbfunc.MccMnc {
	var imsis []dbfunc.MccMnc
	for _, i := range dbfunc.AllImsis {
		// fmt.Println(c.ISO, i.Iso)
		if strings.ToLower(c.Iso) == i.Iso {
			imsis = append(imsis, i)
		}
	}
	return imsis
}
