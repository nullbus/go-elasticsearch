package es

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Cluster struct {
	URL url.URL
}

func NewCluster(address string) (*Cluster, error) {
	parsed, parseErr := url.Parse(address)
	if nil != parseErr {
		return nil, parseErr
	}

	return &Cluster{*parsed}, nil
}

func (c *Cluster) Search(search Search) (*http.Response, error) {
	url := c.URL
	url.Path = search.Path()
	query := search.Query()

	if search.Type() != SEARCH_TYPE_SCROLL {
		query["search_type"] = []string{search.Type().String()}
	}
	url.RawQuery = query.Encode()

	var (
		req       *http.Request
		createErr error
	)

	body := search.Data()
	if body == nil {
		req, createErr = http.NewRequest("GET", url.String(), nil)
	} else {
		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(body); err != nil {
			return nil, err
		}

		req, createErr = http.NewRequest("POST", url.String(), &buffer)
	}

	if nil != createErr {
		return nil, createErr
	}

	return (&http.Client{}).Do(req)
}
