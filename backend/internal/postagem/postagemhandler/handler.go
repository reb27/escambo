package postagemhandler

import (
	"context"
	"encoding/json"
	"escambo/internal/postagem/postagemsvc"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type PostagemService interface {
	GetDetalhesPostagem(ctx context.Context, postID string) (postagemsvc.Postagem, error)
	InsertPostagem(ctx context.Context, post postagemsvc.Postagem) (string, error)
}

type Handler struct {
	Service PostagemService
}

func NewHandler(service PostagemService) *Handler {
	return &Handler{Service: service}
}

// GetDetalhesPostagem retorna os detalhes de uma postagem específica.
// @Summary      Buscar detalhes da postagem
// @Description  Retorna todas as informações de uma postagem com base no ID fornecido.
// @Tags         postagens
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID da Postagem"
// @Success      200  {object}  postagemsvc.Postagem
// @Failure      500  {string}  string  "Erro interno ao buscar a postagem ou ao codificar a resposta"
// @Router       /postagens/{id}/detalhes [get]
func (h *Handler) GetDetalhesPostagem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	post, err := h.Service.GetDetalhesPostagem(r.Context(), postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao buscar postagem: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, fmt.Sprintf("erro ao codificar resposta em JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// InsertPostagem godoc
// @Summary      Insere uma nova postagem
// @Description  Cria uma nova postagem no sistema com os dados fornecidos no corpo da requisição
// @Tags         postagens
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "Postagem inserida com sucesso"
// @Failure      400  {string}  string  "Erro ao decodificar corpo da requisição"
// @Failure      500  {string}  string  "Erro ao salvar postagem"
// @Router       /postagens [post]
func (h *Handler) InsertPostagem(w http.ResponseWriter, r *http.Request) {
	var post postagemsvc.Postagem

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, fmt.Sprintf("erro ao decodificar corpo da requisição: %v", err), http.StatusBadRequest)
		return
	}

	id, err := h.Service.InsertPostagem(r.Context(), post)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao salvar postagem: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Postagem inserida com sucesso",
		"id":      id,
	})
}
