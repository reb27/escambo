package usuariorepo

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r Repository) InsertUsuario(ctx context.Context, usuario WriteUsuario) (string, error) {
	queryUsuario := `
        INSERT INTO usuarios (nome, email, senha, telefone)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (email) DO NOTHING
        RETURNING id;
    `

	var id string
	err := r.DB.QueryRowContext(ctx, queryUsuario,
		usuario.Nome,
		usuario.Email,
		string(usuario.Senha),
		usuario.Telefone,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return "", errors.New("já existe um cadastro com esse e-mail")
	} else if err != nil {
		return "", err
	}

	queryEndereco := `
        INSERT INTO endereco (cep, rua, numero, bairro, cidade, estado, user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (user_id) DO NOTHING;
    `

	_, err = r.DB.ExecContext(ctx, queryEndereco,
		usuario.CEP,
		usuario.Rua,
		usuario.Numero,
		usuario.Bairro,
		usuario.Cidade,
		usuario.Estado,
		id,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r Repository) UpdateUsuario(ctx context.Context, id string, usuario WriteUsuario) error {
	query := `
        UPDATE usuarios
        SET nome = $1, telefone = $2, email = $3, updated_at = NOW()
        WHERE id = $4;
    `
	_, err := r.DB.ExecContext(ctx, query, usuario.Nome, usuario.Telefone, usuario.Email, id)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetUsuario(ctx context.Context, id string) (ReadUsuario, error) {
	query := `
		SELECT id, nome, email, telefone FROM usuarios WHERE id = $1
	`

	var usuario ReadUsuario

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Email,
		&usuario.Telefone,
	)
	if err != nil {
		return ReadUsuario{}, err
	}

	return usuario, nil
}
