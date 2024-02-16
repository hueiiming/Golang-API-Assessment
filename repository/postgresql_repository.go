package repository

import (
	"Golang-API-Assessment/types"
	"Golang-API-Assessment/utils"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"os"
)

//go:generate mockery --name=Repository
type Repository interface {
	Registration(request *types.RegisterRequest) error
	GetCommonStudents(teachers []string) ([]string, error)
	Suspension(request *types.SuspendRequest) error
	GetNotification(request *types.NotificationRequest) ([]string, error)
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository() (*PostgreSQLRepository, error) {
	var connStr string

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	switch env {
	case "local":
		// Connect to local postgresql
		connStr = "user=postgres dbname=postgres password=root sslmode=disable"
	case "prod":
		// Connect to deployed postgresql https://supabase.com
		password := os.Getenv("DB_PASSWORD")
		connStr = "user=postgres.wxmkhkkcxatyzukbfqtw password=" + password + " host=aws-0-ap-southeast-1.pooler.supabase.com port=5432 dbname=postgres"
	}

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

func (r *PostgreSQLRepository) GetCommonStudents(teachers []string) ([]string, error) {
	pqTeachers := pq.StringArray(teachers)
	query := "SELECT student_email FROM REGISTRATIONS WHERE teacher_email = any($1)"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		fmt.Errorf("error preparing statement: %s", err)
		return nil, err
	}

	rows, err := stmt.Query(pqTeachers)
	if err != nil {
		fmt.Errorf("error querying from DB: %s", err)
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

func (r *PostgreSQLRepository) Suspension(request *types.SuspendRequest) error {
	query := "INSERT INTO suspensions (student_email) VALUES ($1)"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(request.Student)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLRepository) GetNotification(request *types.NotificationRequest) ([]string, error) {
	emails, err := utils.ExtractEmails(request.Message)
	if err != nil {
		return nil, err
	}

	pqEmails := pq.StringArray(emails)

	query := `SELECT DISTINCT student_email
				FROM REGISTRATIONS
				WHERE (teacher_email = $1 OR student_email = any($2))
				AND student_email NOT IN (
				    SELECT student_email FROM SUSPENSIONS
				)
				`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		fmt.Errorf("error preparing statement: %s", err)
		return nil, err
	}

	rows, err := stmt.Query(request.Teacher, pqEmails)
	if err != nil {
		fmt.Errorf("error querying from DB: %s", err)
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
	`
	_, err := r.db.Exec(query)
	return err
}
