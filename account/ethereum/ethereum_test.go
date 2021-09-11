package ethereum

import (
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"testing"
)

func TestAccountCreation(t *testing.T) {
	account, err := NewAccount(DefaultDerivationPath)
	if err != nil {
		t.Errorf("Failed to create new ethereum account, reason: %v", err)
	}

	if account == nil {
		t.Errorf("Ethereum account should not be nil")
	} else {
		if account.publicKey == "" {
			t.Errorf("Ethereum public key should not be empty")
		}
		if account.address == "" {
			t.Errorf("Ethereum address should not be empty")
		}
	}
}

func TestAccountImport(t *testing.T) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		t.Errorf("Failed to generate new mnemonic: %v", err)
	}
	account, err := ImportAccountFromMnemonic(mnemonic, DefaultDerivationPath)
	if err != nil {
		t.Errorf("Failed to create new ethereum account, reason: %v", err)
	}

	if account == nil {
		t.Errorf("Ethereum account should not be nil")
	} else {
		if account.publicKey == "" {
			t.Errorf("Ethereum public key should not be empty")
		}
		if account.address == "" {
			t.Errorf("Ethereum address should not be empty")
		}
	}
}