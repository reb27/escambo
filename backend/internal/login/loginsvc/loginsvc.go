package loginsvc

import (
	"errors"
	"escambo/internal/login/loginrepo"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Autenticar(email, senha string) (string, error) // retorna token JWT
}

type service struct {
	repo   loginrepo.Repository
	secret []byte
}

func NewService(repo loginrepo.Repository, jwtSecret []byte) Service {
	return &service{repo: repo, secret: jwtSecret}
}

var (
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
	ErrSenhaInvalida        = errors.New("senha incorreta")
	ErrUsuarioBloqueado     = errors.New("usuário temporariamente bloqueado")
)

func (s *service) Autenticar(email, senha string) (string, error) {
	usuario, err := s.repo.BuscarUsuarioPorEmail(email)
	if err != nil {
		return "", err
	}
	if usuario == nil {
		return "", ErrUsuarioNaoEncontrado
	}

	if usuario.BloqueadoAte != nil && time.Now().UTC().Before(usuario.BloqueadoAte.UTC()) {
		return "", ErrUsuarioBloqueado
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuario.SenhaHash), []byte(senha))
	if err != nil {
		usuario.Tentativas++
		var bloquear *time.Time
		if usuario.Tentativas >= 5 {
			tmp := time.Now().UTC().Add(15 * time.Minute)
			bloquear = &tmp
		}
		_ = s.repo.IncrementarTentativas(usuario.ID, usuario.Tentativas, bloquear)
		return "", ErrSenhaInvalida
	}

	_ = s.repo.ResetarTentativas(usuario.ID)

	token, err := s.gerarToken(usuario.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) gerarToken(usuarioID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuarioID,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString(s.secret)
}
