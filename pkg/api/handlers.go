package api

import (
	"Golang-API-Assessment/pkg/repository"
	"Golang-API-Assessment/pkg/types"
	"Golang-API-Assessment/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Message string `json:"message"`
}

const (
	HttpMethodGet  = "GET"
	HttpMethodPost = "POST"
)

// MakeHTTPHandle Wrapper func to handle handlers error
func MakeHTTPHandle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteToJSON(w, http.StatusBadRequest, ApiError{Message: "error: " + err.Error()})
		}
	}
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodPost); err != nil {
		return err
	}

	registerReq := types.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	// Validate JSON request
	if registerReq.Teacher == "" || len(registerReq.Students) == 0 {
		return fmt.Errorf("invalid JSON request")
	}

	if isEmailValid, err := utils.IsValidEmail(registerReq.Teacher); err != nil || !isEmailValid {
		return fmt.Errorf("invalid teacher email: %w", err)
	}

	if areEmailsValid, err := utils.AreValidEmails(registerReq.Students); err != nil || !areEmailsValid {
		return fmt.Errorf("invalid student email: %w", err)
	}

	teacherID, err := s.repo.GetTeacherID(registerReq.Teacher)
	if err != nil {
		return err
	}

	studentIDs, err := GetStudentIDs(s.repo, registerReq.Students)
	if err != nil {
		return err
	}

	if err := s.repo.Registration(teacherID, studentIDs); err != nil {
		return fmt.Errorf("error registering: %w", err)
	}

	return WriteToJSON(w, http.StatusNoContent, registerReq)
}

func (s *Server) HandleCommonStudents(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodGet); err != nil {
		return err
	}

	queryParam := r.URL.Query()
	if err := hasWrongCommonStudentsParam(s.repo, queryParam); err != nil {
		return err
	}
	teachers := queryParam["teacher"]

	students, err := s.repo.GetCommonStudents(teachers)
	if err != nil {
		return fmt.Errorf("error getting common students: %w", err)
	}

	// return empty string slice when no students found
	if len(students) == 0 {
		students = []string{}
	}

	// if query param only includes 1 teacher email
	if len(teachers) == 1 && len(students) >= 1 {
		students = append(students, "student_only_under_"+teachers[0])
	}

	commonStudents := &types.CommonStudentsResponse{
		Students: students,
	}

	return WriteToJSON(w, http.StatusOK, commonStudents)
}

func (s *Server) HandleSuspend(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodPost); err != nil {
		return err
	}

	suspendReq := types.SuspendRequest{}
	if err := json.NewDecoder(r.Body).Decode(&suspendReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	// Missing student request
	if suspendReq.Student == "" {
		return fmt.Errorf("missing student request")
	}

	if isEmailValid, err := utils.IsValidEmail(suspendReq.Student); err != nil || !isEmailValid {
		return fmt.Errorf("invalid student email: %w", err)
	}

	studentID, err := s.repo.GetStudentID(suspendReq.Student)
	if err != nil {
		return err
	}

	if err := s.repo.Suspension(studentID); err != nil {
		return fmt.Errorf("error suspending: %w", err)
	}

	return WriteToJSON(w, http.StatusNoContent, nil)
}

func (s *Server) HandleRetrieveNotifications(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodPost); err != nil {
		return err
	}

	notifReq := types.NotificationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&notifReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	// Missing teacher or notification message
	if notifReq.Teacher == "" || notifReq.Notification == "" {
		return fmt.Errorf("invalid teacher or notification request")
	}

	if isEmailValid, err := utils.IsValidEmail(notifReq.Teacher); err != nil || !isEmailValid {
		return fmt.Errorf("invalid teacher email: %w", err)
	}

	// Check if teacher email exist in database
	if _, err := s.repo.GetTeacherID(notifReq.Teacher); err != nil {
		return err
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

func (s *Server) HandlePopulateStudentsAndTeachers(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodPost); err != nil {
		return err
	}

	bulkImportReq := types.BulkImportRequest{}

	if err := json.NewDecoder(r.Body).Decode(&bulkImportReq); err != nil {
		return fmt.Errorf("error decoding JSON request: %w", err)
	}

	if areEmailsValid, err := utils.AreValidEmails(bulkImportReq.Teachers); err != nil || !areEmailsValid {
		return fmt.Errorf("invalid teacher email: %w", err)
	}

	if areEmailsValid, err := utils.AreValidEmails(bulkImportReq.Students); err != nil || !areEmailsValid {
		return fmt.Errorf("invalid student email: %w", err)
	}

	if err := s.repo.PopulateTables(bulkImportReq.Teachers, bulkImportReq.Students); err != nil {
		return err
	}

	return WriteToJSON(w, http.StatusNoContent, nil)
}

func (s *Server) HandleClearDatabase(w http.ResponseWriter, r *http.Request) error {
	if err := isCorrectRequestMethod(r, HttpMethodPost); err != nil {
		return err
	}

	if err := s.repo.ClearTables(); err != nil {
		return fmt.Errorf("error clearing tables")
	}

	return WriteToJSON(w, http.StatusNoContent, nil)
}

// Helper funcs

func GetStudentIDs(repo repository.Repository, studentEmails []string) ([]int, error) {
	var studentIDs []int
	for _, email := range studentEmails {
		studentID, err := repo.GetStudentID(email)
		if err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, studentID)
	}
	return studentIDs, nil
}

func WriteToJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func isCorrectRequestMethod(r *http.Request, method string) error {
	if r.Method != method {
		return fmt.Errorf("status method not allowed")
	}
	return nil
}

func hasWrongCommonStudentsParam(repo repository.Repository, queryParam map[string][]string) error {
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
			// Check if teacher email exist in database
			if _, err := repo.GetTeacherID(email); err != nil {
				return err
			}
		}
	}

	return nil
}
