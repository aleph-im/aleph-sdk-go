package ethereum

import (
	"fmt"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const DefaultDerivationPath = "m/44'/60'/0'/0/0"

type ETHAccount struct {
	address string
	publicKey string
}

func NewAccount(derivationPath string) (*ETHAccount, error) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new mnemonic: %v", err)
	}
	return ImportAccountFromMnemonic(mnemonic, derivationPath)
}

func ImportAccountFromMnemonic(mnemonic, derivationPath string) (*ETHAccount, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from mnemonic: %v", err)
	}

	path := hdwallet.MustParseDerivationPath(derivationPath)
	a, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %v", err)
	}

	publicKey, err := wallet.PublicKeyHex(a)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve public key from ethereum accout: %v", err)
	}
	return &ETHAccount{
		address: a.Address.String(),
		publicKey: publicKey,
	}, nil
}
