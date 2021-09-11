package aggregate

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"strings"
)

type AggregateGetResponse = []byte

type AggregateGetConfiguration struct {
	Address   string
	Keys      []string
	APIServer string
}

// Get retrieves AGGREGATE messages on the Aleph network using filters provided in the configuration.
func Get(agc AggregateGetConfiguration) (AggregateGetResponse, error) {
	param := req.Param{}
	if len(agc.Keys) > 0 {
		param["keys"] = strings.Join(agc.Keys, ",")
	}

	requester := req.New()

	response, err := requester.Get(agc.APIServer+"/api/v0/aggregates/"+agc.Address+".json", param)
	if err != nil {
		return nil, fmt.Errorf("GET request has failed: %v", err)
	}

	type T struct {
		Address string `json:"address"`
		Data map[string]interface{} `json:"data"`
	}
	placeholder := T{}
	err = json.Unmarshal(response.Bytes(), &placeholder)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregate message: %v", err)
	}

	bytes, err := json.Marshal(placeholder.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal aggregate response: %v", err)
	}
	return bytes, nil
}

// GetProfile uses Get to retrieve the aggregate message containing the profile of a user using the key filter.
func GetProfile(agc AggregateGetConfiguration) (AggregateGetResponse, error) {
	response, err := Get(AggregateGetConfiguration{
		Address:   agc.Address,
		Keys:      []string{"profile"},
		APIServer: agc.APIServer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve aggregate message: %v", err)
	}
	return response, nil
}