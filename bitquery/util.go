package bitquery

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"ps-chartdata/config"
)

const (
	BITQUERY_URL = "https://graphql.bitquery.io"
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
	req.Header.Set("X-API-KEY", config.Conf.BitqueryAPIKey)

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
