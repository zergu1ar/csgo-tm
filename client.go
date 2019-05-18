package csgo_tm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	MarketClient struct {
		ApiKey        string
		ctx           context.Context
		Destroy       func()
		currency      string
		requestsQueue chan RequestItem
	}
	RequestItem struct {
		Url          string
		Body         io.Reader
		ResponseChan chan RequestResponse
	}
	RequestResponse struct {
		Error  error
		Body   []byte
		Status int
	}
	PingResponse struct {
		Success bool   `json:"success"`
		Ping    string `json:"ping"`
	}

	WSResponse struct {
		Auth    string `json:"wsAuth"`
		Success bool   `json:"success"`
	}
)

const (
	ApiV2Url    = "https://market.csgo.com/api/v2/%s?key=%s&%s"
	ApiV1Url    = "https://market.csgo.com/api/%s/?key=%s&%s"
	CurrencyRUB = "RUB"
	CurrencyUSD = "USD"
	CurrencyEUR = "EUR"
)

func NewClient(apiKey string, currency string, autoPing bool) *MarketClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := &MarketClient{
		ApiKey:        apiKey,
		ctx:           ctx,
		Destroy:       cancel,
		currency:      currency,
		requestsQueue: make(chan RequestItem, 100),
	}
	if autoPing {
		go client.pingHandler()
	}
	go client.requestsWorker()
	return client
}

func (mc *MarketClient) doRequest(uri, extraParamsString string, apiV2 bool) (error, []byte) {
	var baseUrl string
	if apiV2 {
		baseUrl = ApiV2Url
	} else {
		baseUrl = ApiV1Url
	}

	req := RequestItem{
		Url:          fmt.Sprintf(baseUrl, uri, mc.ApiKey, extraParamsString),
		Body:         nil,
		ResponseChan: make(chan RequestResponse),
	}

	mc.requestsQueue <- req

	res := <-req.ResponseChan
	defer close(req.ResponseChan)

	return res.Error, res.Body
}

func (mc *MarketClient) requestsWorker() {
	for {
		select {
		case req := <-mc.requestsQueue:
			req.ResponseChan <- mc.doSyncRequest(req.Url, req.Body)
			time.Sleep(time.Second)
		case <-mc.ctx.Done():
			return
		}
	}
}

func (mc *MarketClient) doSyncRequest(url string, body io.Reader) RequestResponse {
	req, err := http.NewRequest(http.MethodGet, url, body)
	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return RequestResponse{
			Status: res.StatusCode,
			Error:  err,
		}
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return RequestResponse{
			Status: res.StatusCode,
			Error:  err,
		}
	}
	if res.StatusCode != 200 {
		return RequestResponse{
			Status: res.StatusCode,
			Error:  errors.New(res.Status + string(resBody)),
		}
	}
	return RequestResponse{
		Body:   resBody,
		Status: res.StatusCode,
	}
}

func (mc *MarketClient) pingHandler() {
	tick := time.NewTicker(time.Second * 150)
	for {
		select {
		case <-mc.ctx.Done():
			return
		case <-tick.C:
			_ = mc.Ping()
		}
	}
}

func (mc *MarketClient) Ping() error {
	err, body := mc.doRequest("ping", "", true)
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

func (mc *MarketClient) GetWSToken() (error, string) {
	err, body := mc.doRequest("GetWSAuth", "", false)
	if err != nil {
		return err, ""
	}
	resp := WSResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, ""
	}
	if !resp.Success {
		return errors.New("failed to get ws auth"), ""
	}
	return nil, resp.Auth
}
