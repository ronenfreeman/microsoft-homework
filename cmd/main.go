package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// type Rdf struct {
// 	XMLName xml.Name `xml:"RDF"` //ignore the smaspace of the tag
// 	// Channel Channel `xml:"channel"` // we arent using chnanel for  now
// 	Items []Item `xml:"item"`
// }

type BitcoinCurrentRate struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
	Bpi struct {
		USD struct {
			Code        string  `json:"code"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
	} `json:"bpi"`
}

func initEnv() int {
	log.Println("Verifying env...")
	if os.Getenv("BITCOIN_ENDPOINT") == "" {
		log.Fatal("BITCOIN_ENDPOINT not set!")
	}
	if os.Getenv("MINUTES_TO_SLEEP") == "" {
		log.Fatal("MINUTES_TO_SLEEP not set!")
	}
	minutesToSleep, err := strconv.Atoi(os.Getenv("MINUTES_TO_SLEEP"))
	if err != nil {
		log.Fatal("Could not convert SECONDS_TO_SLEEP to int")
	}
	return minutesToSleep

}

func getCurrentBitcoinRate() (float64, error) {
	url := os.Getenv("BITCOIN_ENDPOINT")
	var rate float64

	resp, err := http.Get(url)
	if err != nil {
		return rate, err
	}
	defer resp.Body.Close()
	if string(resp.Status) != "200 OK" {
		return rate, errors.New("Error retreiving endpoint: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if string(resp.Status) != "200 OK" {
		return rate, err
	}
	var bitcoinCurrentRate BitcoinCurrentRate
	if err := json.Unmarshal([]byte(body), &bitcoinCurrentRate); err != nil {
		return rate, err
	}
	rate = bitcoinCurrentRate.Bpi.USD.RateFloat
	return rate, nil
}
func main() {
	minutesToSleep := initEnv()
	var ratesArray []float64

	log.Println("Getting Bitcoin prices every minute and the average price every 10 minutes...")
	for {
		rateFloat, err := getCurrentBitcoinRate()
		if err != nil {
			log.Fatal(err)
		}
		//append to array of rates
		ratesArray = append(ratesArray, rateFloat)
		rateString := fmt.Sprintf("%f", rateFloat) // s == "123.456000"

		log.Println("Current Bitcoin price (USD): " + rateString)

		time.Sleep(time.Duration(minutesToSleep) * time.Minute)
	}
}
