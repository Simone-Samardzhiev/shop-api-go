package http

import (
	"context"
	"net/http"
	"shop-api-go/internal/adapter/config"
	"shop-api-go/internal/adapter/handler/http/middleware"
	"shop-api-go/internal/core/port"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	server *http.Server
}

func NewRouter(
	appConfig *config.AppConfig,
	tokenGenerator port.TokenGenerator,
	userHandler *UserHandler,
	adminHandler *AdminHandler,
	authHandler *AuthHandler,
) (*Router, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("password", validatePassword); err != nil {
			return nil, err
		}
		if err := v.RegisterValidation("min_bytes", validateMinBytesLength); err != nil {
			return nil, err
		}
		if err := v.RegisterValidation("max_bytes", validateMaxBytesLength); err != nil {
			return nil, err
		}
		if err := v.RegisterValidation("user_role", validateUserRole); err != nil {
			return nil, err
		}
	}

	switch appConfig.Environment {
	case config.Production:
		gin.SetMode(gin.ReleaseMode)
	case config.Development:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger())
	jwtMiddleware := middleware.JWTMiddleware(tokenGenerator, "token")

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/register", userHandler.Register)
			user.PATCH("/me/update", userHandler.UpdateAccount)
		}

		admin := v1.Group("/admin")
		admin.Use(jwtMiddleware)
		{
			adminUser := admin.Group("/users")
			{
				adminUser.GET("", adminHandler.GetUsers)
				adminUser.PATCH("/update/:id", adminHandler.UpdateUser)
			}
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh-session", jwtMiddleware, authHandler.RefreshSession)
		}

	}
	return &Router{&http.Server{
		Addr:    appConfig.Port,
		Handler: r,
	}}, nil
}

func (r *Router) Start() error {
	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
