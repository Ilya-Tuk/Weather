package main

import (
	"github.com/Ilya-Tuk/Weather/internal/repositories/memory"
	"github.com/Ilya-Tuk/Weather/internal/services"
	"github.com/Ilya-Tuk/Weather/internal/transport/rest"
)

func main() {
	repo := &memory.Repository{}
	service := services.New(repo)
	service.Init()

	defer service.Close()

	rest.NewServer(service).Run(":8080")
}
