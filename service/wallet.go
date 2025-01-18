package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/hacpy/go-ethereum/common/hexutil"

	"blockchain-mining/types"

	"github.com/hacpy/go-ethereum/crypto"
)

func (s *Service) newWallet() (string, string, error) {
	p256 := elliptic.P256()

	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else if private == nil {
		return "", "", errors.New("private key is nil")
	} else {
		privateKeyBytes := crypto.FromECDSA(private)
		hexutil.Encode(privateKeyBytes)
		fmt.Println(hexutil.Encode(privateKeyBytes))
	}

	return "", "", nil
}

func (s *Service) MakeWallet() *types.Wallet {

	fmt.Println("들어옴")
	var wallet types.Wallet
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = s.newWallet(); err != nil {
		panic(err)
	} else {
		// TODO -> connect repository
		return &wallet
	}

	return &wallet
}
