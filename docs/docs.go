package docs

import (
	"github.com/swaggo/swag"
)

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "Paving Tiles API",
        "description": "API для управления каталогом тротуарной плитки с аутентификацией",
        "version": "1.0",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@paving-tiles.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        }
    },
    "host": "localhost:4200",
    "basePath": "/",
    "schemes": ["http"],
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Введите 'Bearer {ваш JWT токен}' для авторизации"
        }
    },
    "paths": {
        "/auth/register": {
            "post": {
                "tags": ["Auth"],
                "summary": "Регистрация нового пользователя",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email", "password", "name"],
                            "properties": {
                                "email": {"type": "string", "example": "user@example.com"},
                                "password": {"type": "string", "example": "securePassword123"},
                                "name": {"type": "string", "example": "Иван Петров"}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {"description": "Успешная регистрация"},
                    "400": {"description": "Ошибка валидации"}
                }
            }
        },
        "/auth/login": {
            "post": {
                "tags": ["Auth"],
                "summary": "Аутентификация пользователя",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["email", "password"],
                            "properties": {
                                "email": {"type": "string", "example": "user@example.com"},
                                "password": {"type": "string", "example": "securePassword123"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Успешный вход"},
                    "401": {"description": "Неверные данные"}
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "tags": ["Auth"],
                "summary": "Обновление токенов",
                "responses": {
                    "200": {"description": "Токены обновлены"}
                }
            }
        },
        "/api/v1/auth/logout": {
            "post": {
                "tags": ["Auth"],
                "summary": "Выход из текущей сессии",
                "security": [{"BearerAuth": []}],
                "responses": {
                    "200": {"description": "Успешный выход"}
                }
            }
        },
        "/api/v1/auth/logout-all": {
            "post": {
                "tags": ["Auth"],
                "summary": "Выход из всех устройств",
                "security": [{"BearerAuth": []}],
                "responses": {
                    "200": {"description": "Успешный выход со всех устройств"}
                }
            }
        },
        "/api/v1/auth/whoami": {
            "get": {
                "tags": ["Auth"],
                "summary": "Информация о текущем пользователе",
                "security": [{"BearerAuth": []}],
                "responses": {
                    "200": {"description": "Данные пользователя"}
                }
            }
        },
        "/api/v1/tiles": {
            "get": {
                "tags": ["Tiles"],
                "summary": "Получить список плиток",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {"in": "query", "name": "page", "type": "integer", "default": 1, "description": "Номер страницы"},
                    {"in": "query", "name": "limit", "type": "integer", "default": 10, "description": "Количество на странице"}
                ],
                "responses": {
                    "200": {"description": "Список плиток"},
                    "401": {"description": "Не авторизован"}
                }
            },
            "post": {
                "tags": ["Tiles"],
                "summary": "Создать плитку",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["name", "shape", "color", "size", "price_per_m2"],
                            "properties": {
                                "name": {"type": "string", "example": "Брусчатка Классик"},
                                "shape": {"type": "string", "enum": ["square", "rectangle", "hexagon", "circle"], "example": "rectangle"},
                                "color": {"type": "string", "example": "серый"},
                                "size": {"type": "string", "example": "200x100x60"},
                                "material": {"type": "string", "example": "бетон"},
                                "price_per_m2": {"type": "number", "example": 850.50},
                                "stock": {"type": "integer", "example": 1500},
                                "description": {"type": "string", "example": "Классическая брусчатка"},
                                "image_url": {"type": "string", "example": "https://example.com/photo.jpg"}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {"description": "Плитка создана"},
                    "400": {"description": "Ошибка валидации"},
                    "401": {"description": "Не авторизован"}
                }
            }
        },
        "/api/v1/tiles/{id}": {
            "get": {
                "tags": ["Tiles"],
                "summary": "Получить плитку по ID",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "in": "path",
                        "name": "id",
                        "required": true,
                        "type": "integer",
                        "description": "ID плитки",
                        "example": 1
                    }
                ],
                "responses": {
                    "200": {"description": "Данные плитки"},
                    "404": {"description": "Плитка не найдена"}
                }
            },
            "put": {
                "tags": ["Tiles"],
                "summary": "Полное обновление плитки",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "path",
                        "name": "id",
                        "required": true,
                        "type": "integer",
                        "description": "ID плитки",
                        "example": 1
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["name", "shape", "color", "size", "price_per_m2"],
                            "properties": {
                                "name": {"type": "string", "example": "Обновленная плитка"},
                                "shape": {"type": "string", "enum": ["square", "rectangle", "hexagon", "circle"], "example": "square"},
                                "color": {"type": "string", "example": "темно-серый"},
                                "size": {"type": "string", "example": "300x300x60"},
                                "material": {"type": "string", "example": "гранит"},
                                "price_per_m2": {"type": "number", "example": 2500.00},
                                "stock": {"type": "integer", "example": 500},
                                "description": {"type": "string", "example": "Премиальная плитка"},
                                "image_url": {"type": "string", "example": "https://example.com/photo.jpg"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Плитка обновлена"},
                    "400": {"description": "Ошибка валидации"},
                    "404": {"description": "Плитка не найдена"}
                }
            },
            "patch": {
                "tags": ["Tiles"],
                "summary": "Частичное обновление плитки",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "path",
                        "name": "id",
                        "required": true,
                        "type": "integer",
                        "description": "ID плитки",
                        "example": 1
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string", "example": "Новое название"},
                                "price_per_m2": {"type": "number", "example": 2750.00},
                                "stock": {"type": "integer", "example": 450}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "Плитка обновлена"},
                    "400": {"description": "Ошибка валидации"},
                    "404": {"description": "Плитка не найдена"}
                }
            },
            "delete": {
                "tags": ["Tiles"],
                "summary": "Удалить плитку (soft delete)",
                "security": [{"BearerAuth": []}],
                "parameters": [
                    {
                        "in": "path",
                        "name": "id",
                        "required": true,
                        "type": "integer",
                        "description": "ID плитки",
                        "example": 1
                    }
                ],
                "responses": {
                    "204": {"description": "Плитка удалена"},
                    "404": {"description": "Плитка не найдена"}
                }
            }
        }
    },
    "tags": [
        {"name": "Auth", "description": "Операции аутентификации и управления пользователями"},
        {"name": "Tiles", "description": "CRUD операции с плитками"}
    ]
}`

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:4200",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Paving Tiles API",
	Description:      "API для управления каталогом тротуарной плитки с аутентификацией",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
