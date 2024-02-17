package repository

import "Golang-API-Assessment/pkg/types"

//go:generate mockery --name=Repository
type Repository interface {
	Registration(teacherID int, studentID []int) error
	GetCommonStudents(teachers []string) ([]string, error)
	Suspension(studentID int) error
	GetNotification(request *types.NotificationRequest) ([]string, error)
	GetTeacherID(teacherEmail string) (int, error)
	GetStudentID(studentEmail string) (int, error)
	PopulateTables(teacherEmails []string, studentEmails []string) error
	ClearTables() error
}
