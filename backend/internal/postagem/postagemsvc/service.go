package postagemsvc

import (
	"context"
	"escambo/internal/postagem/postagemrepo"
	"fmt"
	"strings"
)

type PostagemRepository interface {
	InsertPostagem(ctx context.Context, post postagemrepo.Postagem) (string, error)
	GetPostagemByID(ctx context.Context, postID string) (postagemrepo.Postagem, error)
	GetPostagens(ctx context.Context, filtro postagemrepo.FiltroPostagem, offset int) ([]postagemrepo.Postagem, error)
	FavoritarPostagem(ctx context.Context, userID, postagemID string) error
	DesfavoritarPostagem(ctx context.Context, userID, postagemID string) error
	GetFavoritosByID(ctx context.Context, userID string) (postagemrepo.Favoritos, error)
	DeletarPostagem(ctx context.Context, postagemID string) error
	UpdatePostagem(ctx context.Context, postagem postagemrepo.PostagemEdicao) error
}

type Service struct {
	PostagemRepo PostagemRepository
}

func NewService(
	postRepo postagemrepo.Repository,
) Service {
	return Service{
		PostagemRepo: postRepo,
	}
}

func (s Service) InsertPostagem(ctx context.Context, postagem Postagem) (string, error) {
	if err := postagem.Validate(); err != nil {
		return "", err
	}

	id, err := s.PostagemRepo.InsertPostagem(ctx, postagemrepo.Postagem{
		Titulo:    postagem.Titulo,
		Descricao: postagem.Descricao,
		UserID:    postagem.UserID,
		Categoria: postagem.Categoria,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Service) GetDetalhesPostagem(ctx context.Context, postagemID string) (Postagem, error) {
	postagem, err := s.PostagemRepo.GetPostagemByID(ctx, postagemID)
	if err != nil {
		return Postagem{}, err
	}

	return Postagem{
		Titulo:    postagem.Titulo,
		Descricao: postagem.Descricao,
		UserID:    postagem.UserID,
		Categoria: postagem.Categoria,
		Imagens:   postagem.Imagens,
	}, nil

}

func (s Service) GetPostagens(ctx context.Context, filtro postagemrepo.FiltroPostagem) ([]postagemrepo.Postagem, error) {
	limite := filtro.Limite
	if limite <= 0 {
		filtro.Limite = 10
	}
	offset := (filtro.Pagina - 1) * limite
	if offset < 0 {
		offset = 0
	}

	filtro.Ordenacao = strings.ToLower(filtro.Ordenacao)
	if filtro.Ordenacao != "asc" && filtro.Ordenacao != "desc" {
		filtro.Ordenacao = "desc"
	}

	postagens, err := s.PostagemRepo.GetPostagens(ctx, filtro, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar postagens: %w", err)
	}

	return postagens, nil
}

func (s Service) UpdatePostagem(ctx context.Context, post postagemrepo.PostagemEdicao) error {
	err := s.PostagemRepo.UpdatePostagem(ctx, post)
	if err != nil {
		return fmt.Errorf("erro ao atualizar postagem: %w", err)
	}
	return nil
}

func (s Service) DeletarPostagem(ctx context.Context, postagemID string) error {
	return s.PostagemRepo.DeletarPostagem(ctx, postagemID)
}

func (s Service) FavoritarPostagem(ctx context.Context, userID, postagemID string) error {
	return s.PostagemRepo.FavoritarPostagem(ctx, userID, postagemID)
}

func (s Service) DesfavoritarPostagem(ctx context.Context, userID, postagemID string) error {
	return s.PostagemRepo.DesfavoritarPostagem(ctx, userID, postagemID)
}

func (s Service) GetFavoritosByID(ctx context.Context, userID string) (postagemrepo.Favoritos, error) {
	return s.PostagemRepo.GetFavoritosByID(ctx, userID)
}
