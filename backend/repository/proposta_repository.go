package repository

import (
	"context"
	"fmt"
	"propulse/model"
	"propulse/shared/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PropostaRepository struct{
	connection *pgxpool.Pool
}

func NewPropostaRepository(connection *pgxpool.Pool) PropostaRepository {
	return PropostaRepository{
		connection: connection,
	}
}

func (pr *PropostaRepository) CriarProposta(proposta model.Proposta) (*model.Proposta, error) {
    currentTime := time.Now()
	proposta.DataCriacao = currentTime
	proposta.LastUpdate = currentTime
    query := ` INSERT INTO propostas ( id, titulo, nome_empresa, nome_cliente, prompt, cores,
	   logo, logo_cliente, status, arquivo_final, data_criacao, last_update
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11, $12
        )
        RETURNING *;`

    row := pr.connection.QueryRow(
        context.Background(),
        query,
        proposta.Id,
        proposta.Titulo,
        proposta.NomeEmpresa,
        proposta.NomeCliente,
        proposta.Prompt,
        proposta.Cores,
        proposta.Logo,
        proposta.LogoCliente,
        proposta.Status,
        proposta.ArquivoFinal,
        proposta.DataCriacao,
        proposta.LastUpdate,
    )

    var p model.Proposta
    err := row.Scan(
        &p.Id,
        &p.Titulo,
        &p.NomeEmpresa,
        &p.NomeCliente,
        &p.Prompt,
        &p.Cores,
        &p.Logo,
        &p.LogoCliente,
        &p.Status,
        &p.ArquivoFinal,
        &p.DataCriacao,
        &p.LastUpdate,
    )
    if err != nil {
        return nil, err
    }
	logger.Info("Proposta criada com sucesso!")

    return &p, nil
}

func (pr *PropostaRepository) UpdateProposta(id uuid.UUID, update model.PropostaUpdate) (*model.Proposta, error) {
    setParts := []string{}
    args := []any{}
    argIndex := 1

    if update.Titulo != nil {
        setParts = append(setParts, fmt.Sprintf("titulo = $%d", argIndex))
        args = append(args, *update.Titulo)
        argIndex++
    }
    if update.Status != nil {
        setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
        args = append(args, *update.Status)
        argIndex++
    }
    if update.ArquivoFinal != nil {
        setParts = append(setParts, fmt.Sprintf("arquivo_final = $%d", argIndex))
        args = append(args, *update.ArquivoFinal)
        argIndex++
    }

    setParts = append(setParts, fmt.Sprintf("last_update = $%d", argIndex))
    args = append(args, time.Now())
    argIndex++

    args = append(args, id)

    query := fmt.Sprintf(
        "UPDATE propostas SET %s WHERE id = $%d RETURNING *",
        strings.Join(setParts, ", "),
        argIndex,
    )

    row := pr.connection.QueryRow(context.Background(), query, args...)
    
    var p model.Proposta
    err := row.Scan(&p.Id, &p.Titulo, &p.NomeEmpresa, &p.NomeCliente, 
                    &p.Prompt, &p.Cores, &p.Logo, &p.LogoCliente, 
                    &p.Status, &p.ArquivoFinal, &p.DataCriacao, &p.LastUpdate)
    if err != nil {
        return nil, err
    }

    return &p, nil
}


func (pr *PropostaRepository) FindByID(id uuid.UUID) (*model.Proposta, error) {
	query := `SELECT * FROM propostas WHERE id = $1`

	var p model.Proposta

	err := pr.connection.QueryRow(context.Background(), query, id).Scan(
            &p.Id,
            &p.Titulo,
            &p.NomeEmpresa,
            &p.NomeCliente,
            &p.Prompt,
            &p.Cores,
            &p.Logo,
            &p.LogoCliente,
            &p.Status,
            &p.ArquivoFinal,
            &p.DataCriacao,
            &p.LastUpdate,
        )
        if err != nil {
            logger.Error("Erro ao fazer scan da proposta", err)
            return &p, err
        }
		return &p, nil
}

func (pr *PropostaRepository) GetAllPropostas() (*[]model.Proposta, error) {
    query := `SELECT * FROM propostas`

    rows, err := pr.connection.Query(context.Background(), query)
    if err != nil {
        logger.Error("Erro ao buscar propostas", err)
        return &[]model.Proposta{}, err
    }

    defer rows.Close()

    var propostas []model.Proposta

    for rows.Next() {
        var p model.Proposta

        err := rows.Scan(
            &p.Id,
            &p.Titulo,
            &p.NomeEmpresa,
            &p.NomeCliente,
            &p.Prompt,
            &p.Cores,
            &p.Logo,
            &p.LogoCliente,
            &p.Status,
            &p.ArquivoFinal,
            &p.DataCriacao,
            &p.LastUpdate,
        )
        if err != nil {
            logger.Error("Erro ao fazer scan da proposta", err)
            return &[]model.Proposta{}, err
        }
        propostas = append(propostas, p)
    }

    if err = rows.Err(); err != nil {
        logger.Error("Erro durante iteração das linhas", err)
        return &[]model.Proposta{}, err
    }
    return &propostas, nil
}

func (pr *PropostaRepository) DeleteProposta(id uuid.UUID) (error) {
	query := `DELETE FROM propostas WHERE id = $1`

	_, err := pr.connection.Exec(context.Background(), query, id); if err != nil {
		logger.Error("Erro ao realizar a exclusão da proposta", err)
		return err
	}
	return nil
}


func (pr *PropostaRepository) UpdateForRegerar(id uuid.UUID, input model.RegerarProposta) (*model.Proposta, error) {
    setParts := []string{}
    args := []interface{}{}
    argIndex := 1

    setParts = append(setParts, fmt.Sprintf("nome_empresa = $%d", argIndex))
    args = append(args, input.NomeEmpresa)
    argIndex++

    setParts = append(setParts, fmt.Sprintf("nome_cliente = $%d", argIndex))
    args = append(args, input.NomeCliente)
    argIndex++

    setParts = append(setParts, fmt.Sprintf("prompt = $%d", argIndex))
    args = append(args, input.Prompt)
    argIndex++

    setParts = append(setParts, fmt.Sprintf("cores = $%d", argIndex))
    args = append(args, input.Cores) 
    argIndex++

    if input.Logo != "" {
        setParts = append(setParts, fmt.Sprintf("logo = $%d", argIndex))
        args = append(args, input.Logo)
        argIndex++
    }

    if input.LogoCliente != "" {
        setParts = append(setParts, fmt.Sprintf("logo_cliente = $%d", argIndex))
        args = append(args, input.LogoCliente)
        argIndex++
    }

    now := time.Now()
    setParts = append(setParts, fmt.Sprintf("last_update = $%d", argIndex))
    args = append(args, now)
    argIndex++

    args = append(args, id)

    query := fmt.Sprintf(
        "UPDATE propostas SET %s WHERE id = $%d RETURNING id, titulo, nome_empresa, nome_cliente, prompt, cores, logo, logo_cliente, status, arquivo_final, data_criacao, last_update",
        strings.Join(setParts, ", "),
        argIndex,
    )

    logger.Info("Executando UPDATE para regeneração", zap.String("query", query), zap.Int("args_count", len(args)))

    row := pr.connection.QueryRow(context.Background(), query, args...)

    var p model.Proposta
    err := row.Scan(
        &p.Id,
        &p.Titulo,           
        &p.NomeEmpresa,
        &p.NomeCliente,
        &p.Prompt,
        &p.Cores,           
        &p.Logo,
        &p.LogoCliente,
        &p.Status,          
        &p.ArquivoFinal,     
        &p.DataCriacao,      
        &p.LastUpdate,      
    )

    if err != nil {
        logger.Error("Erro ao executar UPDATE para regeneração", err, zap.String("id", id.String()))
        return nil, fmt.Errorf("falha ao atualizar proposta para regeneração: %w", err)
    }

    logger.Info("Proposta atualizada com sucesso para regeneração", zap.String("id", p.Id.String()))
    return &p, nil
}