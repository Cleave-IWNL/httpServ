package main

import (
	"log"
	"net/http"
	"os"

	"httpServ/internal/handler"
	"httpServ/internal/repository"
	"httpServ/internal/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepoPostgres(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service)
	r := handler.NewRouter(h)

	http.ListenAndServe(":8080", r)
}
