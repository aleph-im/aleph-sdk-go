package posts

import (
	"fmt"
	"github.com/aleph-im/aleph-sdk-go/accounts"
	"github.com/aleph-im/aleph-sdk-go/messages"
	"github.com/aleph-im/aleph-sdk-go/messages/create"
	"time"
)

// PostPublishConfiguration is used while publishing post messages.
type PostPublishConfiguration struct {
	APIServer       string                 // The API Server endpoint used to handle the request.
	Account         accounts.Account       // The account used to sign the message before broadcasting it.
	StorageEngine   messages.StorageEngine // The storage engine - IPFS or Aleph - used to store the content.
	Channel         string                 // The targeted channel to store the content.
	InlineRequested bool                   // Will the content be inlined ?
	Content         interface{}            // The content stored.
	PostType        string                 // The post type attached to the content (used as a key).
	Ref             string                 // The reference used if amending a message.
}

type ReferencedPostContent struct {
	PostContent
	Ref string `json:"ref,omitempty"`
}

type PostContent struct {
	Address string      `json:"address"`
	Time    float64     `json:"time"`
	Content interface{} `json:"content"`
	Type    string      `json:"type"`
}

// Publish uses the post type - i.e. the key - and value provided in the configuration to publish a post message on the
// Aleph network.
func Publish(configuration PostPublishConfiguration) (*messages.BaseMessage, error) {
	timestamp := time.Now().Unix()
	content := PostContent{
		Address: configuration.Account.GetAddress(),
		Time:    float64(timestamp),
		Content: configuration.Content,
		Type:    configuration.PostType,
	}

	message := messages.BaseMessage{
		Channel:  configuration.Channel,
		Sender:   configuration.Account.GetAddress(),
		Chain:    configuration.Account.GetChain(),
		Type:     messages.MT_POST,
		Time:     float64(timestamp),
		ItemType: configuration.StorageEngine,
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
		return nil, fmt.Errorf("failed to put content into storage engine: %v", err)
	}

	sc := create.SignConfiguration{
		Account:   configuration.Account,
		Message:   &message,
		APIServer: configuration.APIServer,
	}
	err = create.SignAndBroadcast(sc)
	if err != nil {
		return nil, fmt.Errorf("failed to sign and broadcast post message: %v", err)
	}
	return &message, nil
}
