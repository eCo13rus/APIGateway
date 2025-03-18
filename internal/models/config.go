package models

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Services ServicesConfig `json:"services"`
}

// ServerConfig представляет настройки HTTP-сервера
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// ServicesConfig содержит адреса микросервисов
type ServicesConfig struct {
	NewsURL    string `json:"news_service_url"`
	CommentURL string `json:"comment_service_url"`
	CensorURL  string `json:"censor_service_url"`
}
