package posts

import (
	"fmt"
	"ptitluca.com/aleph-sdk-go/accounts"
	"ptitluca.com/aleph-sdk-go/messages"
	"ptitluca.com/aleph-sdk-go/messages/create"
	"time"
)

type PostPublishConfiguration struct {
	APIServer string
	Ref string
	Channel string
	InlineRequested bool
	StorageEngine messages.StorageEngine
	Account accounts.Account
	PostType string
	Content interface{}
}

type ReferencedPostContent struct {
	PostContent
	Ref string `json:"ref,omitempty"`
}

type PostContent struct {
	Address string `json:"address"`
	Time float64 `json:"time"`
	Content interface{} `json:"content"`
	Type string `json:"type"`
}

func Publish(configuration PostPublishConfiguration) error {
	timestamp := time.Now().Unix()
	content := PostContent{
		Address: configuration.Account.GetAddress(),
		Time: float64(timestamp),
		Content: configuration.Content,
		Type:    configuration.PostType,
	}

	message := messages.BaseMessage{
		Channel:       configuration.Channel,
		Sender:        configuration.Account.GetAddress(),
		Chain:         configuration.Account.GetChain(),
		Type:          messages.MT_POST,
		Time: float64(timestamp),
		ItemType:      configuration.StorageEngine,
	}

	pcc := create.PutContentConfiguration{
		Message:         &message,
		Content:         content,
		InlineRequested: configuration.InlineRequested,
		StorageEngine:   configuration.StorageEngine,
		APIServer:       configuration.APIServer,
	}

	if configuration.Ref != "" {
		refContent := ReferencedPostContent{
			PostContent: content,
			Ref:         configuration.Ref,
		}
		pcc.Content = refContent
	}

	err := create.PutContentToStorageEngine(pcc)
	if err != nil {
		return fmt.Errorf("failed to put content into storage engine: %v", err)
	}

	sc := create.SignConfiguration{
		Account:   configuration.Account,
		Message:   &message,
		APIServer: configuration.APIServer,
	}
	err = create.SignAndBroadcast(sc)
	if err != nil {
		return fmt.Errorf("failed to sign and broadcast post message: %v", err)
	}
	return nil
}
