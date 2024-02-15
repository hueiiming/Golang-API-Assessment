package api

import (
	"Golang-API-Assessment/repository"
	"net/http"
)

type Server struct {
	port string
	repo repository.Repository
}

func NewServer(port string, repo repository.Repository) *Server {
	return &Server{
		port: port,
		repo: repo,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/api/register", makeHTTPHandle(s.handleRegister))
	http.HandleFunc("/api/commonstudents", makeHTTPHandle(s.handleCommonStudents))
	http.HandleFunc("/api/suspend", makeHTTPHandle(s.handleSuspend))
	http.HandleFunc("/api/retrievefornotifications", makeHTTPHandle(s.handleRetrieveNotifications))
	return http.ListenAndServe(s.port, nil)
}
