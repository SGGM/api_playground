package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Response struct {
	Meta struct {
		Code       int    `json:"code"`
		Disclaimer string `json:"disclaimer"`
	} `json:"meta"`
	ResponseData struct {
		Date  string             `json:"date"`
		Base  string             `json:"base"`
		Rates map[string]float64 `json:"rates"`
	} `json:"response"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while geting api key from env")
	}

	startApp()

}

func startApp() {
	arg := parseArgs(os.Args)
	if arg == "-h" {
		showHelp()
	} else {
		makeRequest(arg)
	}
}

func showHelp() {
	log.Println("\nHELP FOR CURRENCY EXCHANGE PROGRAMM\nSHOWS RATE FOR USD - YOUR CURRENCY\nusage: main.go [Currecy CODE]")
}

func makeRequest(arg string) {
	url := configurateURL()

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Fatalf("error creating HTTP request^ %v", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request^ %v", err)
	}

	defer resp.Body.Close()

	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatalf("error decoding JSON: %v", err)
	}

	arg = strings.ToUpper(arg)

	argRate, found := data.ResponseData.Rates[arg]
	if !found {
		log.Fatalf("Currency code '%s' not found", arg)
	}
	log.Printf("\nThe rate of %v is %.2f", arg, argRate)
}

func parseArgs(args []string) string {

	if len(args) != 2 {
		log.Fatalf("\nThere should be one argument [RUB] or other currency \nFor more info use main.go -h")
	}
	arg := args[1]

	return arg
}

func configurateURL() string {
	key := os.Getenv("API_KEY")
	url := "https://api.currencybeacon.com/v1/latest?api_key="
	url += key
	return url
}
