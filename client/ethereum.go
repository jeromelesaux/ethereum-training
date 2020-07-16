package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"sync"

	ecb "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	ec "github.com/ethereum/go-ethereum/ethclient"
	"github.com/jeromelesaux/ethereum-training/config"
)

var (
	Auth           *ecb.TransactOpts
	EthClient      *ec.Client
	authenticateDo sync.Once
	SafeNonceTx    *SafeNonce
	PrivateKey     *ecdsa.PrivateKey
)

type SafeNonce struct {
	NonceMutex sync.Mutex
	Nonce      uint64
}

func connect() (auth *ecb.TransactOpts, clt *ec.Client) {
	conf := config.MyConfig

	clt, err := ec.Dial(conf.EthereumEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not connect to Ethereum with error %v\n", err)
		log.Fatalf("Can not connect to Ethereum with error %v\n", err)
	}

	PrivateKey, err = pkToECDSA()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not convert private key  with error %v\n", err)
		log.Fatalf("Can not convert private key with error %v\n", err)
	}

	auth = ecb.NewKeyedTransactor(PrivateKey)

	fmt.Fprintf(os.Stdout, "Last Nonce %d\n", auth.Nonce)
	return
}

func Authenticate() {
	authenticateDo.Do(
		func() {
			Auth, EthClient = connect()
			SafeNonceTx.NonceMutex.Lock()
			nonce, err := EthClient.PendingNonceAt(context.Background(), Auth.From)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while getting nonce from Ethereum with error:%v\n", err)
			}
			SafeNonceTx.Nonce = nonce
			fmt.Fprintf(os.Stdout, "Last Nonce :%d\n", nonce)
			SafeNonceTx.NonceMutex.Unlock()
		})

}

func pkToECDSA() (*ecdsa.PrivateKey, error) {
	c := config.MyConfig
	privateKey, err := crypto.HexToECDSA(c.PrivateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while convert to ECDSA our private key with error :%v\n", err)
	}
	return privateKey, err
}
