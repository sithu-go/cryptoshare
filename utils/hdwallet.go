package utils

import (
	"log"

	"github.com/ygcool/go-hdwallet"
)

type WalletInfo struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
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

	wallet, err := master.GetWallet(hdwallet.CoinType(coinType))
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

	publicKey := wallet.GetKey().PublicHex(false)

	walletInfo := &WalletInfo{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	return walletInfo, nil

}

func NetworkToCoinType(network string) uint32 {
	if network == "ERC20" {
		return hdwallet.ETH
	}
	return hdwallet.TRX
}
