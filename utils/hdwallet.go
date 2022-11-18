package utils

import (
	"log"

	"github.com/ygcool/go-hdwallet"
)

type WalletInfo struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

func GetInfoFromMnemonic(mnemonic string, network string) (*WalletInfo, error) {
	coinType := NetworkToCoinType(network)
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	wallet, err := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(coinType), hdwallet.AddressIndex(0))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	address, err := wallet.GetAddress()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	privateKey, err := wallet.GetPrivateKey()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	walletInfo := &WalletInfo{
		Address:    address,
		PrivateKey: privateKey,
	}
	return walletInfo, nil

}

func NetworkToCoinType(network string) uint32 {
	if network == "ERC20" {
		return hdwallet.ETH
	}
	return hdwallet.TRX
}
