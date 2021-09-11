package ethereum

import (
	"encoding/hex"
	"fmt"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"ptitluca.com/aleph-sdk-go/messages"
)

const DefaultDerivationPath = "m/44'/60'/0'/0/0"

type ETHAccount struct {
	address string
	publicKey string
	wallet *hdwallet.Wallet
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
	a, err := wallet.Derive(path, true)
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
		wallet: wallet,
	}, nil
}

func (account *ETHAccount) GetAddress() string {
	return account.address
}

func (account *ETHAccount) GetChain() messages.ChainType {
	return messages.CT_ETH
}

func (account *ETHAccount) Sign(message *messages.BaseMessage) error {
	buffer := messages.GetVerificationBuffer(message)

	signature, err := account.wallet.SignText(account.wallet.Accounts()[0], buffer)
	if err != nil {
		return fmt.Errorf("failed to sign message: %v", err)
	}
	message.Signature = "0x" + hex.EncodeToString(signature)
	return nil
}