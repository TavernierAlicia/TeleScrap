package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func addCity() {
	//// OPEN CITY FOLDER ////
	files, err := ioutil.ReadDir("./countries/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if err != nil {
			println(err.Error())
		}

		csvFile, _ := os.Open("./countries/" + file.Name())
		reader := csv.NewReader(bufio.NewReader(csvFile))

		//// GET LINES ////
		for {
			line, err := reader.Read()
			println("here!")
			if err != nil {
				if err == io.EOF {
					println("End of file")
					break
				} else {
					log.Fatal(err.Error())
				}
			}
			//sort and rename unconventionnal fields (if still unconventionnal, it )
			citycodestr := strings.Replace(line[0], "-", "", -1)
			citycodex := strings.Replace(citycodestr, "X", "", -1)
			citycodep1 := strings.Replace(citycodex, ")", "", -1)
			citycodep2 := strings.Replace(citycodep1, "(", "", -1)
			citycode4 := strings.Replace(citycodep2, " + 4 digits", "", -1)
			citycode5 := strings.Replace(citycode4, " + 5 digits", "", -1)
			citycode6 := strings.Replace(citycode5, " + 6 digits", "", -1)
			citycodesix := strings.Replace(citycode6, " + six digits", "", -1)
			citycodesfou := strings.Replace(citycodesix, " plus fou", "", -1)
			citycodeso := strings.Replace(citycodesfou, " o", "", -1)
			citycodeseight := strings.Replace(citycodeso, " + eight digit num", "", -1)
			citycodese := strings.Replace(citycodeseight, "e", "", -1)
			citycodes38 := strings.Replace(citycodese, " ? 38XXXX", "", -1)
			citycode, err := strconv.Atoi(citycodes38)
			if err != nil {
				fmt.Println(err)
			}
			city := line[1]
			countryunf := file.Name()
			country := strings.Replace(countryunf, ".csv", "", -1)
			if city != "Description" && city != "" {
				aCity := true
				Database(aCity, citycode, citycodes38, city, country) //, countrycode)
			}
		}
	}
	return
}

func addCountry() {

	//// GET COUNTRY DATA ////
	csvFile, _ := os.Open("countrycode.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	//// GET LINES ////
	for {
		line, err := reader.Read()
		println("here!")
		if err != nil {
			if err == io.EOF {
				println("End of file")
				break
			} else {
				log.Fatal(err.Error())
			}
		}

		//get strings and remove indesirable chars
		naCodestr := strings.Replace(line[8], "-", "", -1)
		naCode, err := strconv.Atoi(naCodestr)
		if err != nil {
			fmt.Println(err)
		}
		countrySpell := line[0]
		countryDash := strings.Replace(countrySpell, "-", "", -1)
		countrySpace := strings.Replace(countryDash, " ", "", -1)
		country := strings.ToLower(countrySpace)

		aCity := false
		Database(aCity, naCode, naCodestr, countrySpell, country)
	}
	return
}

//Database to request database
func Database(aCity bool, code int, codestr string, neededplace string, country string) {

	//connect database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/countrycodes")
	//checking errors
	if err != nil {
		println("Erreur à la connection")
	} else {
		println("It's a success!")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		println("Erreur après le ping")
	} else {
		println("DOUBLE SUCCESS OF THE DEATH")
	}

	//City Query
	if code != 0 {
		if aCity == true {
			cit, err := db.Exec(`INSERT INTO city (city_name, city_code, country_name) VALUES (?, ?, ?)`, neededplace, code, country)
			if err != nil {
				println("Erreur requête")
				return
			}
			n, _ := cit.RowsAffected()
			if n == 0 {
				println("not inserted")
			}

		} else {

			//Country Query
			coun, err := db.Exec(`INSERT INTO country (name, country_spelling, national_code) VALUES (?, ?, ?)`, country, neededplace, code)
			if err != nil {
				println("Erreur requête")
				return
			}
			n, _ := coun.RowsAffected()
			if n == 0 {
				println("not inserted")
			}
		}
	} else {
		if aCity == true {
			//junk Query
			junk, err := db.Exec(`INSERT INTO junk (city_name, country_name, city_code) VALUES (?, ?, ?)`, country, neededplace, codestr)
			if err != nil {
				println("Erreur requête")
				return
			}
			n, _ := junk.RowsAffected()
			if n == 0 {
				println("not inserted")
			}
		}
	}
}

func main() {
	//scrape()

	addCountry()
	addCity()
	//we need to get this one time too
	//https://countrycode.org/customer/countryCode/downloadCountryCodes?country=
	//for great format country and country code
}
