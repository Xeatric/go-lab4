# Лабораторная работа №4: OpenAPI (Swagger) документация для Paving Tiles API

##  Описание работы

В данной лабораторной работе реализована автоматическая генерация OpenAPI (Swagger) документации для REST API сервиса управления каталогом тротуарной плитки. Документация создается с использованием подхода Code-First, где спецификация API генерируется непосредственно из аннотаций в коде.

##  Цели работы

- Изучение спецификации OpenAPI как стандарта описания веб-сервисов
- Освоение принципов автоматической генерации документации на основе кода (Code-First подход)
- Настройка интерактивной документации Swagger UI для тестирования API
- Реализация разделения конфигурации документации для сред разработки и production
- Настройка схем безопасности (JWT, Cookies) в интерфейсе документации

Установка и запуск

1. **Клонируйте репозиторий:**

```bash
   git clone https://github.com/Xeatric/go-lab4.git
   cd go-lab4
```
Отредактируйте .env

Запустите приложение:

```bash
docker-compose up --build
```
## Доступ к документации

Откройте в браузере: http://localhost:4200/swagger/index.html


## Auth (Аутентификация)
    Метод	    Путь	                Описание	    
    POST	/auth/register	            Регистрация нового пользователя	
    POST	/auth/login	                Вход в систему	
    POST	/auth/refresh	            Обновление JWT токенов	
    GET	    /api/v1/auth/whoami	        Информация о текущем пользователе	
    POST	/api/v1/auth/logout 	    Выход из текущей сессии	
    POST	/api/v1/auth/logout-all	    Выход из всех устройств	
## Tiles (Плитки)
      Метод	    Путь	                   Описание	
      GET	        /api/v1/tiles	           Получить список плиток (с пагинацией)	
      POST	    /api/v1/tiles	           Создать новую плитку	
      GET	        /api/v1/tiles/{id}         Получить плитку по ID	
      PUT	        /api/v1/tiles/{id}	       Полное обновление плитки	
      PATCH	    /api/v1/tiles/{id}	       Частичное обновление плитки	
      DELETE	    /api/v1/tiles/{id}	       Удаление плитки (soft delete)	
 
 
 ## Аутентификация
API использует JWT токены для аутентификации. Токены хранятся в HttpOnly cookies после успешного входа.

 ## Получение токена
Выполните запрос POST /auth/login с вашими учетными данными
Токены автоматически устанавливаются в cookies

## Регистрация пользователя
Запрос:

```json
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "Иван Петров"
}
```
Ответ (201 Created):

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "Иван Петров",
    "role": "user",
    "created_at": "2024-01-15T12:30:00Z"
  },
  "message": "registered successfully"
}
```

## Создание плитки
Запрос:

```json
POST /api/v1/tiles
Authorization: Bearer your_token
Content-Type: application/json

{
  "name": "Брусчатка Классик",
  "shape": "rectangle",
  "color": "серый",
  "size": "200x100x60",
  "material": "бетон",
  "price_per_m2": 850.50,
  "stock": 1500,
  "description": "Классическая брусчатка для пешеходных зон",
  "image_url": "https://example.com/brick.jpg"
}
```

Ответ (201 Created):

```json
{
  "id": 1,
  "name": "Брусчатка Классик",
  "shape": "rectangle",
  "color": "серый",
  "size": "200x100x60",
  "material": "бетон",
  "price_per_m2": 850.5,
  "stock": 1500,
  "description": "Классическая брусчатка для пешеходных зон",
  "image_url": "https://example.com/brick.jpg",
  "created_at": "2024-01-15T12:30:00Z",
  "updated_at": "2024-01-15T12:30:00Z"
}
```

## Получение списка плиток с пагинацией

Запрос:

http
GET /api/v1/tiles?page=1&limit=10
Authorization: Bearer your_token
Ответ (200 OK):

```json
{
  "data": [
    {
      "id": 1,
      "name": "Брусчатка Классик",
      "price_per_m2": 850.5,
      "stock": 1500
    }
  ],
  "meta": {
    "total": 42,
    "page": 1,
    "limit": 10,
    "total_pages": 5
  }
}
```
