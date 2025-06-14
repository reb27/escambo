package postagemhandler

import (
	"context"
	"encoding/json"
	"escambo/internal/postagem/postagemrepo"
	"escambo/internal/postagem/postagemsvc"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostagemService interface {
	GetDetalhesPostagem(ctx context.Context, postID string) (postagemsvc.Postagem, error)
	InsertPostagem(ctx context.Context, post postagemsvc.Postagem) (string, error)
	GetPostagens(ctx context.Context, filtro postagemrepo.FiltroPostagem) ([]postagemrepo.Postagem, error)
	FavoritarPostagem(ctx context.Context, userID, postagemID string) error
	DesfavoritarPostagem(ctx context.Context, userID, postagemID string) error
	GetFavoritosByID(ctx context.Context, userID string) (postagemrepo.Favoritos, error)
	DeletarPostagem(ctx context.Context, postagemID string) error
	UpdatePostagem(ctx context.Context, post postagemrepo.PostagemEdicao) error
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
// @Param        body  body     postagemsvc.PostagemWrite true "Nova postagem"
// @Success      201   {object} PostagemResponse "Postagem inserida com sucesso"
// @Failure      400   {string} string "Erro ao decodificar corpo da requisição"
// @Failure      500   {string} string "Erro ao salvar postagem"
// @Router       /postagens [post]
func (h *Handler) InsertPostagem(w http.ResponseWriter, r *http.Request) {
	var post postagemsvc.PostagemWrite

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, fmt.Sprintf("erro ao decodificar corpo da requisição: %v", err), http.StatusBadRequest)
		return
	}

	postInsert := postagemsvc.Postagem{
		Titulo:    post.Titulo,
		Descricao: post.Descricao,
		UserID:    post.UserID,
		Categoria: post.Categoria,
	}

	id, err := h.Service.InsertPostagem(r.Context(), postInsert)
	if err != nil {
		http.Error(w, fmt.Sprintf("erro ao salvar postagem: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PostagemResponse{
		Message: "Postagem inserida com sucesso",
		ID:      id,
	})
}

// DeletarPostagem godoc
// @Summary Deletar uma postagem
// @Description Remove uma postagem do sistema.
// @Tags postagens
// @Produce json
// @Param id path string true "ID da postagem"
// @Success 204 "Postagem deletada com sucesso"
// @Failure 400 {string} string "ID inválido"
// @Failure 500 {string} string "Erro ao deletar postagem"
// @Router /postagens/{id} [delete]
func (h Handler) DeletarPostagem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	postagemID := vars["id"]
	if postagemID == "" {
		http.Error(w, "id da postagem é obrigatório", http.StatusBadRequest)
		return
	}

	err := h.Service.DeletarPostagem(ctx, postagemID)
	if err != nil {
		http.Error(w, "Erro ao deletar postagem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPostagens godoc
// @Summary Lista postagens com filtros opcionais
// @Description Retorna uma lista de postagens de acordo com os filtros fornecidos
// @Tags postagens
// @Accept json
// @Produce json
// @Param categoria query string false "Categoria da postagem"
// @Param ordenacao query string false "Ordernacao das datas DESC ou ASC"
// @Param limite query int false "Número máximo de postagens por página"
// @Param pagina query int false "Número da página de resultados"
// @Success 200 {array} postagemrepo.Postagem
// @Failure      500   {string} string "Erro ao listar postagens"
// @Router /postagens [get]
func (h Handler) GetPostagens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	filtro := postagemrepo.FiltroPostagem{
		Categoria: query.Get("categoria"),
		Ordenacao: query.Get("ordenacao"),
	}

	if limiteStr := query.Get("limite"); limiteStr != "" {
		if limite, err := strconv.Atoi(limiteStr); err == nil {
			filtro.Limite = limite
		}
	}
	if paginaStr := query.Get("pagina"); paginaStr != "" {
		if pagina, err := strconv.Atoi(paginaStr); err == nil {
			filtro.Pagina = pagina
		}
	}

	postagens, err := h.Service.GetPostagens(ctx, filtro)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar postagens: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postagens)
}

// FavoritarPostagem godoc
// @Summary Favoritar uma postagem
// @Description Marca uma postagem como favorita para um usuário.
// @Tags favoritos
// @Accept json
// @Produce json
// @Param request body FavoritarPostagemRequest true "Dados para favoritar postagem"
// @Success 201
// @Failure 400 {string} string "Request inválido"
// @Failure 500 {string} string "Erro ao favoritar"
// @Router /favoritos [post]
func (h Handler) FavoritarPostagem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req FavoritarPostagemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	err := h.Service.FavoritarPostagem(ctx, req.UserID, req.PostagemID)
	if err != nil {
		http.Error(w, "Erro ao favoritar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DesfavoritarPostagem godoc
// @Summary Desfavoritar uma postagem
// @Description Remove uma postagem dos favoritos de um usuário.
// @Tags favoritos
// @Accept json
// @Produce json
// @Param request body DesfavoritarPostagemRequest true "Dados para desfavoritar postagem"
// @Success 200
// @Failure 400 {string} string "Request inválido"
// @Failure 500 {string} string "Erro ao desfavoritar"
// @Router /favoritos [delete]
func (h Handler) DesfavoritarPostagem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req DesfavoritarPostagemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	err := h.Service.DesfavoritarPostagem(ctx, req.UserID, req.PostagemID)
	if err != nil {
		http.Error(w, "Erro ao desfavoritar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFavoritosByUserID godoc
// @Summary Lista favoritos de um usuário
// @Description Retorna a lista de favoritos para um usuário específico.
// @Tags favoritos
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 200 {array} postagemrepo.Favorito
// @Failure 400 {string} string "userID obrigatório"
// @Failure 500 {string} string "Erro ao buscar favoritos"
// @Router /favoritos/{id} [get]
func (h Handler) GetFavoritosByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		http.Error(w, "id obrigatório", http.StatusBadRequest)
		return
	}

	favoritos, err := h.Service.GetFavoritosByID(ctx, userID)
	if err != nil {
		http.Error(w, "Erro ao buscar favoritos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favoritos)
}

// UpdatePostagem atualiza uma postagem existente.
// @Summary Atualiza uma postagem
// @Description Atualiza os dados de uma postagem específica com base no JSON recebido
// @Tags postagens
// @Accept json
// @Produce json
// @Param postagem body postagemrepo.PostagemEdicao true "Dados para edição da postagem"
// @Success 200  "Mensagem de sucesso"
// @Failure 400  "Requisição inválida"
// @Failure 500 "Erro interno ao atualizar postagem"
// @Router /postagens [put]
func (h Handler) UpdatePostagem(w http.ResponseWriter, r *http.Request) {
	var post postagemrepo.PostagemEdicao
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdatePostagem(r.Context(), post)
	if err != nil {
		http.Error(w, "Erro ao atualizar postagem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Postagem atualizada com sucesso"})
}
