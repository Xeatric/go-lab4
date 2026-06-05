package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"paving-tiles-api/internal/auth/middleware"
	"paving-tiles-api/internal/auth/oauth"
	"paving-tiles-api/internal/auth/service"
	"paving-tiles-api/internal/config"
	"paving-tiles-api/internal/dto"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
	yandexOAuth *oauth.YandexOAuth
	config      *config.Config
}

func NewAuthController(authService *service.AuthService, cfg *config.Config) *AuthController {
	yandexOAuth := oauth.NewYandexOAuth(
		cfg.YandexClientID,
		cfg.YandexClientSecret,
		cfg.YandexRedirectURL,
	)

	return &AuthController{
		authService: authService,
		yandexOAuth: yandexOAuth,
		config:      cfg,
	}
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанными email, паролем и именем
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Данные для регистрации"
// @Success      201 {object} map[string]interface{} "Успешная регистрация, токены в cookies"
// @Failure      400 {object} map[string]interface{} "Ошибка валидации или пользователь уже существует"
// @Router       /auth/register [post]

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	ip := ctx.ClientIP()

	response, err := c.authService.Register(&req, userAgent, ip)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Установка HttpOnly cookies
	ctx.SetCookie("access_token", response.AccessToken, 900, "/", "", false, true)
	ctx.SetCookie("refresh_token", response.RefreshToken, 604800, "/", "", false, true)

	ctx.JSON(http.StatusCreated, gin.H{
		"user":    response.User,
		"message": "registered successfully",
	})
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Вход в систему с получением JWT токенов в HttpOnly cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Учетные данные"
// @Success      200 {object} map[string]interface{} "Успешный вход, токены в cookies"
// @Failure      400 {object} map[string]interface{} "Ошибка валидации"
// @Failure      401 {object} map[string]interface{} "Неверный email или пароль"
// @Router       /auth/login [post]

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	ip := ctx.ClientIP()

	response, err := c.authService.Login(&req, userAgent, ip)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Установка HttpOnly cookies
	ctx.SetCookie("access_token", response.AccessToken, 900, "/", "", false, true)
	ctx.SetCookie("refresh_token", response.RefreshToken, 604800, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"user":    response.User,
		"message": "authenticated successfully",
	})
}

// Refresh godoc
// @Summary      Обновление токенов
// @Description  Использует refresh token из cookie или тела запроса для выдачи новой пары токенов
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest false "Refresh токен (опционально, если передан в cookie)"
// @Success      200 {object} map[string]interface{} "Новые токены установлены в cookies"
// @Failure      400 {object} map[string]interface{} "Отсутствует refresh token"
// @Failure      401 {object} map[string]interface{} "Невалидный или истекший refresh token"
// @Router       /auth/refresh [post]

func (c *AuthController) Refresh(ctx *gin.Context) {
	// Пробуем получить refresh token из cookie или из JSON
	var refreshToken string

	// Сначала проверяем cookie
	cookie, err := ctx.Cookie("refresh_token")
	if err == nil && cookie != "" {
		refreshToken = cookie
	} else {
		// Если нет в cookie, пробуем из JSON
		var req dto.RefreshRequest
		if err := ctx.ShouldBindJSON(&req); err == nil {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	ip := ctx.ClientIP()

	response, err := c.authService.Refresh(refreshToken, userAgent, ip)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем новые токены
	ctx.SetCookie("access_token", response.AccessToken, 900, "/", "", false, true)
	ctx.SetCookie("refresh_token", response.RefreshToken, 604800, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"user":    response.User,
		"message": "tokens refreshed successfully",
	})
}

// Logout godoc
// @Summary      Выход из текущей сессии
// @Description  Отзывает refresh token текущей сессии и очищает cookies
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{} "Успешный выход"
// @Failure      400 {object} map[string]interface{} "Нет активной сессии"
// @Router       /auth/logout [post]

func (c *AuthController) Logout(ctx *gin.Context) {
	// Получаем refresh token из cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no active session"})
		return
	}

	if err := c.authService.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Удаляем cookies
	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

// LogoutAll godoc
// @Summary      Выход из всех устройств
// @Description  Отзывает все refresh токены пользователя
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{} "Успешный выход со всех устройств"
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Router       /auth/logout-all [post]
func (c *AuthController) LogoutAll(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.authService.LogoutAll(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Удаляем cookies
	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully logged out from all devices"})
}

// Whoami godoc
// @Summary      Информация о текущем пользователе
// @Description  Возвращает данные аутентифицированного пользователя
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.WhoamiResponse
// @Failure      401 {object} map[string]interface{} "Не авторизован"
// @Router       /auth/whoami [get]

func (c *AuthController) Whoami(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}

	user, err := c.authService.Whoami(userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user":          user,
	})
}

// OAuthLogin godoc
// @Summary      Инициирует OAuth2 вход через Yandex
// @Description  Редиректит пользователя на страницу авторизации Яндекса
// @Tags         Auth
// @Param        provider path string true "Провайдер OAuth (yandex)"
// @Success      302
// @Failure      400 {object} map[string]interface{} "Неподдерживаемый провайдер"
// @Router       /auth/oauth/{provider} [get]

func (c *AuthController) OAuthLogin(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "yandex" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	// Генерируем state
	state := generateState()

	// Сохраняем в cookie с правильными параметрами
	ctx.SetCookie(
		"oauth_state",
		state,
		600,   // maxAge 600 секунд
		"/",   // path
		"",    // domain
		false, // secure (false для http)
		true,  // httpOnly
	)

	// Для отладки
	fmt.Printf("OAuthLogin: state=%s, cookie set\n", state)

	// Формируем URL авторизации
	authURL := fmt.Sprintf(
		"https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		c.config.YandexClientID,
		url.QueryEscape(c.config.YandexRedirectURL),
		state,
	)

	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OAuthCallback godoc
// @Summary      Callback OAuth2 от Яндекса
// @Description  Обрабатывает ответ от Яндекса, создает/находит пользователя и выдает токены
// @Tags         Auth
// @Param        provider path string true "Провайдер OAuth (yandex)"
// @Param        code query string true "Временный код от Яндекса"
// @Param        state query string true "CSRF protection state"
// @Success      200 {object} map[string]interface{} "Успешная аутентификация, токены в cookies"
// @Failure      400 {object} map[string]interface{} "Неверный state параметр"
// @Failure      500 {object} map[string]interface{} "Ошибка обмена кода или получения данных"
// @Router       /auth/oauth/{provider}/callback [get]

func (c *AuthController) OAuthCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "yandex" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	code := ctx.Query("code")
	state := ctx.Query("state")

	// Для отладки
	fmt.Printf("OAuthCallback: received state=%s, code=%s\n", state, code)

	// Получаем сохраненный state из cookie
	savedState, err := ctx.Cookie("oauth_state")
	fmt.Printf("OAuthCallback: savedState from cookie=%s, error=%v\n", savedState, err)

	if err != nil || state == "" || state != savedState {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid state parameter",
			"details": fmt.Sprintf("state=%s, savedState=%s", state, savedState),
		})
		return
	}

	// Удаляем cookie
	ctx.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 2. Обмениваем "code" на Access Token от Яндекса
	tokenResp, err := c.yandexOAuth.ExchangeCode(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	// 3. Получаем информацию о пользователе
	userInfo, err := c.yandexOAuth.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// 4. Ищем или создаем пользователя в нашей БД
	user, err := c.authService.FindOrCreateUser(userInfo.Email, userInfo.RealName, "yandex", userInfo.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user"})
		return
	}

	// 5. Генерируем нашу собственную пару JWT токенов
	userAgent := ctx.GetHeader("User-Agent")
	ip := ctx.ClientIP()
	tokens, err := c.authService.GenerateTokensForUser(user, userAgent, ip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// 6. Устанавливаем наши токены в HttpOnly cookies
	ctx.SetCookie("access_token", tokens.AccessToken, 900, "/", "", false, true)
	ctx.SetCookie("refresh_token", tokens.RefreshToken, 604800, "/", "", false, true)

	// 7. Отправляем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "authenticated successfully",
		"user":    tokens.User,
	})
}

// generateState - генерация случайной строки для CSRF защиты
func generateState() string {
	return "random_state_string_" + time.Now().Format("20060102150405")
}
