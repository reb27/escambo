package imagensrepo

import (
	"context"
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r Repository) SalvarImagem(dados Metadata) error {
	_, err := r.DB.Exec(`
		INSERT INTO post_images (operacao_id, operacao, url)
		VALUES ($1, $2, $3)
	`, dados.ID, dados.Operacao, dados.URL)

	return err
}

func (r Repository) ListarImagens(ctx context.Context, operacaoID, operacao string) ([]string, error) {
	query := `
        SELECT url FROM post_images
        WHERE operacao_id = $1 AND operacao = $2
        ORDER BY created_at
    `

	rows, err := r.DB.QueryContext(ctx, query, operacaoID, operacao)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}
