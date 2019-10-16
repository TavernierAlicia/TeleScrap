package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//to scrape all country links
func scrape() {
	//source link
	resp, err := http.Get("https://countrycode.org")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//get country url
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		//select a href field
		href, hasattr := s.Find("a").First().Attr("href")
		if hasattr {
			//adding the rest of the url to source
			countryURL := "https://countrycode.org" + href

			//access country url
			second, err := http.Get(countryURL)
			if err != nil {
				log.Fatal(err)
			}
			defer second.Body.Close()

			if second.StatusCode != 200 {
				log.Fatalf("Status code error: %d %s", second.StatusCode, resp.Status)
			}
			doc, err := goquery.NewDocumentFromReader(second.Body)
			if err != nil {
				log.Fatal(err)
			}

			//get download file url
			doc.Find("th").Each(func(i int, s *goquery.Selection) {
				//get download link
				href, hasattr := s.Find("a").First().Attr("href")
				if hasattr {
					//add download link
					dlURL := "https://countrycode.org" + href
					matched, err := regexp.Match(`City`, []byte(dlURL))
					if err != nil {
						log.Fatal(err)
					}
					if matched == true {
						fmt.Println(dlURL)

						//download csv
						downloadFromURL(dlURL, countryURL)
					}
					access, err := http.Get(dlURL)
					if err != nil {
						print(err)
					}
					defer access.Body.Close()
				}
			})
		}
	})
}
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

//dl all country csvs for city codes
func downloadFromURL(dlURL string, countryURL string) {
	//get filename
	tokens := strings.Split(countryURL, "/")
	country := tokens[len(tokens)-1]
	fileName := "../countries/" + country + ".csv"
	fmt.Println("Downloading", dlURL, "to", fileName)

	//create file
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(dlURL)
	if err != nil {
		fmt.Println("Error while downloading", dlURL, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", dlURL, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")

	time.Sleep(1 * time.Second)

}

//dl csv file with all country codes
func dlCountries() {
	//url to dl countrycodes.csv
	url := "https://countrycode.org/customer/countryCode/downloadCountryCodes?country="
	//add filename
	fileName := "../countrycode.csv"
	fmt.Println("Downloading", url, "to", fileName)

	//create file
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", fileName, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", fileName, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

//rename unconventionnal countries names
func rename() {
	os.Rename("../countries/uk.csv", "../countries/unitedkingdom.csv")
	os.Rename("../countries/usa.csv", "../countries/unitedstates.csv")
	os.Rename("../countries/uae.csv", "../countries/u.s.virginislands.csv")
	os.Rename("../countries/congo.csv", "../countries/republicofthecongo.csv")
	os.Rename("../countries/congodemocraticrepublic.csv", "../countries/democraticrepublicofthecongo.csv")
	os.Rename("../countries/bosnia.csv", "../countries/bosniaandherzegovina.csv")
	os.Rename("../countries/virginislands.csv", "../countries/britishvirginislands.csv")
	os.Rename("../countries/cocoskeelingislands.csv", "../countries/cocosislands.csv")
	os.Rename("../countries/falklands.csv", "../countries/falklandsislands.csv")
	os.Rename("../countries/pitcairnislands.csv", "../countries/pitcairn.csv")
	os.Rename("../countries/svalbard.csv", "../countries/svalbardandjanmayen.csv")
	os.Rename("../countries/turcsandcaicos.csv", "../countries/turcsandcaicosislands.csv")
	os.Rename("../countries/sthelena.csv", "../countries/sainthelena.csv")
	os.Rename("../countries/stkitts.csv", "../countries/saintkittsandnevis.csv")
	os.Rename("../countries/stlucia.csv", "../countries/saintlucia.csv")
	os.Rename("../countries/stmartin.csv", "../countries/saintmartin.csv")
	os.Rename("../countries/stpierre.csv", "../countries/saintpierreandmiquelon.csv")
	os.Rename("../countries/stvincent.csv", "../countries/saintvincentandthegrenadines.csv")
}

//uncomment to launch
func main() {
	//scrape()
	//dlCountries()
	//rename()
}
