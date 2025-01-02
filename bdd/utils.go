package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func PostGraphQL(ctx context.Context, query string, queryName string, result interface{}) error {
	fmt.Println(query)
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
	fmt.Println(string(body))
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if result == nil {
		return nil
	}
	err = unmarshalWrappedBody(body, queryName, result)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %s - %v", string(body), err)
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
