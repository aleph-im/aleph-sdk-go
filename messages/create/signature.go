package create

import (
	"encoding/json"
	"fmt"

	"github.com/imroc/req"
	"ptitluca.com/aleph-sdk-go/accounts"
	"ptitluca.com/aleph-sdk-go/messages"
)

type SignConfiguration struct {
	Account   accounts.Account
	Message   *messages.BaseMessage
	APIServer string
}

type BroadcastConfiguration struct {
	Message   *messages.BaseMessage
	APIServer string
}

type BroadcastData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

func Broadcast(bc BroadcastConfiguration) error {
	serialized, err := json.Marshal(bc.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	data := BroadcastData{
		Topic: "ALEPH-TEST",
		Data:  string(serialized),
	}
	serializedData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	requester := req.New()

	_, err = requester.Post(bc.APIServer+"/api/v0/ipfs/pubsub/pub", serializedData)
	if err != nil {
		return fmt.Errorf("POST request has failed: %v", err)
	}
	return nil
}

func SignAndBroadcast(sc SignConfiguration) error {
	if sc.Message.Chain == "" {
		sc.Message.Chain = sc.Account.GetChain()
	}

	err := sc.Account.Sign(sc.Message)
	if err != nil {
		return fmt.Errorf("failed to sign message: %v", err)
	}

	bc := BroadcastConfiguration{
		Message:   sc.Message,
		APIServer: sc.APIServer,
	}
	err = Broadcast(bc)
	if err != nil {
		return fmt.Errorf("failed to broadcast message: %v", err)
	}
	return nil
}
