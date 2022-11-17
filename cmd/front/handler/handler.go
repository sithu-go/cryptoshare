package handler

import (
	"cryptoshare/ds"
	"cryptoshare/middleware"
	"cryptoshare/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	R    *gin.Engine
	repo *repository.Repository
}

type HConfig struct {
	R  *gin.Engine
	DS *ds.DataSource
}

func NewHandler(c *HConfig) *Handler {
	repo := repository.NewRepository(c.DS)
	return &Handler{
		R:    c.R,
		repo: repo,
	}
}

func (h *Handler) Register() {
	h.R.Use(middleware.Cors())

	// auth routes
	authHandler := newAuthHandler(h)
	authHandler.register()
}
