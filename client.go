package csgo_tm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	MarketClient struct {
		ApiKey   string
		ctx      context.Context
		Destroy  func()
		currency string
	}

	PingResponse struct {
		Success bool   `json:"success"`
		Ping    string `json:"ping"`
	}
)

const (
	ApiURL      = "https://market.csgo.com/api/v2/%s?key=%s&%s"
	CurrencyRUB = "RUB"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
)

func NewClient(apiKey string, currency string) *MarketClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := &MarketClient{
		ApiKey:   apiKey,
		ctx:      ctx,
		Destroy:  cancel,
		currency: currency,
	}
	go client.pingHandler()
	return client
}

func (mc *MarketClient) doRequest(uri, extraParamsString string) (error, []byte) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ApiURL, uri, mc.ApiKey, extraParamsString), nil)
	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status + string(body)), nil
	}
	return nil, body
}

func (mc *MarketClient) pingHandler() {
	tick := time.NewTicker(time.Second * 150)
	for {
		select {
		case <-mc.ctx.Done():
			return
		case <-tick.C:
			mc.doPing()
		}
	}
}

func (mc *MarketClient) doPing() error {
	err, body := mc.doRequest("ping", "")
	if err != nil {
		return err
	}
	resp := PingResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New("invalid ping answer")
	}
	return nil
}
