package store

import (
	"testing"
)

const APIServer = "https://api2.aleph.im"

func TestStoreMessageRetrieval(t *testing.T) {
	response, err := Get(StoreGetConfiguration{
		FileHash:  "QmQkv43jguT5HLC8TPbYJi2iEmr4MgLgu4nmBoR4zjYb3L",
		APIServer: APIServer,
	})
	if err != nil {
		t.Errorf("Failed to retrieve store message: %v", err)
	}

	expected := "This is just a test."
	if string(response) != expected {
		t.Errorf("Invalid store message retrieved. Expected %s, but got: %s", expected, string(response))
	}
}
