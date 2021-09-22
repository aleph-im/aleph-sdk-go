package create

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/imroc/req"
	"aleph.im/aleph-sdk-go/messages"
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

type FilePushConfiguration struct {
	APIServer     string
	StorageEngine messages.StorageEngine
	Key           string
	Value         io.Reader
}

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

func PushFileToStorageEngine(configuration FilePushConfiguration) (string, error) {
	var buffer bytes.Buffer
	var fileWriter io.Writer
	var err error

	writer := multipart.NewWriter(&buffer)
	if x, ok := configuration.Value.(io.Closer); ok {
		defer x.Close()
	}

	if x, ok := configuration.Value.(*os.File); ok {
		if fileWriter, err = writer.CreateFormFile(configuration.Key, x.Name()); err != nil {
			return "", fmt.Errorf("failed to create form file: %v", err)
		}
	} else {
		if fileWriter, err = writer.CreateFormField(configuration.Key); err != nil {
			return "", fmt.Errorf("failed to create form field: %v", err)
		}
	}

	if _, err = io.Copy(fileWriter, configuration.Value); err != nil {
		return "", fmt.Errorf("failed to perform copy: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	url := configuration.APIServer + "/api/v0/" + strings.ToLower(configuration.StorageEngine) + "/add_file"
	requester := req.New()
	header := make(http.Header)
	header.Set("Content-Type", writer.FormDataContentType())
	response, err := requester.Post(url, &buffer, header)
	if err != nil {
		return "", fmt.Errorf("POST request has failed: %v", err)
	}

	placeholder := &PushResponse{}
	err = json.Unmarshal(response.Bytes(), &placeholder)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal store response: %v", err)
	}
	return placeholder.Hash, nil
}
