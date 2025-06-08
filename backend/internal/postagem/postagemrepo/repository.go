package postagemrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r Repository) InsertPostagem(ctx context.Context, post Postagem) (string, error) {
	query := `
		INSERT INTO postagens (id, titulo, descricao, user_id, categoria)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT unique_titulo_user DO NOTHING
		RETURNING id;
	`

	var id string
	err := r.DB.QueryRowContext(ctx,
		query,
		post.Titulo,
		post.Descricao,
		post.UserID,
		post.Categoria,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return id, err
}

func (r Repository) GetPostagemByID(ctx context.Context, postID string) (Postagem, error) {
	query := `
		SELECT 
			titulo, 
			descricao, 
			user_id, 
			categoria,
			created_at, 
			updated_at,
			imagem_url
		FROM postagens
		WHERE id = $1;
	`

	var post Postagem
	var imagensJSON []byte

	err := r.DB.QueryRowContext(ctx, query, postID).Scan(
		&post.Titulo,
		&post.Descricao,
		&post.UserID,
		&post.Categoria,
		&post.CreatedAt,
		&post.UpdatedAt,
		&imagensJSON,
	)
	if err != nil {
		return Postagem{}, err
	}

	if len(imagensJSON) > 0 {
		if err := json.Unmarshal(imagensJSON, &post.Imagens); err != nil {
			return Postagem{}, fmt.Errorf("erro ao decodificar imagem_url: %w", err)
		}
	}

	return post, nil
}
