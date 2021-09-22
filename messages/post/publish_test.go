package posts

import (
	"encoding/json"
	"github.com/google/uuid"
	"aleph.im/aleph-sdk-go/accounts/ethereum"
	"aleph.im/aleph-sdk-go/messages"
	"testing"
)

func TestPublishPostMessage(t *testing.T) {
	account, err := ethereum.NewAccount(ethereum.DefaultDerivationPath)
	if err != nil {
		t.Fatalf("Failed to create new ethereum account, reason: %v", err)
	}

	content := map[string]string{
		"Body": "Hello World",
	}
	pType := uuid.NewString()
	configuration := PostPublishConfiguration{
		APIServer:       APIServer,
		Channel:         "TEST",
		InlineRequested: true,
		StorageEngine:   messages.SE_IPFS,
		Account:         account,
		PostType:        pType,
		Content:         content,
	}
	err = Publish(configuration)
	if err != nil {
		t.Errorf("Failed to publish post message: %v", err)
	}

	response, err := Get(PostGetterConfiguration{
		Types:      []string{pType},
		APIServer:  APIServer,
		Pagination: 200,
		Page:       1,
	})
	if err != nil {
		t.Errorf("Failed to retrieve post messages: %v", response)
	}

	if response == nil {
		t.Fatalf("Post message should not be nil")
	} else {
		if len(response.Posts) != 1 {
			t.Errorf("No post messages found for this type")
		}
	}
}

func TestAmendPostMessage(t *testing.T) {
	type T struct {
		Body string `json:"body"`
	}

	account, err := ethereum.NewAccount(ethereum.DefaultDerivationPath)
	if err != nil {
		t.Fatalf("Failed to create new ethereum account, reason: %v", err)
	}

	content := T{Body: "Hello World"}
	pType := uuid.NewString()
	configuration := PostPublishConfiguration{
		APIServer:       APIServer,
		Channel:         "TEST",
		InlineRequested: true,
		StorageEngine:   messages.SE_IPFS,
		Account:         account,
		PostType:        pType,
		Content:         content,
	}
	err = Publish(configuration)
	if err != nil {
		t.Errorf("Failed to publish post message: %v", err)
	}

	response, err := Get(PostGetterConfiguration{
		Types:      []string{pType},
		APIServer:  APIServer,
		Pagination: 200,
		Page:       1,
	})
	if err != nil {
		t.Errorf("Failed to retrieve post messages: %v", response)
	}

	if response == nil {
		t.Fatalf("Post message should not be nil")
	} else {
		newContent := T{Body: "Hello World V2"}
		err = Publish(PostPublishConfiguration{
			APIServer:       APIServer,
			Channel:         "TEST",
			InlineRequested: true,
			StorageEngine:   messages.SE_IPFS,
			Account:         account,
			PostType:        "amend",
			Content:         newContent,
			Ref: response.Posts[0].ItemHash,
		})
		if err != nil {
			t.Errorf("Failed to publish post message: %v", err)
		}

		response, err = Get(PostGetterConfiguration{
			Types:      []string{pType},
			APIServer:  APIServer,
			Pagination: 200,
			Page:       1,
		})
		if err != nil {
			t.Errorf("Failed to retrieve amended post messages: %v", response)
		}
		if response == nil {
			t.Fatalf("Post message should not be nil")
		} else {
			expected := T{
				Body: "Hello World V2",
			}
			placeholder := T{}

			marshal, err := json.Marshal(response.Posts[0].Content)
			if err != nil {
				t.Fatalf("Failed to marshal content: %v", err)
			}
			err = json.Unmarshal(marshal, &placeholder)
			if err != nil {
				t.Fatalf("Failed to unmarshal content: %v", err)
			}
			if placeholder.Body != expected.Body {
				t.Errorf("Invalid post message content afer amend. Expected %s, but got %s",
					response.Posts[0].Content, expected)
			}
		}
	}
}