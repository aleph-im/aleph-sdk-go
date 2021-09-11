package messages

type StorageEngine = string

const (
	SE_IPFS    = "IPFS"
	SE_STORAGE = "STORAGE"
)

type ChainType = string

const (
	CT_ETH ChainType = "ETH"
	CT_SOL ChainType = "SOL"
)

type MessageType = string

const (
	MT_AGGREGATE MessageType = "AGGREGATE"
)

type ItemType = string

const (
	IT_INLINE  ItemType = "INLINE"
	IT_IPFS    ItemType = "IPFS"
	IT_STORAGE ItemType = "STORAGE"
)

type MessageConfirmation struct {
	Chain string `json:"chain"`
	Height int64 `json:"height"`
	Hash interface{} `json:"hash"`
}

type BaseMessageContent struct {
	Address string `json:"address"`
	Time float64 `json:"time"`
}

type BaseMessage struct {
	ID			map[string]string `json:"_id"`
	Channel     string      `json:"channel"`
	Sender      string      `json:"sender"`
	Chain       ChainType   `json:"chain"`
	Type        MessageType `json:"type"`
	Time        float64     `json:"time"`
	ItemType    ItemType    `json:"item_type,omitempty"`
	ItemContent string      `json:"item_content"`
	ItemHash    string      `json:"item_hash"`
	Signature   string      `json:"signature"`
	Confirmations []MessageConfirmation `json:"confirmations,omitempty"`
	Confirmed bool `json:"confirmed"`
	Size uint64 `json:"size"`
	HashType string `json:"hash_type,omitempty"`
	Content BaseMessageContent `json:"content,omitempty"`
}

func GetVerificationBuffer(msg *BaseMessage) []byte {
	buffer := msg.Chain + "\n" + msg.Sender + "\n" + msg.Type + "\n" + msg.ItemHash
	return []byte(buffer)
}