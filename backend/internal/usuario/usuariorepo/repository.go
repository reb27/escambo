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
func (r Repository) InsertUsuario(ctx context.Context, usuario Usuario) (string, error) {
	query := `
        INSERT INTO usuarios (id, nome, email, senha, telefone)
        VALUES (gen_random_uuid(), $1, $2, $3, $4)
        ON CONFLICT (email) DO NOTHING
        RETURNING id;
    `

	var id string
	err := r.DB.QueryRowContext(ctx, query, usuario.Nome, usuario.Email, string(usuario.Senha), usuario.Telefone).Scan(&id)
	if err == sql.ErrNoRows {
		return "", errors.New("já existe um cadastro com esse e-mail")
	} else if err != nil {
		return "", err
	}

	return id, nil
}

func (r Repository) UpdateUsuario(ctx context.Context, id string, usuario Usuario) error {
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
