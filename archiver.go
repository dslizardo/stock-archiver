package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const basedURL = "http://www.pse.com.ph/stockMarket/home.html"
const credentials = "phstock_service.json"

// TODO make this dynamic
const driveFolder = "1kfUMYjjhon3k9oDSCdgBXjJV_2P-xSnj"

var csvHeader = []string{"alias", "symbol", "lastTradedPrice"}

// Stock is a struct representation of individual stock information
type Stock struct {
	SecurityAlias   string `json:"securityAlias"`
	SecuritySymbol  string `json:"securitySymbol"`
	LastTradedPrice string `json:"lastTradedPrice"`
}

func getStocks() []Stock {
	log.Print("Getting stocks prices")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", basedURL, nil)
	req.Header.Set("Referer", basedURL)

	params := req.URL.Query()
	params.Set("method", "getSecuritiesAndIndicesForPublic")
	params.Set("ajax", "true")
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	stocks := []Stock{}
	if err := json.Unmarshal(body, &stocks); err != nil {
		log.Fatal(err)
	}

	return stocks
}

func uploadStocks(filename string) {
	b, _ := ioutil.ReadFile(credentials)
	conf, _ := google.JWTConfigFromJSON(b, drive.DriveFileScope)

	client := conf.Client(oauth2.NoContext)
	drivesvc, err := drive.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("error opening %q: %v", filename, err)
	}

	parents := []string{driveFolder}
	log.Print("Uploading ", filename)
	driveFile, err := drivesvc.Files.Create(&drive.File{Name: filename, Parents: parents}).Media(file).Do()

	if err != nil {
		log.Printf("Got drive.File, err: %#v, %v", driveFile, err)

	} else {
		log.Print("Uploaded successfully!")
	}
}

func convertToCSV(stocks []Stock) string {
	today := time.Now()
	filename := today.Format("2019-01-01")
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Problem creating file", err)
	}
	log.Print(filename, "file created")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(csvHeader)

	for _, value := range stocks {
		err := writer.Write([]string{value.SecurityAlias, value.SecuritySymbol, value.LastTradedPrice})
		if err != nil {
			log.Fatal("Problem writing to file", err)
		}
	}

	return filename
}

func archiveStocks() {
	today := time.Now()
	// PSE closing time is 3:30 PM
	closingTime := time.Date(today.Year(), today.Month(), today.Day(), 15, 35, 0, 0, today.Location())
	diff := closingTime.Sub(today)

	if diff < 0 {
		closingTime = closingTime.Add(24 * time.Hour)
		diff = closingTime.Sub(today)
	}

	for {
		time.Sleep(diff)
		diff = 24 * time.Hour
		stocks := getStocks()
		file := convertToCSV(stocks)
		uploadStocks(file)
	}
}

func main() {
	archiveStocks()
}
