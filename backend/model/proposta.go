package model

import (
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"propulse/shared/logger"
)

var hexColorRegex = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

type Proposta struct {
    Id           uuid.UUID `json:"id" validate:"uuid"`
    Titulo       string    `json:"titulo" validate:"required,min=3,max=100"`
    NomeEmpresa  string    `json:"nomeEmpresa" validate:"required"`
    NomeCliente  string    `json:"nomeCliente" validate:"required"`
    Prompt       string    `json:"prompt" validate:"required,min=20"`
    Cores        []string  `json:"cores" validate:"required,dive,hexcolor"`
    Logo         string    `json:"logo" validate:"omitempty,url"`
    LogoCliente  string    `json:"logoCliente" validate:"omitempty,url"`
    Status       string    `json:"status" validate:"required,oneof=rascunho enviado aprovado"`
    ArquivoFinal string    `json:"arquivoFinal"`
    DataCriacao  time.Time `json:"dataCriacao"`
    LastUpdate   time.Time `json:"lastUpdate"`
}

type PropostaUpdate struct {
    Titulo       *string   `json:"titulo" validate:"omitempty,min=3,max=100"`
    Status       *string `json:"status" validate:"omitempty,oneof=rascunho enviado aprovado"`
    ArquivoFinal *string    `json:"arquivoFinal"`
}

type RegerarProposta struct {
    NomeEmpresa  string    `json:"nomeEmpresa" validate:"required"`
    NomeCliente  string    `json:"nomeCliente" validate:"required"`
    Prompt       string    `json:"prompt" validate:"required,min=20"`
    Cores        []string  `json:"cores" validate:"required,dive,hexcolor"`
    Logo         string    `json:"logo" validate:"omitempty,url"`
    LogoCliente  string    `json:"logoCliente" validate:"omitempty,url"`
}

func HexColor(fl validator.FieldLevel) bool {
    color := fl.Field().String()
    return hexColorRegex.MatchString(color)
}

func ValidarStructProposta(p *Proposta) error {
    validate := validator.New()

	logger.Info("Validando cor hex")
    validate.RegisterValidation("hexcolor", func(fl validator.FieldLevel) bool {
        return HexColor(fl)
    })

	err := validate.Struct(p)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            logger.Error("Erro de validação no campo", err,
                zap.String("campo", err.Field()),
                zap.String("regra", err.Tag()),
                zap.String("erro", err.Error()),
            )
        }
        return err
    } else {
        logger.Info("Validação concluída com sucesso!")
        return nil
    }
}

func ValidarStructRegerarProposta(p *RegerarProposta) error {
    validate := validator.New()

	logger.Info("Validando cor hex")
    validate.RegisterValidation("hexcolor", func(fl validator.FieldLevel) bool {
        return HexColor(fl)
    })

	err := validate.Struct(p)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            logger.Error("Erro de validação no campo", err,
                zap.String("campo", err.Field()),
                zap.String("regra", err.Tag()),
                zap.String("erro", err.Error()),
            )
        }
        return err
    } else {
        logger.Info("Validação concluída com sucesso!")
        return nil
    }
}
