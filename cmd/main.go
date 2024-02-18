package main

import (
	"Golang-API-Assessment/pkg/api"
	"Golang-API-Assessment/pkg/repository"
	"flag"
	"log"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// flag listenAddr to allow changing of ports in CLI for local machine
	listenAddr := flag.String("port", "0.0.0.0:"+port, "server address")
	flag.Parse()

	repo, err := repository.NewPostgreSQLRepository()
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(*listenAddr, repo)
	log.Printf("Server running on address: %s", *listenAddr)
	log.Fatal(server.Start())
}
