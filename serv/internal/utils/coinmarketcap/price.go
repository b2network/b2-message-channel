package coinmarketcap

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GetBTCPrice() (float64, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=BTC&convert=USD"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("X-CMC_PRO_API_KEY", "461ba055-ebed-45d7-8366-5efcf57a0acf")
	req.Header.Add("Host", "pro-api.coinmarketcap.com")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var data DataPrice
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}
	if data.Status.ErrorCode != 0 {
		return 0, errors.New("Error getting price")
	}
	return data.Data.BTC.Quote.USD.Price, nil
}

type DataPrice struct {
	Status struct {
		ErrorCode int `json:"error_code"`
	} `json:"status"`
	Data struct {
		BTC struct {
			Quote struct {
				USD struct {
					Price float64 `json:"price"`
				} `json:"USD"`
			} `json:"quote"`
		} `json:"BTC"`
	} `json:"data"`
}
