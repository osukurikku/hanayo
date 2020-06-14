package btcconversions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

// CurrencyRates rates with eur/usd to RUB
var CurrencyRates = struct {
	EUR float64
	USD float64
}{
	// 1 rub to eur/usb on 15.06.2020at00:56
	78.37,
	69.70,
}

// GetRates returns the bitcoin rates as JSON.
func GetRates(c *gin.Context) {
	c.JSON(200, CurrencyRates)
}

func init() {
	go func() {
		for {
			updateRates()
			time.Sleep(time.Hour)
		}
	}()
}

func updateRates() {
	var result map[string]interface{}
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
	}

	CurrencyRates.EUR = result["Valute"].(map[string]interface{})["EUR"].(map[string]interface{})["Value"].(float64)
	CurrencyRates.USD = result["Valute"].(map[string]interface{})["USD"].(map[string]interface{})["Value"].(float64)
}
