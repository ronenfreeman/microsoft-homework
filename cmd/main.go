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

func initEnv() (int, int) {
	log.Println("Verifying env...")
	if os.Getenv("BITCOIN_ENDPOINT") == "" {
		log.Fatal("BITCOIN_ENDPOINT not set!")
	}
	if os.Getenv("MINUTES_TO_SLEEP") == "" {
		log.Fatal("MINUTES_TO_SLEEP not set!")
	}
	if os.Getenv("MINUTES_TO_GET_AVERAGE") == "" {
		log.Fatal("MINUTES_TO_GET_AVERAGE not set!")
	}
	minutesToSleep, err := strconv.Atoi(os.Getenv("MINUTES_TO_SLEEP"))
	if err != nil {
		log.Fatal("Could not convert SECONDS_TO_SLEEP to int")
	}
	minutesToGetAverage, err := strconv.Atoi(os.Getenv("MINUTES_TO_GET_AVERAGE"))
	if err != nil {
		log.Fatal("Could not convert MINUTES_TO_GET_AVERAGE to int")
	}
	//ensure minutesToGetAverage is greater than minutesToSleep
	if minutesToGetAverage <= minutesToSleep {
		log.Fatal("MINUTES_TO_GET_AVERAGE must be greater than MINUTES_TO_SLEEP")
	}

	return minutesToSleep, minutesToGetAverage

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
	minutesToSleep, minutesToGetAverage := initEnv()
	//increment counter to know when to get average
	counter := 1
	var ratesTotal float64
	stringToLog := fmt.Sprintf("Getting Bitcoin rates every %d minute(s) and the average rate every %d minutes...", minutesToSleep, minutesToGetAverage)
	log.Println(stringToLog)
	for {
		rateFloat, err := getCurrentBitcoinRate()
		if err != nil {
			log.Fatal(err)
		}
		//append to ratesTotal
		ratesTotal = ratesTotal + rateFloat
		rateString := fmt.Sprintf("%f", rateFloat)
		log.Println("Current Bitcoin rate (USD): " + rateString)
		// we assume the program wont die ever. so when the array size equals the number of minutes
		// to calculate the average, calculate the average for all values in the array
		if counter == minutesToGetAverage {
			// get average
			averageRate := ratesTotal / float64(minutesToGetAverage)
			stringToLog := fmt.Sprintf("Average Bitcoin rate over %d minute(s) (USD): %f", minutesToGetAverage, averageRate)
			log.Println(stringToLog)
			// reset counter and ratesTotal total
			counter = 0
			ratesTotal = 0
		}
		counter++
		// we assume there are no delays in getting the rate, so the time to sleep is exactly minutesToSleep
		time.Sleep(time.Duration(minutesToSleep) * time.Minute)
		// time.Sleep(time.Duration(minutesToSleep) * time.Second) // use for testing
	}
}
