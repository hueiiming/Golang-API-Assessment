package repository

import (
	"Golang-API-Assessment/types"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

//go:generate mockery --name=Repository
type Repository interface {
	Registration(request *types.RegisterRequest) error
	GetCommonStudents() (*types.CommonStudents, error)
	GetNotification() (*types.Notification, error)
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository() (*PostgreSQLRepository, error) {
	connStr := "user=postgres dbname=postgres password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQLRepository{
		db: db,
	}, nil
}

func (r *PostgreSQLRepository) Registration(request *types.RegisterRequest) error {
	query := "INSERT INTO registrations (teacher_email, student_email) VALUES ($1, $2)"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	for _, studentEmail := range request.Students {
		_, err := stmt.Exec(request.Teacher, studentEmail)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgreSQLRepository) GetCommonStudents(teachers []string) (*types.CommonStudents, error) {
	if len(teachers) == 1 {
		teacherEmail := teachers[0]
		students, err := r.getStudents(teacherEmail)
		if err != nil {
			return nil, err
		}
		students = append(students, "student_only_under_"+teacherEmail)

		return &types.CommonStudents{
			Students: students,
		}, nil

	} else {
		var allStudents []string

		for _, teacher := range teachers {
			students, err := r.getStudents(teacher)
			if err != nil {
				return nil, err
			}
			allStudents = append(allStudents, students...)
		}
		return &types.CommonStudents{
			Students: allStudents,
		}, nil
	}
}

func (r *PostgreSQLRepository) getStudents(teacherEmail string) ([]string, error) {
	query := "SELECT student_email FROM REGISTRATIONS WHERE teacher_email = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		fmt.Errorf("error preparing statement: %s", err)
		return nil, err
	}

	rows, err := stmt.Query(teacherEmail)
	if err != nil {
		fmt.Errorf("error getting teacherEmail: %s", err)
		return nil, err
	}

	var students []string

	for rows.Next() {
		var studentEmail string
		if err := rows.Scan(&studentEmail); err != nil {
			return nil, err
		}
		students = append(students, studentEmail)
	}

	return students, nil
}

func (r *PostgreSQLRepository) GetNotification() (*types.Notification, error) {
	return &types.Notification{}, nil
}

func (r *PostgreSQLRepository) Init() error {
	return r.createTables()
}

func (r *PostgreSQLRepository) createTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS registrations (
		   registration_id SERIAL PRIMARY KEY,
		   teacher_email VARCHAR(255),
		   student_email VARCHAR(255),
		   UNIQUE (teacher_email, student_email)
		);
		
		CREATE TABLE IF NOT EXISTS suspensions (
			 suspension_id SERIAL PRIMARY KEY,
			 student_email VARCHAR(255),
			 UNIQUE (student_email)
		);
		
		CREATE TABLE IF NOT EXISTS notifications (
		   notification_id SERIAL PRIMARY KEY,
		   teacher_email VARCHAR(255),
		   notification_text TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS mentioned_students (
			mention_id SERIAL PRIMARY KEY,
			notification_id INT REFERENCES notifications(notification_id),
			student_email VARCHAR(255)
		);
	`
	_, err := r.db.Exec(query)
	return err
}
