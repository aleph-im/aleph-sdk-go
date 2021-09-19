package aggregate

import (
	"fmt"
	"time"

	"ptitluca.com/aleph-sdk-go/accounts"
	"ptitluca.com/aleph-sdk-go/messages"
	"ptitluca.com/aleph-sdk-go/messages/create"
)

type AggregateContent struct {
	Address string      `json:"address"`
	Key     string      `json:"key"`
	Content interface{} `json:"content"`
	Time    uint64      `json:"time"`
}

type AggregatePublishConfiguration struct {
	Account         accounts.Account
	Key             string
	Content         interface{}
	Chain           messages.ChainType
	Channel         string
	APIServer       string
	InlineRequested bool
	StorageEngine   messages.StorageEngine
}

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
		Chain:   asc.Chain,
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
