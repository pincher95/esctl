package cat

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type CatFielddataResponse struct {
	ID    string `json:"id"`
	Host  string `json:"host"`
	IP    string `json:"ip"`
	Node  string `json:"node"`
	Field string `json:"field"`
	Size  string `json:"size"`
}

func (c *cat) CatFielddata(endpoint, fields, bytes *string) (*[]CatFielddataResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/fielddata?format=json"

		if *fields != "" {
			*endpoint = fmt.Sprintf("_cat/fielddata/%s?format=json", *fields)
		}
	}

	if *bytes != "" {
		*endpoint += fmt.Sprintf("&bytes=%s", *bytes)
	}

	fieldData := make([]CatFielddataResponse, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&fieldData).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get field data: %s", resp.Status())
	}

	return &fieldData, nil
}
