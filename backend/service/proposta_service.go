package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"propulse/model"
	"propulse/repository"
	"propulse/shared/logger"
	"slices"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PropostaService struct {
	repository repository.PropostaRepository
}

var iaURL = os.Getenv("IA_URL")

func NewPropostaService(pr repository.PropostaRepository) PropostaService {
	return PropostaService{
		repository: pr,
	}
}

func (ps *PropostaService) SalvarPDF(propostaID uuid.UUID, pdfData []byte) (string, error) {
	diretorioDestino := filepath.Join("uploads", "propostas")
	if err := os.MkdirAll(diretorioDestino, 0755); err != nil {
		logger.Error("Erro ao criar diretório de destino:", err, zap.String("Path", diretorioDestino))
		return "", err 
	}
	nomeArquivo := fmt.Sprintf("proposta_%s.pdf", propostaID.String())
	filePath := filepath.Join(diretorioDestino, nomeArquivo)
	err := os.WriteFile(filePath, pdfData, 0644); if err != nil {
		logger.Error("Erro ao salvar arquivo PDF no disco:", err)
		return "", err
	}
	logger.Info("PDF salvo com sucesso em:", zap.String("Path", filePath))
	return filePath, nil
}

func (ps *PropostaService) GetAllPropostas() (*[]model.Proposta, error) {
	listaDePropostas, err := ps.repository.GetAllPropostas()
	if err != nil {
		logger.Error("Erro ao consultar propostas", err)
		return &[]model.Proposta{}, err
	}
	return listaDePropostas, nil
}

func (ps *PropostaService) CriarProposta(propostaInput model.Proposta) (*model.Proposta, error) {
	propostaInput.Id = uuid.New()
	propostaOutput, err := ps.repository.CriarProposta(propostaInput)
	if err != nil {
		logger.Error("Erro ao criar proposta!", err)
		return &model.Proposta{}, err
	}
	body, err := json.Marshal(propostaOutput)
	if err != nil {
		logger.Error("Erro ao realizar o Marshal da Proposta", err)
		return &model.Proposta{}, err
	}
	payload := bytes.NewBuffer(body)
	resp, err := http.Post(iaURL + "/gerarproposta/pdf_dynamic", "application/json", payload)
	if err != nil {
		logger.Error("Erro ao gerar proposta", err)
		return &model.Proposta{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("Serviço de IA retornou um erro:", errors.New(resp.Status))
		errorBody, _ := io.ReadAll(resp.Body)
		logger.Error("Detalhe do erro da IA:", errors.New(string(errorBody)))
		return nil, fmt.Errorf("serviço de IA falhou: %s", resp.Status)
	}
	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Erro ao ler o PDF da resposta da IA:", err)
		return nil, err
	}
	filePath, err := ps.SalvarPDF(propostaOutput.Id, pdfBytes)
	if err != nil {
		logger.Error("Erro ao salvar o arquivo PDF:", err)
		return nil, err
	}
	novoStatus := "rascunho"
	updateData := model.PropostaUpdate{
        ArquivoFinal: &filePath,
        Status:       &novoStatus,
    }
	propostaAtualizada, err := ps.repository.UpdateProposta(propostaOutput.Id, updateData)
	if err != nil {
		logger.Error("Erro ao atualizar proposta com caminho do PDF:", err)
		return nil, err
	}
	return propostaAtualizada, nil
}

func (ps *PropostaService) UpdateProposta(id uuid.UUID, update model.PropostaUpdate) (*model.Proposta, error) {
	if update.Status != nil {
		validStatuses := []string{"rascunho", "enviado", "aprovado"}
		statusValido := slices.Contains(validStatuses, *update.Status)
		if !statusValido {
			logger.Error("Status inválido fornecido", nil)
			return nil, fmt.Errorf("status inválido: %s", *update.Status)
		}
	}
	propostaOutput, err := ps.repository.UpdateProposta(id, update)
	if err != nil {
		logger.Error("Erro ao atualizar proposta!", err)
		return nil, err
	}
	return propostaOutput, nil
}

func (ps *PropostaService) FindByID(ParamID string) (*model.Proposta, error) {
	if ParamID == "" {
		return &model.Proposta{}, errors.New("ID não pode ser nulo")
	}
	id, err := uuid.Parse(ParamID)
	if err != nil {
		logger.Error("id não é um UUID", err)
		return &model.Proposta{}, err
	}
	proposta, err := ps.repository.FindByID(id)
	if err != nil {
		logger.Error("Erro ao procurar proposta", err)
		return &model.Proposta{}, err
	}
	return proposta, nil
}

func (ps *PropostaService) DeleteProposta(idParam string) error {
	if idParam == "" {
		return errors.New("ID não pode ser nulo")
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("id não é um UUID", err)
		return err
	}
	err = ps.repository.DeleteProposta(id)
	if err != nil {
		logger.Error("Erro ao deletar proposta", err)
		return err
	}
	return nil
}

func (ps *PropostaService) RegerarProposta(idParam string, input model.RegerarProposta) (*model.Proposta, error) {
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("id não é um UUID válido", err)
		return nil, err
	}
	propostaAtualizada, err := ps.repository.UpdateForRegerar(id, input)
	if err != nil {
		logger.Error("Erro ao atualizar proposta para regerar", err)
		return nil, err
	}
	body, err := json.Marshal(propostaAtualizada)
	if err != nil {
		logger.Error("Erro ao realizar o Marshal da Proposta", err)
		return nil, err
	}

	payload := bytes.NewBuffer(body)
	fullApiUrl := iaURL + "/gerarproposta/pdf_dynamic"

	resp, err := http.Post(fullApiUrl, "application/json", payload)
	if err != nil {
		logger.Error("Erro ao chamar o serviço de IA", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Errorf("serviço de IA falhou: %s - %s", resp.Status, string(errorBody))
		logger.Error("Erro:", errorMsg)
		return nil, errorMsg
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Erro ao ler o PDF da resposta da IA", err)
		return nil, err
	}

	filePath, err := ps.SalvarPDF(id, pdfBytes)
	if err != nil {
		logger.Error("Erro ao salvar o novo arquivo PDF", err)
		return nil, err
	}

	updateData := model.PropostaUpdate{
		ArquivoFinal: &filePath,
	}
	
	propostaComPDF, err := ps.repository.UpdateProposta(id, updateData)
	if err != nil {
		logger.Error("Erro ao atualizar proposta com caminho do PDF regerado", err)
		return nil, err
	}

	return propostaComPDF, nil
}