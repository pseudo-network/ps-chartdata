package bitquery

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	BITQUERY_URL     = "https://graphql.bitquery.io"
	BITQUERY_API_KEY = "BQYug1u2azt1EzuPggXfnhdhzFObRW0g"
)

func Query(query string) ([]byte, error) {

	reqBody, err := json.Marshal(map[string]string{
		"query": query,
	})

	b := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest("POST", BITQUERY_URL, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", BITQUERY_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
