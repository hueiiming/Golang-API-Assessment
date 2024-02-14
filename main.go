package main

import (
	"Golang-API-Assessment/api"
	"Golang-API-Assessment/repository"
	"flag"
	"fmt"
	"log"
)

func main() {
	port := flag.String("port", ":3000", "server address")
	flag.Parse()
	repo := repository.NewPostgreSQLRepository()

	server := api.NewServer(*port, repo)
	fmt.Println("server running on port:", *port)
	log.Fatal(server.Start())
}
