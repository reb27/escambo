package routes

import (
	"context"
	"database/sql"
	"errors"
	"escambo/internal/imagens/imagenshandler"
	"escambo/internal/imagens/imagensrepo"
	"escambo/internal/imagens/imagenssvc"
	"escambo/internal/login/loginhandler"
	"escambo/internal/login/loginrepo"
	"escambo/internal/login/loginsvc"
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
	"net/http"
	"os"
	"strings"

	_ "escambo/docs"

	"github.com/golang-jwt/jwt"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, db *sql.DB, jwtSecret []byte) {
	postRepo := postagemrepo.NewRepository(db)
	postService := postagemsvc.NewService(postRepo)
	postHandler := postagemhandler.NewHandler(postService)

	authRoutes := r.PathPrefix("/").Subrouter()
	authRoutes.Use(AutenticacaoMiddleware(jwtSecret))

	authRoutes.HandleFunc("/postagens/{id}/detalhes", postHandler.GetDetalhesPostagem).Methods("GET")

	r.HandleFunc("/postagens", postHandler.InsertPostagem).Methods("POST")
	r.HandleFunc("/postagens", postHandler.GetPostagens).Methods("GET")
	r.HandleFunc("/postagens/{id}", postHandler.DeletarPostagem).Methods("DELETE")
	r.HandleFunc("/postagens", postHandler.UpdatePostagem).Methods("PUT")

	r.HandleFunc("/favoritos", postHandler.FavoritarPostagem).Methods("POST")
	r.HandleFunc("/favoritos", postHandler.DesfavoritarPostagem).Methods("DELETE")
	r.HandleFunc("/favoritos/{id}", postHandler.GetFavoritosByUserID).Methods("GET")

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
	r.HandleFunc("/postagens/{id}/imagem", imagenshandler.DeleteImagemPostagem).Methods("DELETE")

	authRepository := loginrepo.NewRepo(db)
	authService := loginsvc.NewService(authRepository, jwtSecret)
	authHandler := loginhandler.NewHandler(authService)

	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

}

func AutenticacaoMiddleware(chaveSecreta []byte) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "token não fornecido", http.StatusUnauthorized)
				return
			}

			usuarioID, err := validarTokenERetornarUsuarioID(token, chaveSecreta)
			if err != nil {
				http.Error(w, "token inválido", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "usuarioID", usuarioID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validarTokenERetornarUsuarioID(tokenStr string, chaveSecreta []byte) (string, error) {
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inesperado")
		}
		return chaveSecreta, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("token inválido")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if usuarioID, ok := claims["sub"].(string); ok {
			return usuarioID, nil
		}
	}

	return "", errors.New("usuario_id não encontrado no token")
}
