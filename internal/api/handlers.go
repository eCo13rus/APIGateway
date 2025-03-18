package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/eCo13rus/APIGateway/internal/client"
	"github.com/eCo13rus/APIGateway/internal/models"
)

type Handler struct {
	newsClient    *client.NewsClient
	commentClient *client.CommentClient
	censorClient  *client.CensorClient
}

func NewHandler(newsURL, commentURL, censorURL string) *Handler {
	return &Handler{
		newsClient:    client.NewNewsClient(newsURL),
		commentClient: client.NewCommentClient(commentURL),
		censorClient:  client.NewCensorClient(censorURL),
	}
}

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("s")
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	requestID := r.Context().Value(RequestIDKey).(string)

	newsResponse, err := h.newsClient.GetNews(10, page, searchQuery, requestID)
	if err != nil {
		log.Printf("[%s] Ошибка получения новостей: %v", requestID, err)
		http.Error(w, "Ошибка получения новостей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newsResponse); err != nil {
		log.Printf("[%s] Ошибка кодирования ответа: %v", requestID, err)
	}
}

// GetNewsDetails обрабатывает запрос на получение детальной информации о новости
func (h *Handler) GetNewsDetails(w http.ResponseWriter, r *http.Request) {
	// Получаем ID новости из URL
	vars := mux.Vars(r)
	newsIDStr := vars["id"]
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value(RequestIDKey).(string)

	type result struct {
		data interface{}
		err  error
	}

	resultChan := make(chan result, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		newsDetails, err := h.newsClient.GetNewsDetails(newsID, requestID)
		resultChan <- result{data: newsDetails, err: err}
	}()

	go func() {
		defer wg.Done()
		comments, err := h.commentClient.GetComments(newsID, requestID)
		resultChan <- result{data: comments, err: err}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var newsDetails *models.NewsFullDetailed
	var comments []models.Comment
	var resultErr error

	for res := range resultChan {
		if res.err != nil {
			resultErr = res.err
			continue
		}

		// Проверяем тип результата и записываем в соответствующую переменную
		switch data := res.data.(type) {
		case *models.NewsFullDetailed:
			newsDetails = data
		case []models.Comment:
			comments = data
		}
	}

	if resultErr != nil {
		log.Printf("[%s] Ошибка получения данных: %v", requestID, resultErr)

		if strings.Contains(resultErr.Error(), "не найдена") {
			http.Error(w, resultErr.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		}
		return
	}

	if newsDetails == nil {
		http.Error(w, fmt.Sprintf("Новость с ID %d не найдена", newsID), http.StatusNotFound)
		return
	}

	newsDetails.Comments = comments

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newsDetails); err != nil {
		log.Printf("[%s] Ошибка кодирования ответа: %v", requestID, err)
	}
}

// AddComment обрабатывает запрос на добавление комментария
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsIDStr := vars["id"]
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	requestID := r.Context().Value(RequestIDKey).(string)

	var req struct {
		ParentID *int   `json:"parent_id,omitempty"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.censorClient.CheckContent(req.Content, requestID); err != nil {
		log.Printf("[%s] Ошибка проверки цензуры: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentReq := &models.CommentRequest{
		NewsID:   newsID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := h.commentClient.AddComment(commentReq, requestID); err != nil {
		log.Printf("[%s] Ошибка добавления комментария: %v", requestID, err)
		http.Error(w, "Ошибка добавления комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"status":  "success",
		"message": "Комментарий успешно добавлен",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNewsComments обрабатывает запрос на получение комментариев к новости
func (h *Handler) GetNewsComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsIDStr := vars["id"]
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value(RequestIDKey).(string)
	log.Printf("[%s] Запрос комментариев для новости ID: %d", requestID, newsID)

	comments, err := h.commentClient.GetComments(newsID, requestID)
	if err != nil {
		log.Printf("[%s] Ошибка получения комментариев: %v", requestID, err)
		http.Error(w, "Ошибка получения комментариев", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"comments": comments,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Ошибка кодирования ответа: %v", requestID, err)
	}
}
