package models

// NewsShortDetailed представляет краткую информацию о новости для списков
type NewsShortDetailed struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content,omitempty"`
	PubTime   int64  `json:"pub_time"`
	Link      string `json:"link"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewsFullDetailed представляет детальную информацию о новости
type NewsFullDetailed struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	PubTime   int64     `json:"pub_time"`
	Link      string    `json:"link"`
	Comments  []Comment `json:"comments,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// NewsResponse представляет ответ со списком новостей и пагинацией
type NewsResponse struct {
	News       []NewsShortDetailed `json:"news"`
	Pagination Pagination          `json:"pagination"`
}

// Pagination представляет информацию о пагинации
type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
}
