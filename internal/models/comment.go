package models

// Comment представляет комментарий к новости
type Comment struct {
	ID        int    `json:"id"`
	NewsID    int    `json:"news_id"`
	ParentID  *int   `json:"parent_id,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// CommentRequest представляет запрос на создание нового комментария
type CommentRequest struct {
	NewsID   int    `json:"news_id"`
	ParentID *int   `json:"parent_id,omitempty"`
	Content  string `json:"content"`
}
