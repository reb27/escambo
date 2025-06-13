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
func (r Repository) GetPostagens(ctx context.Context, filtro FiltroPostagem) ([]Postagem, error) {
	query := `
		SELECT 
			p.titulo, 
			p.descricao, 
			p.user_id, 
			p.categoria,
			p.created_at, 
			p.updated_at,
			p.imagem_url,
			u.nome as nome_usuario,
			e.cidade,
			e.estado,
			e.bairro
		FROM postagens p
		JOIN usuarios u ON u.id = p.user_id
		JOIN endereco e ON e.user_id = u.id
		WHERE p.postagem_ativa = TRUE
			AND ($1::text IS NULL OR p.categoria = $1)
			AND ($2::text IS NULL OR e.cidade ILIKE '%' || $2 || '%')
			AND ($3::text IS NULL OR e.estado ILIKE '%' || $3 || '%')
			AND ($4::text IS NULL OR p.titulo ILIKE '%' || $4 || '%' OR p.descricao ILIKE '%' || $4 || '%')
			AND ($5::timestamp IS NULL OR p.created_at >= $5)
			AND ($6::timestamp IS NULL OR p.created_at <= $6)
		ORDER BY 
			CASE WHEN $7 = 'data' THEN p.created_at
			     WHEN $7 = 'relevancia' THEN p.updated_at
			     ELSE p.created_at
			END DESC
	`

	rows, err := r.DB.QueryContext(ctx, query,
		filtro.Categoria,
		filtro.Cidade,
		filtro.Estado,
		filtro.PalavraChave,
		filtro.De,
		filtro.Ate,
		filtro.Ordenacao,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var postagens []Postagem

	for rows.Next() {
		var post Postagem
		var imagensJSON []byte

		err := rows.Scan(
			&post.Titulo,
			&post.Descricao,
			&post.UserID,
			&post.Categoria,
			&post.CreatedAt,
			&post.UpdatedAt,
			&imagensJSON,
			&post.NomeUsuario,
			&post.Cidade,
			&post.Estado,
			&post.Bairro,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear linha: %w", err)
		}

		if len(imagensJSON) > 0 {
			if err := json.Unmarshal(imagensJSON, &post.Imagens); err != nil {
				return nil, fmt.Errorf("erro ao decodificar imagem_url: %w", err)
			}
		}

		postagens = append(postagens, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro no loop de leitura: %w", err)
	}

	return postagens, nil
}
