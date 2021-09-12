package posts

import (
	"encoding/json"
	"testing"
)

const APIServer = "https://api2.aleph.im"

func TestPostMessageRetrievalNoFilters(t *testing.T) {
	pgc := PostGetterConfiguration{
		APIServer:  APIServer,
		Pagination: 10,
		Page:       1,
	}

	message, err := Get(pgc)
	if err != nil {
		t.Errorf("Failed to retrieve post message, reason: %v", err)
	}
	if message == nil {
		t.Errorf("Retrieved POST message should not be nil")
	} else {
		expected := 10
		messageNumber := len(message.Posts)
		if messageNumber != expected {
			t.Errorf("Invalid number of POST messages retrieved, expected [%d], but got [%d]", expected, messageNumber)
		}
	}
}

// Testing using https://explorer.aleph.im/address/ETH/0xF4c48E1B83233722F3609953EaF9800d0e3a1d8E/message/POST/a02179378327aa554f647f75a1dd5d62255fca8184a36ac66a1b3c3eaf343cd8
func TestPostMessageRetrievalWithFilters(t *testing.T) {
	pgc := PostGetterConfiguration{
		Types:      []string{"test_archetype"},
		APIServer:  APIServer,
		Pagination: 10,
		Page:       1,
		Hashes:     []string{"a02179378327aa554f647f75a1dd5d62255fca8184a36ac66a1b3c3eaf343cd8"},
	}

	message, err := Get(pgc)
	if err != nil {
		t.Errorf("Failed to retrieve post message, reason: %v", err)
	}
	if message == nil {
		t.Errorf("Retrieved POST message should not be nil")
	} else {
		expected := 1
		messageNumber := len(message.Posts)
		if messageNumber != expected {
			t.Errorf("Invalid number of POST messages retrieved, expected [%d], but got [%d]", expected, messageNumber)
		}

		type T struct {
			Name string `json:"name"`
			Description string `json:"description"`
		}
		placeholder := T{}
		marshal, err := json.Marshal(message.Posts[0].Content)
		if err != nil {
			t.Fatalf("Failed to marshal post message, reason: %v", err)
		}
		err = json.Unmarshal(marshal, &placeholder)
		if err != nil {
			t.Fatalf("Failed to unmarshal post message, reason: %v", err)
		}

		expectedDescription := "Je vous demande de vous arreter!"
		if placeholder.Description != expectedDescription {
			t.Errorf("Invalid post message description. Expected %s but got %s", expectedDescription, placeholder.Description)
		}
	}
}
