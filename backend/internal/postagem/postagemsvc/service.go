package postagemsvc

import (
	"context"
	"escambo/internal/postagem/postagemrepo"
)

type PostagemRepository interface {
	InsertPostagem(ctx context.Context, post postagemrepo.Postagem) (string, error)
	GetPostagemByID(ctx context.Context, postID string) (postagemrepo.Postagem, error)
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
		Imagens:   postagem.Imagens,
		UserID:    postagem.UserID,
		Categoria: postagem.Categoria,
	}, nil

}
