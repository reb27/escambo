package usuariohandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"escambo/internal/usuario/usuariorepo"
	"escambo/internal/usuario/usuariosvc"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Handler struct {
	usuarioService *usuariosvc.Service
}

func NewHandler(usuarioService *usuariosvc.Service) *Handler {
	return &Handler{usuarioService: usuarioService}
}

// InsertUsuario godoc
// @Summary      Cadastra um novo usuário
// @Description  Insere um usuário no sistema com base nos dados fornecidos no corpo da requisição
// @Tags         usuarios
// @Accept       json
// @Produce      json
// @Param        proposta  body     usuariorepo.WriteUsuario  true  "Dados da proposta"
// @Success      201  {object}  map[string]string  "Usuário inserido com sucesso e ID retornado"
// @Failure      400  {string}  string  "Erro ao decodificar corpo da requisição"
// @Failure      409  {string}  string  "Já existe um cadastro com esse e-mail"
// @Failure      500  {string}  string  "Erro interno do servidor"
// @Router       /usuarios [post]
func (h *Handler) InsertUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario usuariorepo.WriteUsuario

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	id, err := h.usuarioService.InsertUsuario(r.Context(), usuario)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "já existe") {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Usuário inserido com sucesso",
		"id":      id,
	})
}

// UpdateUsuario godoc
// @Summary      Atualiza dados de um usuário
// @Description  Atualiza as informações de um usuário identificado pelo ID
// @Tags         usuarios
// @Accept       json
// @Produce      plain
// @Param        id       path   string                true  "ID do usuário"
// @Param        usuario  body   usuariorepo.WriteUsuario   true  "Dados atualizados do usuário"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {string}  string  "Erro ao decodificar corpo da requisição"
// @Failure      500  {string}  string  "Erro interno do servidor"
// @Router       /usuarios/{id} [put]
func (h *Handler) UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var usuario usuariorepo.WriteUsuario

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Erro ao decodificar corpo da requisição", http.StatusBadRequest)
		return
	}

	err = h.usuarioService.UpdateUsuario(r.Context(), id, usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUsuario godoc
// @Summary      Busca um usuário
// @Description  Retorna os dados de um usuário identificado pelo ID
// @Tags         usuarios
// @Accept       json
// @Produce      json
// @Param        id   path   string  true  "ID do usuário"
// @Success      200  {object}  usuariorepo.ReadUsuario
// @Failure      400  {string}  string  "ID inválido"
// @Failure      404  {string}  string  "Usuário não encontrado"
// @Failure      500  {string}  string  "Erro interno do servidor"
// @Router       /usuarios/{id} [get]
func (h *Handler) GetUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.usuarioService.GetUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// GetMe godoc
// @Summary      Retorna o usuário autenticado
// @Description  Retorna os dados do usuário com base no token de autenticação
// @Tags         usuarios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  usuariorepo.ReadUsuario
// @Failure      401  {string}  string  "Não autorizado"
// @Failure      500  {string}  string  "Erro interno do servidor"
// @Router       /me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("usuarioID").(string)
	if !ok || userID == "" {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	usuario, err := h.usuarioService.GetUsuario(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}
