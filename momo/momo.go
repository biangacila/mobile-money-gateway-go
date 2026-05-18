package momo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Client interface {
	CreateAPIUser(callbackHost string) (string, error)
	GetAPIUser(apiUserRef string) (string, error)
	CreateAPIKey(apiUserID string) (string, error)
	GenerateCollectionToken() (string, error)
	RequestToPay(token string, payload RequestToPayPayload) (string, error)
}

type ClientImpl struct {
	SubscriptionKey string
	APIUser         string
	APIKey          string
	HttpClient      *http.Client

	baseURL           string
	targetEnvironment string
}

func NewClient(subscriptionKey string, APIUser string, APIKey string, baseURL string, targetEnvironment string) *ClientImpl {
	return &ClientImpl{
		SubscriptionKey:   subscriptionKey,
		APIUser:           APIUser,
		APIKey:            APIKey,
		HttpClient:        &http.Client{},
		baseURL:           baseURL,
		targetEnvironment: targetEnvironment,
	}
}

// 1. Create API User
func (m *ClientImpl) CreateAPIUser(callbackHost string) (string, error) {
	apiUserID := uuid.New().String()

	payload := map[string]string{
		"providerCallbackHost": callbackHost,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/v1_0/apiuser",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Reference-Id", apiUserID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", m.SubscriptionKey)

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create api user failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return apiUserID, nil
}

// 2. Create API Key
func (m *ClientImpl) CreateAPIKey(apiUserID string) (string, error) {
	req, err := http.NewRequest(
		"POST",
		baseURL+"/v1_0/apiuser/"+apiUserID+"/apikey",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", m.SubscriptionKey)

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create api key failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result struct {
		APIKey string `json:"apiKey"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.APIKey, nil
}

// 3. Generate Collection Access Token
func (m *ClientImpl) GenerateCollectionToken() (string, error) {
	authString := m.APIUser + ":" + m.APIKey
	authBase64 := base64.StdEncoding.EncodeToString([]byte(authString))

	req, err := http.NewRequest(
		"POST",
		baseURL+"/collection/token/",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+authBase64)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.SubscriptionKey)

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (m *ClientImpl) RequestToPay(token string, payload RequestToPayPayload) (string, error) {
	referenceID := uuid.New().String()

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/collection/v1_0/requesttopay",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("X-Target-Environment", targetEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.SubscriptionKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("request to pay failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return referenceID, nil
}

// 5. Check Payment Status
func (m *ClientImpl) GetRequestToPayStatus(token, referenceID string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		baseURL+"/collection/v1_0/requesttopay/"+referenceID,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Target-Environment", targetEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.SubscriptionKey)

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status check failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
