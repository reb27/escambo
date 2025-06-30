package postagemrepo

import (
	"time"
)

type Postagem struct {
	ID                  string    `json:"id"`
	Titulo              string    `json:"titulo"`
	Descricao           string    `json:"descricao"`
	Imagens             []string  `json:"imagens"`
	UserID              string    `json:"user_id"`
	Status              bool      `json:"ativa"`
	NomeUsuario         string    `json:"nome_usuario"`
	Categoria           string    `json:"categoria"`
	CreatedAt           time.Time `json:"criacao_em"`
	MarcadoComoFavorito bool      `json:"marcado_como_favorito"`
	Endereco
}

type PostagemEdicao struct {
	ID        string  `json:"id,omitempty"`
	Titulo    *string `json:"titulo"`
	Descricao *string `json:"descricao"`
	Status    *bool   `json:"status"`
	Categoria *string `json:"categoria"`
}

type Endereco struct {
	Cidade string `json:"cidade"`
	Estado string `json:"estado"`
	Bairro string `json:"bairro"`
}

type FiltroPostagem struct {
	Categoria  string
	Ordenacao  string
	UsuarioID  string
	Limite     int
	Pagina     int
	BuscaTexto string
}

type Favorito struct {
	ID         string    `json:"id"`
	PostagemID string    `json:"postagem_id"`
	CriadoEm   time.Time `json:"criado_em"`
}

type Favoritos []Postagem
