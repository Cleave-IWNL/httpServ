package main

import (
	"net/http"

	"httpServ/internal/handler"
	"httpServ/internal/repository"
	"httpServ/internal/service"
)

func main() {
	repo := repository.NewRepoInMemory()
	service := service.NewService(repo)
	h := handler.NewHandler(service)
	r := handler.NewRouter(h)

	http.ListenAndServe(":8080", r)
}
