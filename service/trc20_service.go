package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type trc20Service struct {
}

func newTRC20Service() *trc20Service {
	return &trc20Service{}
}

type TRC20AccountInfo struct {
	Balance uint64              `json:"balance"`
	TRC20   []map[string]string `json:"trc20"`
}

func (s *trc20Service) GetAccountInfo(address string) {

	reqUrl := fmt.Sprintf("%s/v3/tron/account/%s", tatumBaseURL, address)
	req, err := http.NewRequest("GET", reqUrl, nil)

	if err != nil {
		fmt.Println()
		return
	}
	req.Header.Add("x-api-key", tatumApiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	accountInfo := &TRC20AccountInfo{}
	err = json.NewDecoder(res.Body).Decode(accountInfo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%+v\n", accountInfo)
}
