package store

import (
	"fmt"

	"github.com/imroc/req"
)

type StoreGetConfiguration struct {
	FileHash  string
	APIServer string
}

func Get(sgc StoreGetConfiguration) ([]byte, error) {
	requester := req.New()

	response, err := requester.Get(sgc.APIServer + "/api/v0/storage/raw/" + sgc.FileHash + "?find")
	if err != nil {
		return nil, fmt.Errorf("GET request has failed: %v", err)
	}
	return response.Bytes(), nil
}
