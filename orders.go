package csgo_tm

import (
	"encoding/json"
	"fmt"
)

type (
	DeleteOrderResponse struct {
		Success       bool `json:"success"`
		DeletedOrders int  `json:"deleted_orders"`
	}
	ProcessOrderResponse struct {
		Success bool   `json:"success"`
		Way     string `json:"way"`
		Price   int    `json:"price"`
	}
	GetOrdersResponse struct {
		Success bool     `json:"success"`
		Orders  []*Order `json:"Orders"`
	}
	Order struct {
		ClassId    string `json:"i_classid"`
		InstanceId string `json:"i_instanceid"`
		HashName   string `json:"i_market_hash_name"`
		Name       string `json:"i_market_name"`
		Price      string `json:"o_price"`
		State      string `json:"o_state"`
	}
	OrderFor struct {
		Partner string
		Token   string
	}
)

func (mc *MarketClient) ProcessOrder(item *Order, orderFor *OrderFor) (error, *ProcessOrderResponse) {
	url := fmt.Sprintf("ProcessOrder/%v/%v/%v", item.ClassId, item.InstanceId, item.Price)
	extraParams := ""
	if orderFor != nil {
		extraParams = fmt.Sprintf("partner=%s&token=%s", orderFor.Partner, orderFor.Token)
	}
	err, body := mc.doRequest(url, extraParams, false)
	if err != nil {
		return err, nil
	}
	res := ProcessOrderResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) DeleteOrders() (error, *DeleteOrderResponse) {
	err, body := mc.doRequest("DeleteOrders", "", false)
	if err != nil {
		return err, nil
	}
	res := DeleteOrderResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) GetOrders() (error, *GetOrdersResponse) {
	err, body := mc.doRequest("GetOrders", "", false)
	if err != nil {
		return err, nil
	}
	res := GetOrdersResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}
