package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func PostGraphQL(ctx context.Context, query string, queryName string, result interface{}) error {
	// dev.PrintJson(query)
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer, got %T", result)
	}

	url := ctx.Value(url).(string)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(query))
	req.Header.Set("Content-Type", "application/json")

	token := ctx.Value(token)
	if token != nil && token.(string) != "" {
		req.Header.Set("Authorization", "Bearer "+token.(string))
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

func LoginBDDUser(graphqlURL string) string {
	password := os.Getenv("BDD_USER_PASSWORD")
	email := os.Getenv("BDD_USER_EMAIL")
	query := fmt.Sprintf(`{ "query": "query { loginEmail(email: \"%s\", password: \"%s\") }" }`, email, password)

	ctx := context.WithValue(context.Background(), url, graphqlURL)
	var t *string
	err := PostGraphQL(ctx, query, "loginEmail", &t)
	if t == nil {
		panic("login failed: received nil token")
	}
	if err != nil {
		panic(fmt.Sprintf("login failed: %v", err))
	}

	return *t
}
