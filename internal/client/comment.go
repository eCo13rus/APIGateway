package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/eCo13rus/APIGateway/internal/models"
)

type CommentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCommentClient(baseURL string) *CommentClient {
	return &CommentClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetComments получает список комментариев для новости
func (c *CommentClient) GetComments(newsID int, requestID string) ([]models.Comment, error) {
	reqURL := fmt.Sprintf("%s/api/comments/news/%d", c.baseURL, newsID)

	params := url.Values{}
	if requestID != "" {
		params.Add("request_id", requestID)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервис комментариев вернул ошибку: %d", resp.StatusCode)
	}

	var response struct {
		Comments []models.Comment `json:"comments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return response.Comments, nil
}

func (c *CommentClient) AddComment(comment *models.CommentRequest, requestID string) error {
	reqURL := fmt.Sprintf("%s/api/comments", c.baseURL)

	params := url.Values{}
	if requestID != "" {
		params.Add("request_id", requestID)
	}

	body, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("ошибка кодирования запроса: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL+"?"+params.Encode(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервис комментариев вернул ошибку: %d", resp.StatusCode)
	}

	return nil
}
