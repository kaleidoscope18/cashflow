package bdd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PostGraphQL(url string, query string, queryName string, result interface{}) error {
	fmt.Println(query)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(query))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("problem status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if result == nil {
		return nil
	}
	err = unmarshalWrappedBody(body, queryName, result)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return nil
}

func unmarshalWrappedBody(body []byte, queryName string, target interface{}) error {
	var outerMap map[string]map[string]json.RawMessage
	err := json.Unmarshal(body, &outerMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal outer JSON: %v", err)
	}

	result, ok := outerMap["data"][queryName]
	if !ok {
		return fmt.Errorf("%s not found in body data %s", queryName, outerMap)
	}

	if err := json.Unmarshal(result, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target struct: %v", err)
	}

	return nil
}
