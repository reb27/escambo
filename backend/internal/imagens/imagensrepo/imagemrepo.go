package imagensrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

func (r Repository) RemoverImagemDaPostagem(ctx context.Context, postagemID string, targetURL string) error {
	var imagens []string

	querySelect := `SELECT imagem_url FROM postagens WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, querySelect, postagemID)

	var imagensJSON []byte
	if err := row.Scan(&imagensJSON); err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao buscar imagem_url: %w", err)
	}

	if err := json.Unmarshal(imagensJSON, &imagens); err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao fazer unmarshal do json: %w", err)
	}

	var novasImagens []string
	for _, img := range imagens {
		if img != targetURL {
			novasImagens = append(novasImagens, img)
		}
	}

	imagensAtualizadasJSON, err := json.Marshal(novasImagens)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("erro ao fazer marshal das novas imagens: %w", err)
	}

	queryUpdate := `UPDATE postagens SET imagem_url = $1 WHERE id = $2`
	_, err = r.DB.ExecContext(ctx, queryUpdate, imagensAtualizadasJSON, postagemID)
	if err != nil {
		log.Println(err)

		return fmt.Errorf("erro ao atualizar imagem_url: %w", err)
	}

	log.Println("SUCCESS update image list")

	return nil
}
