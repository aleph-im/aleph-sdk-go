package ethereum

import (
	"encoding/hex"
	"fmt"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"aleph.im/aleph-sdk-go/messages"
)

const DefaultDerivationPath = "m/44'/60'/0'/0/0"

// The ETHAccount structure represents an ethereum account.
type ETHAccount struct {
	address   string
	publicKey string
	wallet    *hdwallet.Wallet
}

// NewAccount creates a new ethereum account using a derivation path.
// The default one you can use is DefaultDerivationPath.
//
// It creates the account using a generated mnemonic.
func NewAccount(derivationPath string) (*ETHAccount, error) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new mnemonic: %v", err)
	}
	return ImportAccountFromMnemonic(mnemonic, derivationPath)
}

// ImportAccountFromMnemonic imports an ethereum account using a given mnemonic.
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
		address:   a.Address.String(),
		publicKey: publicKey,
		wallet:    wallet,
	}, nil
}

// GetAddress implements the GetAddress method of the Account interface.
// It returns the account's address.
func (account *ETHAccount) GetAddress() string {
	return account.address
}

// GetChain returns the chain this account is attached to.
func (account *ETHAccount) GetChain() messages.ChainType {
	return messages.CT_ETH
}

// Sign implements the Sign method of the Account interface.
// It extracts parts of the message to be sign, and uses the ethereum wallet to sign the extracted buffer.
// it returns the resulting signature.
func (account *ETHAccount) Sign(message *messages.BaseMessage) error {
	buffer := messages.GetVerificationBuffer(message)

	signature, err := account.wallet.SignText(account.wallet.Accounts()[0], buffer)
	if err != nil {
		return fmt.Errorf("failed to sign message: %v", err)
	}
	message.Signature = "0x" + hex.EncodeToString(signature)
	return nil
}
