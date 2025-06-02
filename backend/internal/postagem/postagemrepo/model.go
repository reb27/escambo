package postagemrepo

import (
	"time"
)

type Postagem struct {
	Titulo    string    `json:"titulo"`
	Descricao string    `json:"descricao"`
	Imagens   []string  `json:"imagens"`
	UserID    string    `json:"user_id"`
	Categoria string    `json:"categoria"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
