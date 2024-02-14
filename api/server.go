package api

import (
	"Golang-API-Assessment/repository"
	"net/http"
)

type Server struct {
	listenAddr string
	repo       repository.Repository
}

func NewServer(listenAddr string, repo repository.Repository) *Server {
	return &Server{
		listenAddr: listenAddr,
		repo:       repo,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/api/register", s.handleRegister)
	http.HandleFunc("/api/commonstudents", s.handleCommonStudents)
	http.HandleFunc("/api/suspend", s.handleSuspend)
	http.HandleFunc("/api/retrievefornotifications", s.handleRetrieveNotifications)
	return http.ListenAndServe(s.listenAddr, nil)
}
