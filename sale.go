package csgo_tm

import (
	"encoding/json"
	"fmt"
	"math"
)

type (
	AddToSaleResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		ItemId  string `json:"item_id"`
	}
	SetPriceResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	ItemsResponse struct {
		Success bool `json:"success"`
		Items   []*Item
	}
	Item struct {
		Id         string `json:"item_id"`
		AssetId    string `json:"assetid"`
		ClassId    string `json:"classid"`
		InstanceId string `json:"instanceid"`
		HashName   string `json:"market_hash_name"`
		Position   int    `json:"position"`
		Price      int    `json:"price"`
		Currency   string `json:"currency"`
		Status     int    `json:"status"`
		LiveTime   int    `json:"live_time"`
	}
)

const (
	ItemStatusInTrade       = 1
	ItemStatusNeedTrade     = 2
	ItemStatusWaitingSeller = 3
	ItemStatusNeedReceive   = 4
)

func (mc *MarketClient) AddToSale(item InventoryItem, price float64) (error, *AddToSaleResponse) {
	extraParams := "id=%s&price=%d&cur=%s"
	err, body := mc.doRequest(
		"add-to-sale",
		fmt.Sprintf(extraParams, item.Id, int(math.Round(price))*100, mc.currency),
		true,
	)
	if err != nil {
		return err, nil
	}
	res := AddToSaleResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) SetPrice(itemId string, price float64) (error, *SetPriceResponse) {
	err, body := mc.doRequest(
		"set-price",
		fmt.Sprintf("item_id=%s&price=%d&cur=%s", itemId, int(math.Round(price))*100, mc.currency),
		true,
	)
	if err != nil {
		return err, nil
	}
	res := SetPriceResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) Items() (error, *ItemsResponse) {
	err, body := mc.doRequest(
		"items",
		"",
		true,
	)
	if err != nil {
		return err, nil
	}
	res := ItemsResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}
