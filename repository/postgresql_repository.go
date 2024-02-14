package repository

import (
	"Golang-API-Assessment/types"
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository interface {
	Registration(request types.RegisterRequest) error
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

func (r *PostgreSQLRepository) Registration(request types.RegisterRequest) error {
	query := `
		INSERT INTO registrations
		(teacher_email, student_email)
		VALUES ($1, $2)
`
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

func (r *PostgreSQLRepository) GetCommonStudents() (*types.CommonStudents, error) {
	return &types.CommonStudents{
		Students: []string{
			"test@gmail.com",
			"test2@gmail.com",
		},
	}, nil
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
