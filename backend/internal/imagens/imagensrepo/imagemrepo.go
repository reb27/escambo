package imagensrepo

import "database/sql"

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
