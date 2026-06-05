@"
# Paving Tiles API

API для управления каталогом тротуарной плитки с аутентификацией.

## Запуск проекта

### Через Docker Compose

`docker-compose up --build`
## Документация API

После запуска документация доступна по адресу:
http://localhost:4200/swagger/index.html

## Endpoints

 Auth
- POST /auth/register - регистрация
- POST /auth/login - вход
- POST /auth/refresh - обновление токенов
- GET /api/v1/auth/whoami - информация о пользователе
- POST /api/v1/auth/logout - выход
- POST /api/v1/auth/logout-all - выход со всех устройств

 Tiles (требуют авторизации)
- GET /api/v1/tiles - список плиток
- POST /api/v1/tiles - создание
- GET /api/v1/tiles/{id} - получение по ID
- PUT /api/v1/tiles/{id} - полное обновление
- PATCH /api/v1/tiles/{id} - частичное обновление
- DELETE /api/v1/tiles/{id} - удаление
