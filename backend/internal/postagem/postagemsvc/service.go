package postagemsvc

import (
	"context"
	"escambo/internal/postagem/postagemrepo"
)

type PostagemRepository interface {
	InsertPostagem(ctx context.Context, post postagemrepo.Post) (string, error)
	GetPostagemByID(ctx context.Context, postID string) (postagemrepo.Post, error)
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

func (s Service) InsertPostagem(ctx context.Context, post Postagem) (string, error) {
	if err := post.Validate(); err != nil {
		return "", err
	}

	id, err := s.PostagemRepo.InsertPostagem(ctx, postagemrepo.Post{
		Titulo:       post.Titulo,
		Descricao:    post.Descricao,
		ImagemBase64: post.ImagemBase64,
		UserID:       post.UserID,
		Categoria:    post.Categoria,
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
		Titulo:       postagem.Titulo,
		Descricao:    postagem.Descricao,
		ImagemBase64: postagem.ImagemBase64,
		UserID:       postagem.UserID,
		Categoria:    postagem.Categoria,
	}, nil
}
