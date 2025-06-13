package postagemsvc

import (
	"context"
	"errors"
	"escambo/internal/postagem/postagemrepo"
)

type PostagemRepository interface {
	InsertPostagem(ctx context.Context, post postagemrepo.Postagem) (string, error)
	GetPostagemByID(ctx context.Context, postID string) (postagemrepo.Postagem, error)
	GetPostagens(ctx context.Context, filtro postagemrepo.FiltroPostagem) ([]postagemrepo.Postagem, error)
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
	postagens, err := s.PostagemRepo.GetPostagens(ctx, filtro)
	if err != nil {
		return []postagemrepo.Postagem{}, errors.New("erro ao listar postagens")
	}

	return postagens, nil
}
