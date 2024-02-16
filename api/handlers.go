package api

import (
	"Golang-API-Assessment/types"
	"Golang-API-Assessment/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Message string `json:"message"`
}

func makeHTTPHandle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteToJSON(w, http.StatusBadRequest, ApiError{Message: "error: " + err.Error()})
		}
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("status method not allowed")
	}

	registerReq := types.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	if isEmailValid, err := utils.IsValidEmail(registerReq.Teacher); err != nil || !isEmailValid {
		return fmt.Errorf("invalid teacher email: %w", err)
	}

	if areEmailsValid, err := utils.AreValidEmails(registerReq.Students); err != nil || !areEmailsValid {
		return fmt.Errorf("invalid student email: %w", err)
	}

	if err := s.repo.Registration(&registerReq); err != nil {
		return fmt.Errorf("error registering: %w", err)
	}

	return WriteToJSON(w, http.StatusNoContent, registerReq)
}

func (s *Server) handleCommonStudents(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("status method not allowed")
	}

	queryParam := r.URL.Query()

	if err := hasWrongParam(queryParam); err != nil {
		return err
	}

	teachers := queryParam["teacher"]

	students, err := s.repo.GetCommonStudents(teachers)
	if err != nil {
		return fmt.Errorf("error getting common students: %w", err)
	}

	if len(students) == 0 {
		students = []string{}
	}

	if len(teachers) == 1 && len(students) >= 1 {
		students = append(students, "student_only_under_"+teachers[0])
	}

	commonStudents := &types.CommonStudentsResponse{
		Students: students,
	}

	return WriteToJSON(w, http.StatusOK, commonStudents)
}

func (s *Server) handleSuspend(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("status method not allowed")
	}

	suspendReq := types.SuspendRequest{}
	if err := json.NewDecoder(r.Body).Decode(&suspendReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	if suspendReq.Student == "" {
		return fmt.Errorf("invalid request")
	}

	if err := s.repo.Suspension(&suspendReq); err != nil {
		return fmt.Errorf("error suspending: %w", err)
	}

	return WriteToJSON(w, http.StatusNoContent, nil)
}

func (s *Server) handleRetrieveNotifications(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("status method not allowed")
	}

	notifReq := types.NotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&notifReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	if notifReq.Teacher == "" || notifReq.Message == "" {
		return fmt.Errorf("invalid request")
	}

	if isEmailValid, err := utils.IsValidEmail(notifReq.Teacher); err != nil || !isEmailValid {
		return fmt.Errorf("invalid teacher email: %w", err)
	}

	recipients, err := s.repo.GetNotification(&notifReq)
	if err != nil {
		return fmt.Errorf("error getting notification: %w", err)
	}

	if len(recipients) == 0 {
		recipients = []string{}
	}

	notification := &types.NotificationResponse{
		Recipients: recipients,
	}

	return WriteToJSON(w, http.StatusOK, notification)
}

func WriteToJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func removeEmptyStrings(s []string) []string {
	var result []string
	for _, str := range s {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func hasWrongParam(queryParam map[string][]string) error {
	for paramName := range queryParam {
		if paramName != "teacher" {
			return fmt.Errorf("invalid query param")
		}
	}
	for _, emails := range queryParam {
		for _, email := range emails {
			if email == "" {
				return fmt.Errorf("empty query param")
			}
		}
	}

	return nil
}
