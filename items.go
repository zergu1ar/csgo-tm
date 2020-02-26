package csgo_tm

import (
	"encoding/json"
	"fmt"
)

type (
	ItemHistoryResponse struct {
		Success bool           `json:"success"`
		Max     int            `json:"max"`
		Min     int            `json:"min"`
		Average int            `json:"average"`
		Number  int            `json:"number"`
		History []*ItemHistory `json:"history"`
	}
	ItemHistory struct {
		Price string `json:"l_price"`
		Time  string `json:"l_time"`
	}
	ItemInfoResponse struct {
		ClassId    string       `json:"classid"`
		InstanceId string       `json:"instanceid"`
		HashName   string       `json:"market_hash_name"`
		MinPrice   string       `json:"min_price"`
		Offers     []*ItemOffer `json:"offers"`
	}
	ItemOffer struct {
		Price   string `json:"price"`
		Count   string `json:"count"`
		MyCount string `json:"my_count"`
	}
	BuyOffersResponse struct {
		Success   bool        `json:"success"`
		BestOffer string      `json:"best_offer"`
		Offers    []*BuyOffer `json:"offers"`
	}
	BuyOffer struct {
		Price   string `json:"o_price"`
		Count   string `json:"c"`
		MyCount string `json:"my_count"`
	}
)

func (mc *MarketClient) ItemHistory(classId, instanceId string) (error, *ItemHistoryResponse) {
	err, body := mc.doRequest(
		fmt.Sprintf("ItemHistory/%s_%s", classId, instanceId),
		"",
		false,
	)
	if err != nil {
		return err, nil
	}
	res := ItemHistoryResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) ItemInfo(classId, instanceId, lang string) (error, *ItemInfoResponse) {
	err, body := mc.doRequest(
		fmt.Sprintf("ItemInfo/%s_%s/%s", classId, instanceId, lang),
		"",
		false,
	)
	if err != nil {
		return err, nil
	}
	res := ItemInfoResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) BuyOffers(classId, instanceId string) (error, *BuyOffersResponse) {
	err, body := mc.doRequest(
		fmt.Sprintf("BuyOffers/%s_%s", classId, instanceId),
		"",
		false,
	)
	if err != nil {
		return err, nil
	}
	res := BuyOffersResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}