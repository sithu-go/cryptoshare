package service

import (
	"fmt"
	"net/http"
)

type trc20Service struct {
}

func newTRC20Service() *trc20Service {
	return &trc20Service{}
}

func (s *trc20Service) GetAccountInfo(address string) {
	reqUrl := fmt.Sprintf("%s/v3/tron/account/", address)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}
	req.Header.Add("x-api-key", apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
}
