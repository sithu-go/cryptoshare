package service

import (
	"context"
	"crypto/ecdsa"
	"cryptoshare/conf"
	"cryptoshare/dto"
	"cryptoshare/model"
	"cryptoshare/utils/token"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

type erc20Service struct {
	EtherClient *ethclient.Client
}

func newERC20Service() *erc20Service {
	client, err := ethclient.Dial(fmt.Sprintf("%v/v3/%v", conf.INFURA_BASE_URL, conf.INFURA_API_KEY))
	if err != nil {
		log.Fatal(err)
	}

	return &erc20Service{
		EtherClient: client,
	}
}

func (s *erc20Service) Close() {
	s.EtherClient.Close()
	log.Println("Ether client disconnected successfully.")
}

func (s *erc20Service) GetBalance(address string) (*dto.ERC20BalanceResp, error) {
	account := common.HexToAddress(address)
	balance, err := s.EtherClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Println("err on ", err)
		return nil, err
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethVal := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	//get usdt Balance
	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	instance, err := token.NewToken(contractAddress, s.EtherClient)
	if err != nil {
		log.Println("err on ", err)
		return nil, err
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, account)
	if err != nil {
		log.Println("Error retrieving balance")
		return nil, err
	}

	fbalance = new(big.Float)
	fbalance.SetString(bal.String())
	usdtBal := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(6)))

	ETH, _ := ethVal.Float64()
	USDT, _ := usdtBal.Float64()

	return &dto.ERC20BalanceResp{
		USDTBalance: USDT,
		ETHBalance:  ETH,
	}, nil
}

func (s *erc20Service) TransferERC20USDT(transferReq *dto.TransferReq) (string, error) {
	privateKey, err := crypto.HexToECDSA(transferReq.PrivateKey)
	if err != nil {
		log.Println(err, "Error parsing HexToECDSA")
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("error casting public key to ECDSA")
		return "", nil
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := s.EtherClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err, "Fail Checking Transaction Pending state")
		return "", err
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := s.EtherClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err, "Error gettig suggestion gas price")
		return "", err
	}

	toAddress := common.HexToAddress(transferReq.ToAddress)
	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.New256()
	// hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	amount := new(big.Int)
	amount.SetString(fmt.Sprintf("%v000000", transferReq.Amount), 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := s.EtherClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Data:     data,
		GasPrice: gasPrice,
		Value:    value,
	})
	// gasLimit := uint64(21000)
	if err != nil {
		log.Println("Error estimating gas price")
		return "", err
	}

	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)

	chainID, err := s.EtherClient.NetworkID(context.Background())
	if err != nil {
		log.Println(err, ":Error getting NetworkID")
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println(err, "Error While signing transaction")
		return "", err
	}

	err = s.EtherClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err, "Error while sending transaction")
		return "", err
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
	return signedTx.Hash().Hex(), nil

}

func (s *erc20Service) TransferERC20ETH(transferReq *dto.TransferReq) (string, error) {
	privateKey, err := crypto.HexToECDSA(transferReq.PrivateKey)
	if err != nil {
		log.Println(err, "Error parsing HexToECDSA")
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("Address B ", fromAddress)

	nonce, err := s.EtherClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err, "Fail Checking Transaction Pending state")
		return "", err
	}

	// leave 5 USDT Amount of ETH For gas price
	transferReq.Amount = transferReq.Amount - 0.0042
	if transferReq.Amount <= 0 {
		log.Println("Not have 5 usdt equal balance of eth")
		return "", errors.New("we set the limit at lest eth balance equal to 5 USDT for gas price to avoid risky transaction")
	}
	value := big.NewInt(int64(transferReq.Amount * 1000000000000000000)) // in wei (1 eth) 1000000000000000000 = 1 eth
	gasLimit := uint64(21000)
	gasPrice, err := s.EtherClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err, "Error gettig suggestion gas price")
		return "", err
	}

	fmt.Println("AMOUNT", value)
	fmt.Println("Gas Limit", gasLimit)
	fmt.Println("Gas Price", gasPrice)

	toAddress := common.HexToAddress(transferReq.ToAddress)
	var data []byte

	if value.Cmp(gasPrice) != 1 {
		return "", errors.New("not enough ETH Balance")
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := s.EtherClient.NetworkID(context.Background())
	if err != nil {
		log.Println(err, "Error getting network ID")
		return "", nil
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println(err, "Error While signing transaction")
		return "", err
	}

	err = s.EtherClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err, "Error while sending transaction")
		return "", err
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil

}

// privatekey must be from approved address (bank)
func (s *erc20Service) TransferETHFromApprovedBankAddr(transferReq *dto.TransferReq) (string, error) {

	privateKey, err := crypto.HexToECDSA(transferReq.PrivateKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println(ok)
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	approvedAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("Address B ", approvedAddress)
	nonce, err := s.EtherClient.PendingNonceAt(context.Background(), approvedAddress)
	fromAddress := common.HexToAddress(transferReq.FromAddress)
	if err != nil {
		log.Println(err)
		return "", err
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := s.EtherClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("Gas Price ", gasPrice)

	toAddress := common.HexToAddress(transferReq.ToAddress)
	// contract address
	tokenAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7") // @todo 授权钱包地址,需要改成参数

	transferFnSignature := []byte("transferFrom(address,address,uint256)")
	hash := sha3.New256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	log.Println(hexutil.Encode(methodID))

	paddedAddressFrom := common.LeftPadBytes(fromAddress.Bytes(), 32)

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	log.Println(hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString(fmt.Sprintf("%v000000", transferReq.Amount), 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	log.Println(hexutil.Encode(paddedAmount))

	var data []byte
	//data = append(data, []byte("0x")...)
	data = append(data, methodID...)
	data = append(data, paddedAddressFrom...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := s.EtherClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Data:     data,
		Value:    value,
		GasPrice: gasPrice,
	})
	if err != nil {
		return "", err
	}
	log.Println("Gas Limit ", gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	log.Println(tx)
	//chainID, err := s.EtherClient.NetworkID(context.Background())
	chainID := big.NewInt(1) // chain id ==> prod
	if err != nil {
		log.Println(err)
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = s.EtherClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("tx sent: ", signedTx.Hash().Hex())

	return signedTx.Hash().Hex(), nil
}

func (s *erc20Service) ERC20CheckTransactionStatus(txID string) (*dto.ERC20TransStatusResp, error) {
	ticker := time.NewTicker(time.Second * 10)
	res := &dto.ERC20TransStatusResp{}

	var times int

	for range ticker.C {
		blockHash := common.HexToHash(txID)
		_, isPending, err := s.EtherClient.TransactionByHash(context.TODO(), blockHash)
		if err != nil {
			log.Println(err, "Error getting trasaction")
			times++
			if times > 20 {
				return nil, err
			}
			continue
		}
		fmt.Println(isPending, "Check Pending")
		if !isPending {
			receipt, err := s.EtherClient.TransactionReceipt(context.Background(), blockHash)
			if err != nil {
				log.Println(err, "Error getting transaction receipt")
				return nil, err
			}
			if receipt.Status == 1 {
				res.State = model.StateSuccess
				res.StateMessage = "Success"
			} else {
				res.State = model.StateFail
				res.StateMessage = "Fail"
			}
			return res, nil
		}

	}
	return res, nil

}
