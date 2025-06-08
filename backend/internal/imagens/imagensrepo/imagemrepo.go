package imagensrepo

import (
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

func (r Repository) SalvarImagem(dados Metadata) error {
	var currentJSON []byte
	var urls []string

	var querySelect string
	var queryUpdate string

	switch dados.Operacao {
	case "troca":
		querySelect = "SELECT imagem_url FROM propostas WHERE id = $1"
		queryUpdate = "UPDATE propostas SET imagem_url = $1 WHERE id = $2"

	case "postagem":
		querySelect = "SELECT imagem_url FROM postagens WHERE id = $1"
		queryUpdate = "UPDATE postagens SET imagem_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"

	default:
		return fmt.Errorf("operação desconhecida: %s", dados.Operacao)
	}

	err := r.DB.QueryRow(querySelect, dados.ID).Scan(&currentJSON)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if len(currentJSON) > 0 {
		err = json.Unmarshal(currentJSON, &urls)
		if err != nil {
			return fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
		}
	}

	urls = append(urls, dados.URL)

	newJSON, err := json.Marshal(urls)
	if err != nil {
		return fmt.Errorf("erro ao fazer marshal do JSON: %w", err)
	}

	_, err = r.DB.Exec(queryUpdate, string(newJSON), dados.ID)
	return err
}
