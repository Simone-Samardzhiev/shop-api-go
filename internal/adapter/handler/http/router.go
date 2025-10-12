package http

import (
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/core/port"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	*gin.Engine
	appConfig *config.AppConfig
}

func NewRouter(
	appConfig *config.AppConfig,
	tokenGenerator port.TokenGenerator,
	userHandler *UserHandler,
	adminHandler *AdminHandler,
	authHandler *AuthHandler,
) *Router {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("min_bytes", validateMinBytesLength)
		_ = v.RegisterValidation("max_bytes", validateMaxBytesLength)
		_ = v.RegisterValidation("user_role", validateUserRole)
	}

	switch appConfig.Environment {
	case config.Production:
		gin.SetMode(gin.ReleaseMode)
	case config.Development:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(zapLogger())
	jwtMiddleware := newJwtMiddleware(tokenGenerator, "token")

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/register", userHandler.Register)
			user.PATCH("/me/change-username", userHandler.ChangeUsername)
		}

		admin := v1.Group("/admin")
		admin.Use(jwtMiddleware)
		{
			admin.GET("/pagination-by-offset", adminHandler.GetUsersByOffsetPagination)
			admin.GET("/pagination-by-time", adminHandler.GetUsersByTimePagination)
			admin.GET("/search/by-username", adminHandler.SearchUserByUsername)
			admin.GET("/search/by-email", adminHandler.SearchUserByEmail)
			admin.GET("/by-id", adminHandler.GetUserById)
			admin.PATCH("/update/username", adminHandler.UpdateUsername)
			admin.PATCH("/update/email", adminHandler.UpdateEmail)
			admin.PATCH("/update/password", adminHandler.UpdatePassword)
			admin.PATCH("/update/role", adminHandler.UpdateRole)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh-session", jwtMiddleware, authHandler.RefreshSession)
		}
	}
	return &Router{r, appConfig}
}

func (r *Router) Start() error {
	return r.Run(r.appConfig.Port)
}
