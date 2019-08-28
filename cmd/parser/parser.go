package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	UsdCode     = "USDRUB"
	EurCode     = "EURRUB"
	ExternalApi = "https://quoteorg.fxclub.org/info/ru?symbols=" + UsdCode + "," + EurCode
	rate        rawData
	Debug       = false
)

type ExternalApiResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   struct {
		Instruments []struct {
			Code      string  `json:"code"`
			Name      string  `json:"name"`
			Tradable  bool    `json:"tradable"`
			Rate      float64 `json:"rate"`
			Link      string  `json:"link"`
			Tradelink string  `json:"tradelink"`
		} `json:"instruments"`
		Count int `json:"count"`
	} `json:"data"`
}

type rawData struct {
	Usd float64
	Eur float64
}

func toLog(v ...interface{}) {
	if Debug {
		log.Println(v...)
	}
}

func parseExternalApi(rate *rawData) {

	resp, err := http.Get(ExternalApi)

	if err != nil {

		toLog("Get External Api Error => ", err.Error())

	} else {
		defer func() {

			err := resp.Body.Close()

			if err != nil {

				toLog("Response Body Close Error => " + err.Error())

			}

		}()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {

			toLog("Response Body Read Error => " + err.Error())

		} else {
			var respData ExternalApiResponse

			err = json.NewDecoder(strings.NewReader(string(body))).Decode(&respData)

			if err != nil {

				toLog("Response Decode Json Error => " + err.Error())

			} else {
				if respData.Code == 200 {

					rate.Usd = 0
					rate.Eur = 0

					for _, Instrument := range respData.Data.Instruments {

						if UsdCode == Instrument.Code {

							rate.Usd = Instrument.Rate

						} else if EurCode == Instrument.Code {

							rate.Eur = Instrument.Rate

						}

					}

					db, err := sql.Open("sqlite3", "go_test.sqlite")

					if err != nil {

						toLog("Connect DB Error => " + err.Error())

					} else {

						Stmt, err := db.Prepare("INSERT INTO raw_data(usd, eur) values(?,?)")

						if err != nil {

							toLog("Prepare Query DB Error=> " + err.Error())

						} else {

							_, err = Stmt.Exec(rate.Usd, rate.Eur)

							if err != nil {

								toLog("Insert DB Error=> " + err.Error())

							} else {

								toLog("DB Inserted => " + fmt.Sprintf("%f | %f", rate.Usd, rate.Eur))

							}

							err = Stmt.Close()

							if err != nil {

								toLog("Statements Close DB Error=> " + err.Error())

							}

						}

						err = db.Close()

						if err != nil {

							toLog("Close DB Error=> " + err.Error())

						}

					}

				}
			}

		}
	}

}

func main() {

	for {

		parseExternalApi(&rate)

		toLog(" THE END LOOP ")

		time.Sleep(20 * time.Second)

	}

}
