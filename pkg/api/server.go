package api

import (
	"Golang-API-Assessment/pkg/repository"
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
	http.HandleFunc("/api/register", MakeHTTPHandle(s.HandleRegister))
	http.HandleFunc("/api/commonstudents", MakeHTTPHandle(s.HandleCommonStudents))
	http.HandleFunc("/api/suspend", MakeHTTPHandle(s.HandleSuspend))
	http.HandleFunc("/api/retrievefornotifications", MakeHTTPHandle(s.HandleRetrieveNotifications))
	http.HandleFunc("/api/populatestudentsandteachers", MakeHTTPHandle(s.HandlePopulateStudentsAndTeachers))
	http.HandleFunc("/api/cleardatabase", MakeHTTPHandle(s.HandleClearDatabase))
	return http.ListenAndServe(s.listenAddr, nil)
}
