package handler

import (
	"propulse/repository"
	"propulse/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServices(db *pgxpool.Pool, router *gin.Engine) {
	PropostaRepo := repository.NewPropostaRepository(db)
	PropostaService := service.NewPropostaService(PropostaRepo)
	PropostaHandler := NewPropostaHandler(PropostaService)
	PropostaHandler.RegisterRoutes(router)
}