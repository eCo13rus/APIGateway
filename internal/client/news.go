package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eCo13rus/APIGateway/internal/models"
)

// NewsClient представляет клиент для взаимодействия с сервисом новостей
type NewsClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewNewsClient создает новый клиент для сервиса новостей
func NewNewsClient(baseURL string) *NewsClient {
	return &NewsClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetNews получает список новостей с возможностью поиска и пагинации
func (c *NewsClient) GetNews(count int, page int, searchQuery, requestID string) (*models.NewsResponse, error) {
	reqURL := fmt.Sprintf("%s/api/news/%d", c.baseURL, count)

	params := url.Values{}
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}
	if searchQuery != "" {
		params.Add("s", searchQuery)
	}
	if requestID != "" {
		params.Add("request_id", requestID)
	}

	log.Printf("[%s] Отправка запроса к сервису новостей: %s", requestID, reqURL+"?"+params.Encode())

	req, err := http.NewRequest(http.MethodGet, reqURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[%s] Ответ от сервиса новостей: %d", requestID, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервис новостей вернул ошибку: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	log.Printf("[%s] Тело ответа: %s", requestID, string(bodyBytes))

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var newsResponse models.NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResponse); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &newsResponse, nil
}

// GetNewsDetails получает детальную информацию о новости по ID
func (c *NewsClient) GetNewsDetails(newsID int, requestID string) (*models.NewsFullDetailed, error) {
	reqURL := fmt.Sprintf("%s/api/news/detail/%d", c.baseURL, newsID)

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("новость с ID %d не найдена", newsID)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервис новостей вернул ошибку: %d", resp.StatusCode)
	}

	var newsDetail models.NewsFullDetailed
	if err := json.NewDecoder(resp.Body).Decode(&newsDetail); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &newsDetail, nil
}
