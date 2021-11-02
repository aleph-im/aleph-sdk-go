package aggregate

import (
	"fmt"
	"github.com/aleph-im/aleph-sdk-go/accounts"
	"github.com/aleph-im/aleph-sdk-go/messages"
	"github.com/aleph-im/aleph-sdk-go/messages/create"
	"time"
)

type AggregateContent struct {
	Address string      `json:"address"`
	Key     string      `json:"key"`
	Content interface{} `json:"content"`
	Time    uint64      `json:"time"`
}

// AggregatePublishConfiguration is used while publishing aggregate messages.
type AggregatePublishConfiguration struct {
	APIServer       string                 // The API Server endpoint used to handle the request.
	Account         accounts.Account       // The account used to sign the message before broadcasting it.
	StorageEngine   messages.StorageEngine // The storage engine - IPFS or Aleph - used to store the content.
	Channel         string                 // The targeted channel to store the content.
	InlineRequested bool                   // Will the content be inlined ?
	Content         interface{}            // The content stored.
	Key             string                 // The key attached to the content.
}

// Publish uses the key and value provided in the configuration to publish an aggregate message on the Aleph network.
func Publish(asc AggregatePublishConfiguration) (*messages.BaseMessage, error) {
	timestamp := time.Now().Unix()
	content := AggregateContent{
		Address: asc.Account.GetAddress(),
		Key:     asc.Key,
		Content: asc.Content,
		Time:    uint64(timestamp),
	}

	msg := messages.BaseMessage{
		Channel: asc.Channel,
		Chain:   asc.Account.GetChain(),
		Sender:  asc.Account.GetAddress(),
		Type:    messages.MT_AGGREGATE,
		Time:    float64(timestamp),
	}

	pcc := create.PutContentConfiguration{
		Message:         &msg,
		Content:         content,
		InlineRequested: asc.InlineRequested,
		StorageEngine:   asc.StorageEngine,
		APIServer:       asc.APIServer,
	}
	err := create.PutContentToStorageEngine(pcc)
	if err != nil {
		return nil, fmt.Errorf("failed to put content into storage engine: %v", err)
	}

	sc := create.SignConfiguration{
		Account:   asc.Account,
		Message:   &msg,
		APIServer: asc.APIServer,
	}
	err = create.SignAndBroadcast(sc)
	if err != nil {
		return nil, fmt.Errorf("failed to sign and broadcast msg: %v", err)
	}
	return &msg, nil
}
