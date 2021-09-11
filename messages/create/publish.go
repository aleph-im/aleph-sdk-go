package create

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"net/http"
	"ptitluca.com/aleph-sdk-go/messages"
	"strings"
)

type PushResponse struct {
	Hash string `json:"hash"`
}

type StandardPushConfiguration struct {
	Value         interface{}
	APIServer     string
	StorageEngine messages.StorageEngine
}

type PutContentConfiguration struct {
	Message         *messages.BaseMessage
	Content         interface{}
	InlineRequested bool
	StorageEngine   messages.StorageEngine
	APIServer       string
}

// PushToStorageEngine sends the provided content to the selected storage engine (given in the configuration).
func PushToStorageEngine(spc StandardPushConfiguration) (*PushResponse, error) {
	url := spc.APIServer + "/api/v0/" + strings.ToLower(spc.StorageEngine) + "/add_json"
	requester := req.New()

	serialized, err := json.Marshal(spc.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %v", err)
	}
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	response, err := requester.Post(url, serialized, header)
	if err != nil {
		return nil, fmt.Errorf("POST request has failed: %v", err)
	}

	buffer := &PushResponse{}
	err = response.ToJSON(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %v", err)
	}
	return buffer, nil
}

func PutContentToStorageEngine(pcc PutContentConfiguration) error {
	if pcc.InlineRequested {
		serialized, err := json.Marshal(pcc.Content)
		if err != nil {
			return fmt.Errorf("failed to marshal content: %v", err)
		}

		if len(serialized) > 150000 {
			pcc.InlineRequested = false
		} else {
			pcc.Message.ItemType = messages.IT_INLINE
			pcc.Message.ItemContent = string(serialized)

			hasher := sha256.New()
			_, err := hasher.Write(serialized)
			if err != nil {
				return fmt.Errorf("failed to write: %v", err)
			}

			pcc.Message.ItemHash = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	if !pcc.InlineRequested {
		pcc.Message.ItemType = pcc.StorageEngine

		spc := StandardPushConfiguration{
			Value:         pcc.Content,
			APIServer:     pcc.APIServer,
			StorageEngine: pcc.StorageEngine,
		}
		response, err := PushToStorageEngine(spc)
		if err != nil {
			return fmt.Errorf("failed to push to desired storage engine: %v", err)
		}
		pcc.Message.ItemHash = response.Hash
	}
	return nil
}
