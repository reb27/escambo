package imagenshandler

import (
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

type ImagensService interface {
	UploadImagem(file multipart.File, header *multipart.FileHeader, ID, operacao string) (string, error)
}

type Handler struct {
	Service ImagensService
}

func NewHandler(service ImagensService) *Handler {
	return &Handler{Service: service}
}

// UploadImagemPostagem faz upload de uma imagem associada a uma postagem específica.
// @Summary Upload de imagem de postagem
// @Description Recebe um arquivo de imagem e o UUID de uma postagem. Retorna status 202 em caso de sucesso.
// @Tags         postagens
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "UUID da postagem"
// @Param image formData file true "Arquivo de imagem (máx 10MB)"
// @Success 202 "Status Accepted (202)"
// @Failure 400 {string} string "Erro na leitura da imagem ou ID inválido"
// @Failure 500 {string} string "Erro interno no upload"
// @Router /postagens/{id}/imagem [post]
func (h *Handler) UploadImagemPostagem(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erro ao ler imagem: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	vars := mux.Vars(r)
	postagemID := vars["id"]

	_, err = h.Service.UploadImagem(file, header, postagemID, "postagem")
	if err != nil {
		http.Error(w, "Erro ao fazer upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// UploadImagemTroca faz upload de uma imagem associada a uma troca específica.
// @Summary Upload de imagem de troca
// @Description Recebe um arquivo de imagem e o UUID de uma troca. Retorna status 202 em caso de sucesso.
// @Tags         trocas
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "UUID da troca"
// @Param image formData file true "Arquivo de imagem (máx 10MB)"
// @Success 202 "Status Accepted (202)"
// @Failure 400 {string} string "Erro na leitura da imagem ou ID inválido"
// @Failure 500 {string} string "Erro interno no upload"
// @Router /trocas/{id}/imagem [post]
func (h *Handler) UploadImagemTroca(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erro ao ler imagem: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	vars := mux.Vars(r)
	trocaID := vars["id"]

	_, err = h.Service.UploadImagem(file, header, trocaID, "troca")
	if err != nil {
		http.Error(w, "Erro ao fazer upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
