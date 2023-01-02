package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
)

func ConvertFCYToSGD(original float64, currency string) (float64, error) {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	key := os.Getenv("KEY")
	url := fmt.Sprintf("https://api.freecurrencyapi.com/v1/latest?apikey=%v&currencies=SGD&base_currency=%v", key, currency)
	log.Info(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	data := map[string]interface{}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	rate := data["data"].(map[string]interface{})["SGD"].(float64)
	return original * rate, nil
}
