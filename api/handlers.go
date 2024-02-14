package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {}
func (s *Server) handleCommonStudents(w http.ResponseWriter, r *http.Request) {
	user := s.repo.GetCommonStudents()
	json.NewEncoder(w).Encode(user)
}
func (s *Server) handleSuspend(w http.ResponseWriter, r *http.Request)               {}
func (s *Server) handleRetrieveNotifications(w http.ResponseWriter, r *http.Request) {}
