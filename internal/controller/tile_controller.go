package controller

import (
	"net/http"
	"strconv"

	"paving-tiles-api/internal/auth/middleware"
	"paving-tiles-api/internal/dto"
	"paving-tiles-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TileController struct {
	service service.TileService
}

func NewTileController(service service.TileService) *TileController {
	return &TileController{service: service}
}

// GetTiles godoc
// @Summary      Получить список плиток текущего пользователя
// @Description  Возвращает пагинированный список плиток, принадлежащих авторизованному пользователю
// @Tags         Tiles
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество элементов на странице" default(10) maximum(100)
// @Success      200 {object} dto.PaginationResponse
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router       /tiles [get]

func (c *TileController) GetTiles(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetTiles(userID, paginationReq.Page, paginationReq.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTileByID godoc
// @Summary      Получить плитку по ID
// @Description  Возвращает одну плитку, если она принадлежит текущему пользователю
// @Tags         Tiles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID плитки"
// @Success      200 {object} dto.TileResponse
// @Failure      400 {object} map[string]interface{} "Неверный ID"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      403 {object} map[string]interface{} "Доступ запрещен (не своя плитка)"
// @Failure      404 {object} map[string]interface{} "Плитка не найдена"
// @Router       /tiles/{id} [get]

func (c *TileController) GetTileByID(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tile, err := c.service.GetTileByID(userID, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// CreateTile godoc
// @Summary      Создать новую плитку
// @Description  Добавляет новую плитку в каталог текущего пользователя
// @Tags         Tiles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTileRequest true "Данные новой плитки"
// @Success      201 {object} dto.TileResponse
// @Failure      400 {object} map[string]interface{} "Ошибка валидации"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router       /tiles [post]

func (c *TileController) CreateTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.CreateTile(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tile)
}

// UpdateTile godoc
// @Summary      Полное обновление плитки
// @Description  Заменяет все поля существующей плитки
// @Tags         Tiles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID плитки"
// @Param        request body dto.UpdateTileRequest true "Новые данные плитки"
// @Success      200 {object} dto.TileResponse
// @Failure      400 {object} map[string]interface{} "Ошибка валидации или ID"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      403 {object} map[string]interface{} "Доступ запрещен"
// @Failure      404 {object} map[string]interface{} "Плитка не найдена"
// @Router       /tiles/{id} [put]

func (c *TileController) UpdateTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.UpdateTile(userID, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// PatchTile godoc
// @Summary      Частичное обновление плитки
// @Description  Обновляет только указанные поля плитки
// @Tags         Tiles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "ID плитки"
// @Param        request body dto.PatchTileRequest true "Поля для обновления"
// @Success      200 {object} dto.TileResponse
// @Failure      400 {object} map[string]interface{} "Ошибка валидации или ID"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      403 {object} map[string]interface{} "Доступ запрещен"
// @Failure      404 {object} map[string]interface{} "Плитка не найдена"
// @Router       /tiles/{id} [patch]

func (c *TileController) PatchTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.PatchTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.PatchTile(userID, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// DeleteTile godoc
// @Summary      Удалить плитку (soft delete)
// @Description  Помечает плитку как удаленную (не удаляет физически)
// @Tags         Tiles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID плитки"
// @Success      204 "Нет содержимого (успешное удаление)"
// @Failure      400 {object} map[string]interface{} "Неверный ID"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Failure      403 {object} map[string]interface{} "Доступ запрещен"
// @Failure      404 {object} map[string]interface{} "Плитка не найдена"
// @Router       /tiles/{id} [delete]

func (c *TileController) DeleteTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.service.DeleteTile(userID, uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
