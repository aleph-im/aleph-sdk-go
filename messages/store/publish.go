package store

import (
	"fmt"
	"github.com/aleph-im/aleph-sdk-go/accounts"
	"github.com/aleph-im/aleph-sdk-go/messages"
	"github.com/aleph-im/aleph-sdk-go/messages/create"
	"io"
	"time"
)

type StorePublishConfiguration struct {
	APIServer     string
	Account       accounts.Account
	StorageEngine messages.StorageEngine
	Channel       string
	File          io.Reader
}

type StoreContent struct {
	Address     string  `json:"address"`
	Time        float64 `json:"time"`
	ItemType    string  `json:"item_type"`
	ItemHash    string  `json:"item_hash"`
	Size        int64   `json:"size"`
	ContentType string  `json:"content_type"`
	Ref         string  `json:"ref,omitempty"`
}

func Publish(configuration StorePublishConfiguration) (string, error) {
	hash, err := create.PushFileToStorageEngine(create.FilePushConfiguration{
		APIServer:     configuration.APIServer,
		StorageEngine: configuration.StorageEngine,
		Key:           "file",
		Value:         configuration.File,
	})
	if err != nil {
		return "", fmt.Errorf("failed to push file into specified storage engine: %v", err)
	}

	timestamp := time.Now().Unix()
	content := StoreContent{
		Address:     configuration.Account.GetAddress(),
		Time:        float64(timestamp),
		ItemType:    configuration.StorageEngine,
		ItemHash:    hash,
		Size:        0,
		ContentType: "",
	}
	message := messages.BaseMessage{
		Channel:     configuration.Channel,
		Sender:      configuration.Account.GetAddress(),
		Chain:       configuration.Account.GetChain(),
		Type:        messages.MT_STORE,
		Time:        float64(timestamp),
		ItemType:    configuration.StorageEngine,
		ItemContent: "",
		ItemHash:    "",
	}

	err = create.PutContentToStorageEngine(create.PutContentConfiguration{
		Message:         &message,
		Content:         content,
		InlineRequested: true,
		StorageEngine:   configuration.StorageEngine,
		APIServer:       configuration.APIServer,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put content into specified storage engine: %v", err)
	}

	err = create.SignAndBroadcast(create.SignConfiguration{
		Account:   configuration.Account,
		Message:   &message,
		APIServer: configuration.APIServer,
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign and broadcast store message: %v", err)
	}
	return hash, err
}
