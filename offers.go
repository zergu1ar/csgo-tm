package csgo_tm

import (
	"encoding/json"
)

type (
	TradeRequestGiveP2PResponse struct {
		Success bool   `json:"success"`
		Hash    string `json:"hash"`
		Offer   Offer  `json:"offer"`
		Error   string  `json:"error"`
	}
	TradeRequestGiveP2PAllResponse struct {
		Success bool    `json:"success"`
		Hash    string  `json:"hash"`
		Offers  []Offer `json:"offers"`
		Error   string  `json:"error"`
	}
	TradesResponse struct {
		Success bool    `json:"success"`
		Trades  []Trade `json:"trades"`
	}
	Trade struct {
		Dir       string `json:"dir"`
		TradeId   string `json:"trade_id"`
		BotId     string `json:"bot_id"`
		Timestamp int    `json:"timestamp"`
	}
	Offer struct {
		Partner           int         `json:"partner"`
		Token             string      `json:"token"`
		TradeOfferMessage string      `json:"tradeoffermessage"`
		Items             []TradeItem `json:"items"`
	}
	TradeItem struct {
		AppId     int    `json:"appid"`
		ContextId int    `json:"contextid"`
		AssetsId  string `json:"assetid"`
		Amount    int    `json:"amount"`
	}
)

func (mc *MarketClient) TradeRequestGiveP2P() (error, *TradeRequestGiveP2PResponse) {
	err, body := mc.doRequest("trade-request-give-p2p", "")
	if err != nil {
		return err, nil
	}
	res := TradeRequestGiveP2PResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) TradeRequestGiveP2PAll() (error, *TradeRequestGiveP2PAllResponse) {
	err, body := mc.doRequest("trade-request-give-p2p-all", "")
	if err != nil {
		return err, nil
	}
	res := TradeRequestGiveP2PAllResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}

func (mc *MarketClient) Trades() (error, *TradesResponse) {
	err, body := mc.doRequest("trades", "")
	if err != nil {
		return err, nil
	}
	res := TradesResponse{}
	err = json.Unmarshal(body, &res)
	return err, &res
}
