package posts

import (
	"fmt"
	"strings"

	"github.com/imroc/req"
)

type PostGetterResponseBody struct {
	ID                map[string]string `json:"_id"`
	Chain             string            `json:"chain"`
	ItemHash          string            `json:"item_hash"`
	Sender            string            `json:"sender"`
	Type              string            `json:"type"`
	Channel           string            `json:"channel"`
	Confirmed         bool              `json:"confirmed"`
	Content           interface{}       `json:"content"`
	ItemContent       string            `json:"item_content"`
	ItemType          string            `json:"item_type"`
	Signature         string            `json:"signature"`
	Size              uint64            `json:"size"`
	Time              float64           `json:"time"`
	OriginalItemHash  string            `json:"original_item_hash"`
	OriginalSignature string            `json:"original_signature"`
	OriginalType      string            `json:"original_type"`
	Hash              string            `json:"hash"`
	Address           string            `json:"address"`
}

type PostGetterResponse struct {
	Posts             []PostGetterResponseBody `json:"posts"`
	PaginationPage    uint64                   `json:"pagination_page"`
	PaginationTotal   uint64                   `json:"pagination_total"`
	PaginationPerPage uint64                   `json:"pagination_per_page"`
	PaginationItem    string                   `json:"pagination_item"`
}

type PostGetterConfiguration struct {
	Types      []string
	APIServer  string
	Pagination uint64
	Page       uint64
	Refs       []string
	Addresses  []string
	Tags       []string
	Hashes     []string
}

// Get retrieves a post message from the Aleph network.
func Get(pgc PostGetterConfiguration) (*PostGetterResponse, error) {
	param := req.Param{
		"pagination": pgc.Pagination,
		"page":       pgc.Page,
	}

	if len(pgc.Types) > 0 {
		param["types"] = strings.Join(pgc.Types, ",")
	}
	if len(pgc.Refs) > 0 {
		param["refs"] = strings.Join(pgc.Refs, ",")
	}
	if len(pgc.Addresses) > 0 {
		param["addresses"] = strings.Join(pgc.Addresses, ",")
	}
	if len(pgc.Tags) > 0 {
		param["tags"] = strings.Join(pgc.Tags, ",")
	}
	if len(pgc.Hashes) > 0 {
		param["hashes"] = strings.Join(pgc.Hashes, ",")
	}

	requester := req.New()

	response, err := requester.Get(pgc.APIServer+"/api/v0/posts.json", param)
	if err != nil {
		return nil, fmt.Errorf("GET request has failed: %v", err)
	}

	buffer := &PostGetterResponse{}
	err = response.ToJSON(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response to JSON: %v", err)
	}
	return buffer, nil
}
