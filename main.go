package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/hatrnuhn/pijar-crud/internal/database"
	"github.com/joho/godotenv"
)

type Config struct {
	db *database.DB
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DBPATH")

	rChi := chi.NewRouter()
	rDB := chi.NewRouter()

	var err error
	cfg := Config{}
	cfg.db, err = database.NewDB(dbPath)
	if err != nil {
		log.Fatal("couldn't initialize database")
	}

	rChi.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	rDB.Post("/produks", cfg.handleCreate)
	rDB.Get("/produks", cfg.handleRead)
	rDB.Put("/produks/{produkName}", cfg.handleUpdate)
	rDB.Delete("/produks/{produkName}", cfg.handleDel)

	rChi.Mount("/database", rDB)

	fmt.Printf("Starting server at: http://localhost:%s\n", port)
	http.ListenAndServe(port, rChi)
}
