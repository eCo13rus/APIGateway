# Микросервисный новостной агрегатор

Система микросервисов на **Go** для сбора, хранения и отображения новостей с возможностью комментирования.

## 📌 Структура проекта
Система состоит из четырех микросервисов:

- **API Gateway** (`порт 8080`) – маршрутизатор запросов.
- **News Aggregator** (`порт 8081`) – сервис для сбора и хранения новостей.
- **Comments Service** (`порт 8082`) – сервис комментариев.
- **Censorship Service** (`порт 8083`) – сервис цензурирования контента.

## 📋 Требования

- **Go** `1.22` или выше
- **PostgreSQL** `14` или выше

## 🚀 Установка и запуск

### 1️⃣ Клонирование репозиториев:

```sh
git clone https://github.com/eCo13rus/news-aggregator
git clone https://github.com/eCo13rus/comments_service
git clone https://github.com/eCo13rus/censorship_service
git clone https://github.com/eCo13rus/APIGateway


### cd news-aggregator
## Создание/Настройка/подключение БД - файл для миграции ~/news_aggregator/migrations/init_DB.sql

### cd comments_service
## Создание/Настройка/подключение БД - файл для миграции ~/comments_service/migrations/001_init_schema.sql


### cd news-aggregator
### go run cmd/main.go

### cd censorship_service
### go run cmd/main.go

### cd comments_service
### go run cmd/main.go

### cd APIGateway
### go run cmd/main.go

### Добавьте коллецию в Postman запросов/ответов - ~/news_Aggregator_collection






