package propostasvc

import (
	"context"
	"escambo/internal/proposta/propostarepo"
	"fmt"
)

type PropostaRepository interface {
	InsertProposta(ctx context.Context, proposta propostarepo.PropostaWriteModel) (*propostarepo.Proposta, error)
	UpdatePropostaStatus(ctx context.Context, propostaID string, status string) error
	GetPropostas(ctx context.Context, filter propostarepo.PropostasQueryFilter) ([]propostarepo.PropostaFormatada, error)
}

type Service struct {
	repo PropostaRepository
}

func NewService(repo PropostaRepository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) GetPropostas(ctx context.Context, filter PropostasFilter) ([]propostarepo.PropostaFormatada, error) {
	return s.repo.GetPropostas(ctx, propostarepo.PropostasQueryFilter{
		UsuarioID: filter.UsuarioID,
		Status:    filter.Status, //opcional
		Tipo:      filter.Tipo,   //obrigatorio
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	})
}

func (s *Service) UpdatePropostaStatus(ctx context.Context, propostaID string, status string) error {
	if !statusMap[status] {
		return fmt.Errorf("status inválido: '%s'. Os valores permitidos são: 'aceita' ou 'recusada'", status)
	}

	return s.repo.UpdatePropostaStatus(ctx, propostaID, status)
}

func (s *Service) InsertProposta(ctx context.Context, proposta propostarepo.PropostaWriteModel) (*propostarepo.Proposta, error) {
	if err := proposta.Validate(); err != nil {
		return nil, err
	}

	return s.repo.InsertProposta(ctx, proposta)
}
