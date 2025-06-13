package postagemrepo

import (
	"time"
)

type Postagem struct {
	Titulo      string    `json:"titulo"`
	Descricao   string    `json:"descricao"`
	Imagens     []string  `json:"imagens"`
	UserID      string    `json:"user_id"`
	NomeUsuario string    `json:"nome_usuario"`
	Categoria   string    `json:"categoria"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Endereco
}

type Endereco struct {
	Cidade string `json:"cidade"`
	Estado string `json:"estado"`
	Bairro string `json:"bairro"`
}

type FiltroPostagem struct {
	Categoria    string
	Cidade       string
	Estado       string
	PalavraChave string
	Ordenacao    string
	De           time.Time
	Ate          time.Time
}
