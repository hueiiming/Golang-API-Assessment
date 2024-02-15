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
	http.HandleFunc("/api/register", makeHTTPHandle(s.handleRegister))
	http.HandleFunc("/api/commonstudents", makeHTTPHandle(s.handleCommonStudents))
	http.HandleFunc("/api/suspend", makeHTTPHandle(s.handleSuspend))
	http.HandleFunc("/api/retrievefornotifications", makeHTTPHandle(s.handleRetrieveNotifications))
	return http.ListenAndServe(s.listenAddr, nil)
}
