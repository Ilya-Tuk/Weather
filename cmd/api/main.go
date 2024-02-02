package main

import (
	"github.com/Ilya-Tuk/Weather/internal/transport/rest"
	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/repositories/memory"
	"github.com/Ilya-Tuk/Weather/internal/services"
)

func main() {
	repo := &memory.Repository{}
	service := services.New(repo)

	rest.NewServer(service).Run(":8080")
}