package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"

	"blockchain-mining/types"

	"github.com/hacpy/go-ethereum/common/hexutil"
	"github.com/hacpy/go-ethereum/crypto"
)

func (s *Service) newKeyPair() (string, string, error) {
	p256 := elliptic.P256()

	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else if private == nil {
		return "", "", errors.New(types.PkNil)
	} else {
		privateKeyBytes := crypto.FromECDSA(private)
		privateKey := hexutil.Encode(privateKeyBytes)

		importedPrivateKey, err := crypto.HexToECDSA(privateKey[2:])

		if err != nil {
			return "", "", err
		}

		coordPublicKey := importedPrivateKey.Public()
		publicKeyECDSA, ok := coordPublicKey.(*ecdsa.PublicKey)

		if !ok {
			return "", "", errors.New("error casting public key type")
		}

		publicKeyBytes := crypto.PubkeyToAddress(*publicKeyECDSA)
		publicKey := hexutil.Encode(publicKeyBytes[:])
		return privateKey, publicKey, nil
	}
}

func (s *Service) MakeWallet() *types.Wallet {
	var wallet types.Wallet
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = s.newKeyPair(); err != nil {
		return nil
	} else if err = s.repository.CreateNewWallet(&wallet); err != nil {
		return nil
	} else {
		return &wallet
	}
}

func (s *Service) GetWallet(pk string) (*types.Wallet, error) {
	if wallet, err := s.repository.GetWallet(pk); err != nil {
		return nil, err
	} else {
		return wallet, nil
	}
}
