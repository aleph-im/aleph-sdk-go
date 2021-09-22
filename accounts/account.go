package accounts

import "aleph.im/aleph-sdk-go/messages"

// The Account interface is implemented by ETHAccount.
// In the future, it will be implemented by every account supported by Aleph.
// The Sign method is used to sign published messages.
type Account interface {
	GetAddress() string
	GetChain() messages.ChainType
	Sign(message *messages.BaseMessage) error
}
