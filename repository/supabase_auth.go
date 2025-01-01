package repository

import (
	"bytes"
	"cashflow/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (repo *supabaseRepository) LoginEmail(ctx context.Context, email string, password string) (string, error) {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	req, err := http.NewRequest("POST", os.Getenv("SUPABASE_URL")+"/auth/v1/token?grant_type=password", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return accessToken, nil
}

func (repo *supabaseRepository) ValidateSession(ctx context.Context, token string) (models.User, error) {
	req, _ := http.NewRequest("GET", os.Getenv("SUPABASE_URL")+"/auth/v1/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", os.Getenv("SUPABASE_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return models.User{}, err
		}

		var userInfo models.User
		err = json.Unmarshal(body, &userInfo)
		if err != nil {
			return models.User{}, err
		}

		return userInfo, nil
	} else if resp.StatusCode == http.StatusUnauthorized {
		return models.User{}, fmt.Errorf("unauthorized")
	}

	log.Println("unexpected status code: ", resp.StatusCode)
	return models.User{}, fmt.Errorf("error with status code: %d", resp.StatusCode)
}
