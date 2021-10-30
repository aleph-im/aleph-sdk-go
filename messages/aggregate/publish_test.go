package aggregate

import (
	"encoding/json"
	"github.com/aleph-im/aleph-sdk-go/accounts/ethereum"
	"github.com/aleph-im/aleph-sdk-go/messages"
	"testing"
)

func TestAggregateSubmit(t *testing.T) {
	type T struct {
		A int `json:"A"`
		B int `json:"B"`
	}
	account, err := ethereum.NewAccount(ethereum.DefaultDerivationPath)
	if err != nil {
		t.Fatalf("Failed to create new ethereum account: %v", err)
	}

	content := T{
		A: 1,
		B: 2,
	}
	asc := AggregatePublishConfiguration{
		Account:         account,
		Key:             "satoshi",
		Content:         content,
		Channel:         "TEST",
		APIServer:       APIServer,
		InlineRequested: true,
		StorageEngine:   messages.SE_STORAGE,
	}
	err = Publish(asc)
	if err != nil {
		t.Errorf("Failed to submit aggregate message: %v", err)
	}

	sgc := AggregateGetConfiguration{
		APIServer: APIServer,
		Address:   account.GetAddress(),
		Keys:      []string{"satoshi"},
	}
	msg, err := Get(sgc)
	if err != nil {
		t.Errorf("Failed to retrieve store msg, reason: %v", err)
	}
	if msg == nil {
		t.Errorf("Aggregate msg should not be nil")
	} else {
		placeholder := map[string]T{}
		err = json.Unmarshal(msg, &placeholder)
		if err != nil {
			t.Errorf("Failed to unmarshal aggregate message, reason: %v", err)
		}

		if placeholder["satoshi"].A != 1 && placeholder["satoshi"].B != 2 {
			t.Errorf("Failed to retrieve value of A or B.")
		}
	}
}
