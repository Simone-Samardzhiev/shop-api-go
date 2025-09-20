package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	*gin.Engine
}

func NewRouter(userHandler *UserHandler) *Router {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("min_bytes", validateMinBytesLength)
		_ = v.RegisterValidation("max_bytes", validateMaxBytesLength)
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.Use(errorMiddleware())
	{
		user := v1.Group("/users")
		{
			user.POST("/register", userHandler.Register)
		}

	}
	return &Router{r}
}

func (r *Router) Start(addr string) error {
	return r.Run(addr)
}
