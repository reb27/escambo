package postagemsvc

import (
	"errors"
)

type Postagem struct {
	Titulo    string   `json:"titulo"`
	Descricao string   `json:"descricao"`
	Imagens   []string `json:"imagens"`
	UserID    string   `json:"user_id"`
	Categoria string   `json:"categoria"`
}

func (p *Postagem) Validate() error {
	if p.Titulo == "" {
		return errors.New("o título é obrigatório")
	}

	if p.Descricao == "" {
		return errors.New("a descrição é obrigatória")
	}

	if p.Categoria == "" {
		return errors.New("a categoria é obrigatória")
	}

	return nil
}
