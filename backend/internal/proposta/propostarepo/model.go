package propostarepo

import (
	"errors"
	"time"
)

type PropostaWriteModel struct {
	PostagemID     string `json:"postagem_id"`
	RemetenteID    string `json:"interessado_id"`
	DestinatarioID string `json:"dono_postagem_id"`
	Descricao      string `json:"descricao"`
	Categoria      string `json:"categoria"`
	Nome           string `json:"nome"`
}

type PropostasQueryFilter struct {
	UsuarioID string
	Status    *string
	Tipo      string
}

type Produto struct {
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
	Descricao string `json:"descricao"`
	Usuario   string `json:"usuario"`
}

type PropostaFormatada struct {
	ProdutoPostagem      Produto `json:"produto_postagem"`
	ProdutoPropostaTroca Produto `json:"produto_proposta_troca"`
	Status               string  `json:"status"`
}

func (p PropostaWriteModel) Validate() error {
	if p.RemetenteID == "" {
		return errors.New("remetenteID é obrigatório")
	}
	if p.DestinatarioID == "" {
		return errors.New("destinatarioID é obrigatório")
	}

	return nil
}

type Proposta struct {
	ID             string    `json:"id"`
	PostagemID     string    `json:"postagem_id"`
	InteressadoID  string    `json:"interessado_id"`
	DonoPostagemID string    `json:"dono_postagem_id"`
	Descricao      string    `json:"descricao"`
	Nome           string    `json:"nome"`
	Categoria      string    `json:"categoria"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
