package http

import (
	"net/http"

	"perpustakaan/internal/delivery/http/middleware"
	"perpustakaan/internal/domain"
	"perpustakaan/pkg/jwt"
	"perpustakaan/pkg/response"

	"github.com/labstack/echo/v4"
)

// UserHandler represent the httphandler for user
type UserHandler struct {
	UUsecase domain.UserUsecase
}

// NewUserHandler will initialize the user endpoints
func NewUserHandler(e *echo.Echo, us domain.UserUsecase, tokenMaker jwt.TokenMaker) {
	handler := &UserHandler{
		UUsecase: us,
	}

	// Public routes
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	// Protected routes
	authGrp := e.Group("/api")
	authGrp.Use(middleware.AuthMiddleware(tokenMaker))
	authGrp.GET("/profile", handler.GetProfile)
}

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var req authRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusUnprocessableEntity, err.Error())
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UUsecase.Register(c.Request().Context(), &user)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusCreated, "user registered successfully", nil)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusUnprocessableEntity, err.Error())
	}

	token, err := h.UUsecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error())
	}

	return response.Success(c, http.StatusOK, "login successful", map[string]string{
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.UUsecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, http.StatusOK, "profile retrieved successfully", user)
}
