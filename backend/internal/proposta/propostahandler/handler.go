package propostahandler

import (
	"context"
	"encoding/json"
	"errors"
	"escambo/internal/proposta/propostarepo"
	"escambo/internal/proposta/propostasvc"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type PropostaSvc interface {
	GetPropostas(ctx context.Context, filter propostasvc.PropostasFilter) ([]propostarepo.PropostaFormatada, error)
	InsertProposta(ctx context.Context, proposta propostarepo.PropostaWriteModel) (*propostarepo.Proposta, error)
	UpdatePropostaStatus(ctx context.Context, propostaID string, status string) error
}

type Handler struct {
	svc PropostaSvc
}

func NewHandler(svc PropostaSvc) *Handler {
	return &Handler{svc: svc}
}

// GetPropostas godoc
// @Summary      Lista propostas do usuário
// @Description  Retorna propostas enviadas ou recebidas por um usuário com base no tipo e status
// @Tags         trocas
// @Param        id    path     string  true  "ID do usuário"
// @Param        tipo  query    string  true  "Tipo de proposta (enviadas ou recebidas)"
// @Param        status query   string  false "Status da proposta (pendente, aceita, recusada)"
// @Produce      json
// @Success      200  {array}   []propostarepo.PropostaFormatada
// @Failure      500  {string}  string
// @Router       /trocas/{id}/historico [get]
func (h *Handler) GetPropostas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	params, err := getParams(r)
	if err != nil {
		http.Error(w, "parametros invalidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	propostas, err := h.svc.GetPropostas(r.Context(), propostasvc.PropostasFilter{
		UsuarioID: id,
		Status:    params.Status,
		Tipo:      params.Tipo,
	})
	if err != nil {
		http.Error(w, "erro ao buscar propostas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(propostas); err != nil {
		http.Error(w, "erro ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
	}
}

// InsertProposta godoc
// @Summary      Cadastra nova proposta
// @Description  Registra uma proposta de troca com base nos dados enviados
// @Tags         trocas
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]string  "Proposta criada com sucesso e ID retornado"
// @Failure      400  {string}  string  "Body inválido"
// @Failure      500  {string}  string  "Erro ao salvar proposta"
// @Router       /trocas [post]
func (h *Handler) InsertProposta(w http.ResponseWriter, r *http.Request) {
	var proposta propostarepo.PropostaWriteModel

	if err := json.NewDecoder(r.Body).Decode(&proposta); err != nil {
		http.Error(w, "body inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.svc.InsertProposta(r.Context(), proposta)
	if err != nil {
		http.Error(w, "erro ao salvar proposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdatePropostaStatus godoc
// @Summary      Atualiza o status da proposta
// @Description  Altera o status da proposta (aceita, recusada, pendente)
// @Tags         trocas
// @Accept       json
// @Param        id    path     string  true  "ID da proposta"
// @Param        body  body     StatusUpdatePayload true "Novo status"
// @Success      204   {string} string  "Atualizado com sucesso"
// @Failure      400   {string} string  "Erro de validação"
// @Failure      500   {string} string  "Erro interno"
// @Router       /trocas/{id}/status [put]
func (h *Handler) UpdatePropostaStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var payload StatusUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "body inválido", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdatePropostaStatus(r.Context(), id, payload.Status); err != nil {
		http.Error(w, "erro ao atualizar status da proposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getParams(r *http.Request) (GetPropostasParams, error) {
	status := strings.ToLower(r.URL.Query().Get("status"))

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	tipo := strings.ToLower(r.URL.Query().Get("tipo"))
	if tipo == "" {
		return GetPropostasParams{}, errors.New("o parâmetro 'tipo' deve possuir os valores 'enviadas' ou 'recebidas'")
	}

	return GetPropostasParams{
		Status: statusPtr,
		Tipo:   tipo,
	}, nil
}
