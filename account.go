package csgo_tm

import (
	"encoding/json"
)

type (
	GetMoneyResponse struct {
		Money    float32 `json:"money"`
		Currency string  `json:"currency"`
		Success  bool    `json:"success"`
	}
	GoOfflineResponse struct {
		Success bool `json:"success"`
	}
	UpdateInventoryResponse struct {
		Success bool `json:"success"`
	}
	MyInventoryResponse struct {
		Success bool            `json:"success"`
		Items   []InventoryItem `json:"items"`
	}
	InventoryItem struct {
		Id             string  `json:"id"`
		ClassId        string  `json:"classid"`
		InstanceId     string  `json:"instanceid"`
		MarketHashName string  `json:"market_hash_name"`
		MarketPrice    float64 `json:"market_price"`
		Tradable       int     `json:"tradable"`
	}
)

func (mc *MarketClient) GetMoney() (error, *GetMoneyResponse) {
	err, body := mc.doRequest("get-money", "", true)
	if err != nil {
		return err, nil
	}
	res := GetMoneyResponse{}
	err = json.Unmarshal(body, &res)
	return nil, &res
}

func (mc *MarketClient) GoOffline() (error, *GoOfflineResponse) {
	err, body := mc.doRequest("go-offline", "", true)
	if err != nil {
		return err, nil
	}
	res := GoOfflineResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) UpdateInventory() (error, *UpdateInventoryResponse) {
	err, body := mc.doRequest("update-inventory", "", true)
	if err != nil {
		return err, nil
	}
	res := UpdateInventoryResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) MyInventory() (error, *MyInventoryResponse) {
	err, body := mc.doRequest("my-inventory/", "", true)
	if err != nil {
		return err, nil
	}
	res := MyInventoryResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}
