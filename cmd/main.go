package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/eCo13rus/APIGateway/internal/api"
	"github.com/eCo13rus/APIGateway/internal/models"
)

func main() {
	config, err := loadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	handler := api.NewHandler(
		config.Services.NewsURL,
		config.Services.CommentURL,
		config.Services.CensorURL,
	)

	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	server := api.NewServer(handler, addr)

	if err := server.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func loadConfig(path string) (*models.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла конфигурации: %v", err)
	}
	defer file.Close()

	var config models.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("ошибка декодирования конфигурации: %v", err)
	}

	return &config, nil
}
