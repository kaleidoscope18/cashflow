package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// todo automate extracting the query name from the query
func PostGraphQL(query string, queryName string, token string, result interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer, got %T", result)
	}

	jsonValue, _ := json.Marshal(map[string]string{"query": query})
	req, _ := http.NewRequest("POST", os.Getenv("GRAPHQL_URL"), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("problem status %d - %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		formattedError := fmt.Sprintf("failed to read response body: %v", err)
		log.Println(formattedError)
		return fmt.Errorf(formattedError)
	}

	if strings.Contains(string(body), `"errors"`) {
		var response map[string][]map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}

		if errors, ok := response["errors"]; ok {
			return fmt.Errorf(errors[0]["message"].(string))
		}
	}
	if result == nil {
		return nil
	}
	err = unmarshalWrappedBody(body, queryName, result)
	if err != nil {
		formattedError := fmt.Sprintf("failed to parse JSON response: %s - %v", string(body), err)
		log.Println(formattedError)
		return fmt.Errorf(formattedError)
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
