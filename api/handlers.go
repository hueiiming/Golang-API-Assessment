package api

import (
	"Golang-API-Assessment/types"
	"encoding/json"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Message string `json:"message"`
}

func makeHTTPHandle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteToJSON(w, http.StatusBadRequest, ApiError{Message: err.Error()})
		}
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		registerReq := types.RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
			return err
		}

		if err := s.repo.Registration(&registerReq); err != nil {
			return err
		}

		return WriteToJSON(w, http.StatusNoContent, registerReq)
	}
	return nil
}
func (s *Server) handleCommonStudents(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		queryParam := r.URL.Query()
		teachers := queryParam["teacher"]

		var allStudents []string

		if len(teachers) == 1 {
			teacherEmail := teachers[0]
			students, err := s.repo.GetCommonStudents(teacherEmail)
			if err != nil {
				return err
			}
			students = append(students, "student_only_under_"+teacherEmail)
			allStudents = append(allStudents, students...)

		} else {
			var students []string

			for _, teacherEmail := range teachers {
				currStudent, err := s.repo.GetCommonStudents(teacherEmail)
				if err != nil {
					return err
				}
				students = append(students, currStudent...)
				allStudents = append(allStudents, students...)

			}
		}

		commonStudents := &types.CommonStudents{
			Students: allStudents,
		}

		return WriteToJSON(w, http.StatusOK, commonStudents)
	}
	return nil
}
func (s *Server) handleSuspend(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *Server) handleRetrieveNotifications(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteToJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
