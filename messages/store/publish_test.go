package store

import (
	"os"
	"ptitluca.com/aleph-sdk-go/accounts/ethereum"
	"ptitluca.com/aleph-sdk-go/messages"
	"testing"
)

func TestPublish(t *testing.T) {
	account, err := ethereum.NewAccount(ethereum.DefaultDerivationPath)
	if err != nil {
		t.Fatalf("Faield to create new ethereum account, reason: %v", err)
	}

	file, err := os.Open("publish_test_content.txt")
	if err != nil {
		t.Fatalf("Failed to open file, reason: %v", err)
	}
	hash, err := Publish(StorePublishConfiguration{
		Channel:       "TEST",
		Account:       account,
		File:          file,
		StorageEngine: messages.SE_STORAGE,
		APIServer:     APIServer,
	})
	if err != nil {
		t.Errorf("Failed to publish store message, reason: %v", err)
	}

	response, err := Get(StoreGetConfiguration{
		FileHash:  hash,
		APIServer: APIServer,
	})
	if err != nil {
		t.Errorf("Failed to retrieve store message: %v", err)
	}

	expected := "Store message test !"
	if string(response) != expected {
		t.Errorf("Invalid store message retrieved. Expected %s, but got: %s", expected, string(response))
	}
}
