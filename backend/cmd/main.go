package main

import (
	"database/sql"
	"escambo/internal/email"
	"escambo/internal/routes"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Loading envs:", err)
	}

	postgresURI := os.Getenv("DATABASE_URL")
	if postgresURI == "" {
		log.Panic("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		log.Panic("Error opening database connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic("Error pinging the database: ", err)
	}

	go email.IniciarEnvioPeriodico(db, 1*time.Minute)

	r := mux.NewRouter()

	routes.RegisterRoutes(r, db)

	corsHandler := cors.New(cors.Options{AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor rodando na porta " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
	select {}

}
