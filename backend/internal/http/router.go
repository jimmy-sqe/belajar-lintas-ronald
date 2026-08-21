package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler"
)

// Handlers groups HTTP handlers wired by app.NewServer.
type Handlers struct {
	// boilerplate:axis=auth option=jwt START
	Auth *handler.AuthHandler
	// boilerplate:axis=auth option=jwt END
	// boilerplate:axis=sample-app option=todo-app START
	Todo *handler.TodoHandler
	// boilerplate:axis=sample-app option=todo-app END
}

// MiddlewareStack groups middleware applied to /v1 routes (non-public).
type MiddlewareStack struct {
	Auth      echo.MiddlewareFunc
	Recover   echo.MiddlewareFunc
	RequestID echo.MiddlewareFunc
}

// RegisterRoutes attaches handlers and middleware to the Echo instance.
func RegisterRoutes(e *echo.Echo, h *Handlers, mw *MiddlewareStack) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// boilerplate:axis=auth option=jwt START
	// Public auth routes (no JWT middleware)
	if h.Auth != nil {
		authGroup := e.Group("/v1/auth")
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/renew", h.Auth.Renew)
		authGroup.POST("/logout", h.Auth.Logout)
	}
	// boilerplate:axis=auth option=jwt END

	v1 := e.Group("/v1")
	if mw.RequestID != nil {
		v1.Use(mw.RequestID)
	}
	if mw.Recover != nil {
		v1.Use(mw.Recover)
	}
	if mw.Auth != nil {
		v1.Use(mw.Auth)
	}

	// boilerplate:axis=sample-app option=todo-app START
	if h.Todo != nil {
		todos := v1.Group("/todos")
		todos.POST("", h.Todo.Create)
		todos.GET("", h.Todo.List)
		todos.GET("/:id", h.Todo.GetByID)
		todos.PUT("/:id", h.Todo.Update)
		todos.DELETE("/:id", h.Todo.Delete)
	}
	// boilerplate:axis=sample-app option=todo-app END
}
