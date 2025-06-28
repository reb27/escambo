package loginrepo

import (
	"database/sql"
	"time"
)

type Usuario struct {
	ID           string
	Email        string
	SenhaHash    string
	Tentativas   int
	BloqueadoAte *time.Time
}

type Repository interface {
	BuscarUsuarioPorEmail(email string) (*Usuario, error)
	IncrementarTentativas(usuarioID string, tentativas int, bloquear *time.Time) error
	ResetarTentativas(usuarioID string) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repo{db}
}

func (r *repo) BuscarUsuarioPorEmail(email string) (*Usuario, error) {
	var u Usuario
	err := r.db.QueryRow(`
		SELECT id, email, senha, tentativas_login, bloqueado_ate
		FROM usuarios
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.SenhaHash, &u.Tentativas, &u.BloqueadoAte)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *repo) IncrementarTentativas(usuarioID string, tentativas int, bloquear *time.Time) error {
	if bloquear != nil {
		_, err := r.db.Exec(`
			UPDATE usuarios SET tentativas_login = $1, bloqueado_ate = $2 WHERE id = $3
		`, tentativas, *bloquear, usuarioID)
		return err
	}
	_, err := r.db.Exec(`
		UPDATE usuarios SET tentativas_login = $1 WHERE id = $2
	`, tentativas, usuarioID)
	return err
}

func (r *repo) ResetarTentativas(usuarioID string) error {
	_, err := r.db.Exec(`
		UPDATE usuarios SET tentativas_login = 0, bloqueado_ate = NULL WHERE id = $1
	`, usuarioID)
	return err
}
