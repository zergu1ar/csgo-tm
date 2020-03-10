package csgo_tm

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	DbResponse struct {
		Time string `json:"time"`
		Db   string `json:"db"`
	}
	CSGOItem struct {
		ClassId    string
		InstanceId string
		Rarity     string
		Quality    string
		Name       string
		HashName   string
		Color      string
		Image      string
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

func (mc *MarketClient) GetCSGOItems() (error, []*CSGOItem) {
	res, err := http.Get(
		"https://market.csgo.com/itemdb/current_730.json",
	)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	itemsUrl := DbResponse{}
	err = json.Unmarshal(resBody, &itemsUrl)
	if err != nil {
		return err, nil
	}

	// get request
	csvDataUrl := fmt.Sprintf("https://market.csgo.com/itemdb/%s", itemsUrl.Db)
	resp, err := http.Get(csvDataUrl)
	if err != nil {
		return err, nil
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return err, nil
	}

	var items []*CSGOItem
	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}
		items = append(items, &CSGOItem{
			ClassId:    row[0],
			InstanceId: row[1],
			Rarity:     row[5],
			Quality:    row[6],
			Name:       row[10],
			HashName:   row[12],
			Color:      row[13],
			Image:      fmt.Sprintf("//cdn.csgo.com/item/%s/%s", url.QueryEscape(row[10]), "/300.png"),
		})
	}
	return nil, items
}
