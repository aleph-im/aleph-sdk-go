package aggregate

import (
	"fmt"
	"time"

	"aleph.im/aleph-sdk-go/accounts"
	"aleph.im/aleph-sdk-go/messages"
	"aleph.im/aleph-sdk-go/messages/create"
)

type AggregateContent struct {
	Address string      `json:"address"`
	Key     string      `json:"key"`
	Content interface{} `json:"content"`
	Time    uint64      `json:"time"`
}

// AggregatePublishConfiguration is used while publishing aggregate messages.
//
// APIServer	- The API Server endpoint used to handle the request.
// Account		- The account used to sign the message before broadcasting it.
// StorageEngine	- The storage engine - IPFS or Aleph - used to store the content.
// Channel		- The targeted channel to store the content.
// InlineRequest	- Will the content be inlined ?
// Content		- The content stored.
// Key			- The key attached to the content.
type AggregatePublishConfiguration struct {
	APIServer       string
	Account         accounts.Account
	StorageEngine   messages.StorageEngine
	Channel         string
	InlineRequested bool
	Content         interface{}
	Key             string
}

// Publish uses the key and value provided in the configuration to publish an aggregate message on the Aleph network.
func Publish(asc AggregatePublishConfiguration) error {
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
		return fmt.Errorf("failed to put content into storage engine: %v", err)
	}

	sc := create.SignConfiguration{
		Account:   asc.Account,
		Message:   &msg,
		APIServer: asc.APIServer,
	}
	err = create.SignAndBroadcast(sc)
	if err != nil {
		return fmt.Errorf("failed to sign and broadcast msg: %v", err)
	}
	return nil
}
