package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Product    string       `json:"product"`
	Init       string       `json:"init"`
	DataSeries []DataSeries `json:"dataseries"`
}

type DataSeries struct {
	Timepoint   int64   `json:"timepoint"`
	Cloudcover  int64   `json:"cloudcover"`
	LiftedIndex int64   `json:"lifted_index"`
	PrecType    string  `json:"prec_type"`
	PrecAmount  int64   `json:"prec_amount"`
	Temp2m      int64   `json:"temp2m"`
	Rh2m        string  `json:"rh2m"`
	Wind10m     Wind10m `json:"wind10m"`
	Weather     string  `json:"weather"`
}

type Wind10m struct {
	Direction string `json:"direction"`
	Speed     int64  `json:"speed"`
}

func ExtractWeather() *[]Response {

	var responses []Response

	responseObject, _ := extractPage()
	responses = append(responses, *responseObject)

	return &responses
}

func extractPage() (*Response, int) {

	searchUrl := "https://www.7timer.info/bin/api.pl?lon=174.766&lat=-36.843&product=civil&output=json"

	req, err := http.NewRequest(http.MethodGet, searchUrl, nil)
	if err != nil {
		fmt.Print("Request failed: " + err.Error())
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		fmt.Print("Response failed: " + err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//fmt.Println(string(responseData))

	var responseObject Response
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	returnedRows := len(responseObject.DataSeries)
	fmt.Println("Returned rows: " + fmt.Sprint(returnedRows))

	return &responseObject, returnedRows
}
