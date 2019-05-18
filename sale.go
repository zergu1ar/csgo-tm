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
)

func (mc *MarketClient) AddToSale(item InventoryItem, price float64) (error, *AddToSaleResponse) {
	extraParams := "id=%s&price=%d&cur=%s"
	err, body := mc.doRequest("add-to-sale", fmt.Sprintf(extraParams, item.Id, int(math.Round(price))*100, mc.currency), true)
	if err != nil {
		return err, nil
	}
	res := AddToSaleResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}
