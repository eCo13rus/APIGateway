package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CensorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCensorClient(baseURL string) *CensorClient {
	return &CensorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CheckContent проверяет текст на цензуру
func (c *CensorClient) CheckContent(content string, requestID string) error {
	reqURL := fmt.Sprintf("%s/api/censor", c.baseURL)

	params := url.Values{}
	if requestID != "" {
		params.Add("request_id", requestID)
	}

	body := struct {
		Content string `json:"content"`
	}{
		Content: content,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("ошибка кодирования запроса: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL+"?"+params.Encode(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("текст не прошел проверку цензуры")
		}

		return fmt.Errorf("текст не прошел проверку цензуры: %s", errorResp.Message)
	}

	return nil
}
