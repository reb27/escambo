package routes

import (
	"database/sql"
	"escambo/internal/postagem/postagemhandler"
	"escambo/internal/postagem/postagemrepo"
	"escambo/internal/postagem/postagemsvc"
	"escambo/internal/proposta/propostahandler"
	"escambo/internal/proposta/propostarepo"
	"escambo/internal/proposta/propostasvc"
	"escambo/internal/usuario/usuariohandler"
	"escambo/internal/usuario/usuariorepo"
	"escambo/internal/usuario/usuariosvc"

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

	propostaRepo := propostarepo.NewRepository(db)
	propostaService := propostasvc.NewService(propostaRepo)
	propostaHandler := propostahandler.NewHandler(&propostaService)

	r.HandleFunc("/trocas", propostaHandler.InsertProposta).Methods("POST")
	r.HandleFunc("/trocas/{id}/historico", propostaHandler.GetPropostas).Methods("GET")
	r.HandleFunc("/trocas/{id}/status", propostaHandler.UpdatePropostaStatus).Methods("PUT")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

}
