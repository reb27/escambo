package loginhandler

import (
	"encoding/json"
	"escambo/internal/login/loginsvc"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Handler struct {
	svc loginsvc.Service
}

func NewHandler(svc loginsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// Login godoc
// @Summary      Login de usuário
// @Description  Autentica um usuário e retorna token JWT para autenticação nas rotas protegidas.
// @Tags         autenticação
// @Accept       json
// @Produce      json
// @Param        loginRequest  body      LoginRequest  true  "Email e senha do usuário"
// @Success      200           {object}  LoginResponse "Token JWT gerado"
// @Failure      400           {string}  string        "Email e senha são obrigatórios"
// @Failure      401           {string}  string        "Email ou senha inválidos"
// @Failure      403           {string}  string        "Usuário temporariamente bloqueado"
// @Failure      500           {string}  string        "Erro interno"
// @Router       /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fmt.Println(req.Email)
	fmt.Println(req.Senha)
	if err != nil {
		http.Error(w, "Erro ao decodificar payload", http.StatusInternalServerError)
		return
	}
	if req.Email == "" || req.Senha == "" {
		http.Error(w, "Email e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	token, err := h.svc.Autenticar(req.Email, req.Senha)
	fmt.Println(err)

	if err != nil {
		fmt.Println(err)
		switch err {
		case loginsvc.ErrUsuarioNaoEncontrado, loginsvc.ErrSenhaInvalida:
			http.Error(w, "Email ou senha inválidos", http.StatusUnauthorized)
		case loginsvc.ErrUsuarioBloqueado:
			http.Error(w, "Usuário temporariamente bloqueado", http.StatusForbidden)
		default:
			http.Error(w, "Erro interno", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
