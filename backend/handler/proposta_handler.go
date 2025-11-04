package handler

import (
	"net/http"
	"propulse/model"
	"propulse/service"
	"propulse/shared/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PropostaHandler struct {
	propostaService service.PropostaService
}

func NewPropostaHandler(service service.PropostaService) PropostaHandler {
	return PropostaHandler{
		propostaService: service,
	}
}

func (p *PropostaHandler) CriarProposta(ctx *gin.Context) {
	var proposta model.Proposta
	err := ctx.BindJSON(&proposta); if err != nil {
		logger.Error("Erro ao realizar o bind do JSON", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	err = model.ValidarStructProposta(&proposta); if err != nil {
		logger.Error("Erro ao passar no validador de Struct da Proposta", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	propostaOutput, err := p.propostaService.CriarProposta(proposta); if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, propostaOutput)
}

func (p *PropostaHandler) GetAllPropostas(ctx *gin.Context) {
	listasDePropostas, err := p.propostaService.GetAllPropostas(); if err != nil {
		logger.Error("Erro ao buscar propostas", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, listasDePropostas)
}

func (p *PropostaHandler) FindByID(ctx *gin.Context) {
	ParamID := ctx.Param("id")
	proposta, err := p.propostaService.FindByID(ParamID); if err != nil {
		logger.Error("Erro para encontrar proposta", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, proposta)
}

func (p *PropostaHandler) UpdateProposta(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }
    var update model.PropostaUpdate
    if err := ctx.BindJSON(&update); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    proposta, err := p.propostaService.UpdateProposta(id, update)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, proposta)
}

func (p *PropostaHandler) DeleteProposta(ctx *gin.Context) {
	idParam := ctx.Param("id")
	err := p.propostaService.DeleteProposta(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (p *PropostaHandler) RegerarProposta(ctx *gin.Context) {
	var propostaInput model.RegerarProposta
	idParam := ctx.Param("id")
	err := ctx.BindJSON(&propostaInput); if err != nil {
		logger.Error("Erro ao realizar o bind do JSON", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	err = model.ValidarStructRegerarProposta(&propostaInput); if err != nil {
		logger.Error("Erro ao passar no validador de Struct da Proposta", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	propostaOutput, err := 	p.propostaService.RegerarProposta(idParam, propostaInput); if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, propostaOutput)
}

func (h *PropostaHandler) RegisterRoutes(router *gin.Engine) {
	propostaRoutes := router.Group("/proposta")
	{
		propostaRoutes.POST("/", h.CriarProposta)
		propostaRoutes.POST("/:id/regerar")
		propostaRoutes.GET("/", h.GetAllPropostas)
		propostaRoutes.GET("/:id", h.FindByID)
		propostaRoutes.PATCH("/:id", h.UpdateProposta)
		propostaRoutes.DELETE("/:id", h.DeleteProposta)
	}
}