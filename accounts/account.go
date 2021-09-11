package accounts

import "ptitluca.com/aleph-sdk-go/messages"

type Account interface {
	GetAddress() string
	GetChain() messages.ChainType
	Sign(message *messages.BaseMessage) error
}