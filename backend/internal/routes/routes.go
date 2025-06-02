package routes

import (
	"database/sql"
	"escambo/internal/imagens/imagenshandler"
	"escambo/internal/imagens/imagensrepo"
	"escambo/internal/imagens/imagenssvc"
	"escambo/internal/postagem/postagemhandler"
	"escambo/internal/postagem/postagemrepo"
	"escambo/internal/postagem/postagemsvc"
	"escambo/internal/proposta/propostahandler"
	"escambo/internal/proposta/propostarepo"
	"escambo/internal/proposta/propostasvc"
	"escambo/internal/usuario/usuariohandler"
	"escambo/internal/usuario/usuariorepo"
	"escambo/internal/usuario/usuariosvc"
	"log"
	"os"

	_ "escambo/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, db *sql.DB) {
	postRepo := postagemrepo.NewRepository(db)
	postService := postagemsvc.NewService(postRepo)
	postHandler := postagemhandler.NewHandler(postService)

	r.HandleFunc("/postagens/{id}/detalhes", postHandler.GetDetalhesPostagem).Methods("GET")
	r.HandleFunc("/postagens", postHandler.InsertPostagem).Methods("POST")

	userRepo := usuariorepo.NewRepository(db)
	usuarioService := usuariosvc.NewService(userRepo)
	usuarioHandler := usuariohandler.NewHandler(usuarioService)

	r.HandleFunc("/usuarios", usuarioHandler.InsertUsuario).Methods("POST")
	r.HandleFunc("/usuarios/{id}", usuarioHandler.UpdateUsuario).Methods("PUT")
	r.HandleFunc("/usuarios/{id}", usuarioHandler.GetUsuario).Methods("GET")

	propostaRepo := propostarepo.NewRepository(db)
	propostaService := propostasvc.NewService(propostaRepo)
	propostaHandler := propostahandler.NewHandler(&propostaService)

	r.HandleFunc("/trocas", propostaHandler.InsertProposta).Methods("POST")
	r.HandleFunc("/trocas/{id}/historico", propostaHandler.GetPropostas).Methods("GET")
	r.HandleFunc("/trocas/{id}/status", propostaHandler.UpdatePropostaStatus).Methods("PUT")

	bucket := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	imagensRepo := imagensrepo.NewRepository(db)
	imagensService, err := imagenssvc.NewService(bucket, region, imagensRepo)
	if err != nil {
		log.Fatalf("Erro ao criar imagensService: %v", err)
	}
	imagenshandler := imagenshandler.NewHandler(imagensService)

	r.HandleFunc("/trocas/{id}/imagem", imagenshandler.UploadImagemTroca).Methods("POST")
	r.HandleFunc("/postagens/{id}/imagem", imagenshandler.UploadImagemPostagem).Methods("POST")

	r.HandleFunc("/trocas/{id}/imagens", imagenshandler.GetImagensTroca).Methods("GET")
	r.HandleFunc("/postagens/{id}/imagens", imagenshandler.GetImagensPostagem).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

}
