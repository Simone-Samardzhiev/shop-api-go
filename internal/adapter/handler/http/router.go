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
	authHandler *AuthHandler,
) *Router {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("min_bytes", validateMinBytesLength)
		_ = v.RegisterValidation("max_bytes", validateMaxBytesLength)
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
			user.Use(jwtMiddleware)
			user.GET("/pagination-by-offset", userHandler.GetUsersByOffsetPagination)
			user.GET("/pagination-by-time", userHandler.GetUsersByTimePagination)
			user.GET("/search/by-username", userHandler.SearchUserByUsername)
			user.GET("/search/by-email", userHandler.SearchUserByEmail)
		}
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/refresh-session", jwtMiddleware, authHandler.RefreshSession)
		}
	}
	return &Router{r, appConfig}
}

func (r *Router) Start() error {
	return r.Run(r.appConfig.Port)
}
