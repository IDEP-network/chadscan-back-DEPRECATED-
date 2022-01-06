package cmc

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/everstake/cosmoscan-api/log" //"log"
)

const apiURL = "https://pro-api.coinmarketcap.com"

type (
	CMC struct {
		cfg    config.Config
		client *http.Client
	}
	CurrenciesResponse struct {
		Status struct {
			ErrorCode    int    `json:"error_code"`
			ErrorMessage string `json:"error_message,omitempty"`
		} `json:"status"`
		Data []Currency `json:"data"`
	}
	Currency struct {
		CirculatingSupply decimal.Decimal `json:"circulating_supply"`
		CMCRank           int             `json:"cmc_rank"`
		TotalSupply       decimal.Decimal `json:"total_supply"`
		Symbol            string          `json:"symbol"`
		Quote             map[string]struct {
			MarketCap        decimal.Decimal `json:"market_cap"`
			PercentChange1h  decimal.Decimal `json:"percent_change_1h"`
			PercentChange7d  decimal.Decimal `json:"percent_change_7d"`
			PercentChange24h decimal.Decimal `json:"percent_change_24h"`
			Price            decimal.Decimal `json:"price"`
			Volume24h        decimal.Decimal `json:"volume_24h"`
		} `json:"quote"`
	}
)

func NewCMC(cfg config.Config) *CMC {
	return &CMC{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (cmc *CMC) request(endpoint string, data interface{}) error {
	url := fmt.Sprintf("%s%s", apiURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %s", err.Error())
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", cmc.cfg.CMCKey)
	resp, err := cmc.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %s", err.Error())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(d, data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func (cmc *CMC) GetMarketData() (price, volume24h decimal.Decimal, err error) {
	var currencyResp CurrenciesResponse
	err = cmc.request("/v1/cryptocurrency/listings/latest", &currencyResp)
	log.Info("GetMarketData() entered")
	if currencyResp.Status.ErrorCode != 0 {
		return price, volume24h, fmt.Errorf("error code: %d, msg: %s", currencyResp.Status.ErrorCode, currencyResp.Status.ErrorMessage)
	}

	for _, currency := range currencyResp.Data {
		if strings.ToLower(currency.Symbol) == config.Currency {
			quote, ok := currency.Quote["USD"]
			if !ok {
				return price, volume24h, fmt.Errorf("not found USD quote")
			}
			price = quote.Price.Truncate(8)
			volume24h = quote.Volume24h.Truncate(2)
			break
		}
	}

	return price, volume24h, err
}
