package repository

import (
	"Golang-API-Assessment/pkg/types"
	"Golang-API-Assessment/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"os"
)

//go:generate mockery --name=Repository
type Repository interface {
	Registration(teacherID int, studentID []int) error
	GetCommonStudents(teachers []string) ([]string, error)
	Suspension(studentID int) error
	GetNotification(request *types.NotificationRequest) ([]string, error)
	GetTeacherID(teacherEmail string) (int, error)
	GetStudentID(studentEmail string) (int, error)
	PopulateTables() error
	ClearTables() error
}

type PostgreSQLRepository struct {
	Db *sql.DB
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
		Db: db,
	}, nil
}

func (r *PostgreSQLRepository) Registration(teacherID int, studentIDs []int) error {
	query := "INSERT INTO REGISTRATION (teacher_id, student_id) VALUES ($1, $2)"
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return err
	}

	for _, studentID := range studentIDs {
		_, err := stmt.Exec(teacherID, studentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgreSQLRepository) GetCommonStudents(teachers []string) ([]string, error) {
	pqTeachers := pq.StringArray(teachers)
	query := `
		SELECT s.student_email 
		FROM REGISTRATION r
		INNER JOIN STUDENT s
		ON s.student_id = r.student_id
		INNER JOIN TEACHER t
		ON r.teacher_id = t.teacher_id
		WHERE t.teacher_email = any($1)`

	stmt, err := r.Db.Prepare(query)
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

func (r *PostgreSQLRepository) Suspension(studentID int) error {
	query := "INSERT INTO SUSPENSION (student_id) VALUES ($1)"

	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(studentID)
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

	query := `SELECT DISTINCT s.student_email
				FROM STUDENT s
				LEFT JOIN REGISTRATION r
				ON s.student_id = r.student_id
				LEFT JOIN TEACHER t
				ON t.teacher_id = r.teacher_id
				WHERE (t.teacher_email = $1 OR s.student_email = any($2))
				AND s.student_email NOT IN (
				    SELECT s.student_email 
				    FROM SUSPENSION sp
				    INNER JOIN STUDENT s
				    ON sp.student_id = s.student_id
				)
				`

	stmt, err := r.Db.Prepare(query)
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
		CREATE TABLE IF NOT EXISTS student (
		    student_id SERIAL PRIMARY KEY,
			student_email VARCHAR(255) UNIQUE
		);
		
		CREATE TABLE IF NOT EXISTS teacher (
			teacher_id SERIAL PRIMARY KEY,
			teacher_email VARCHAR(255) UNIQUE
		);
		
		CREATE TABLE IF NOT EXISTS registration (
			registration_id SERIAL PRIMARY KEY,
			teacher_id INT,
			student_id INT,
			FOREIGN KEY (teacher_id) REFERENCES teacher(teacher_id),
			FOREIGN KEY (student_id) REFERENCES student(student_id),
			UNIQUE (teacher_id, student_id)
		);
		
		CREATE TABLE IF NOT EXISTS suspension (
			suspension_id SERIAL PRIMARY KEY,
			student_id INT,
			FOREIGN KEY (student_id) REFERENCES student(student_id),
			UNIQUE (student_id)
		);
	`
	_, err := r.Db.Exec(query)
	return err
}

func (r *PostgreSQLRepository) GetStudentID(studentEmail string) (int, error) {
	query := "SELECT student_id FROM STUDENT WHERE student_email = $1"

	var studentID int
	err := r.Db.QueryRow(query, studentEmail).Scan(&studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("student with email %s not found", studentEmail)
		}
		return 0, fmt.Errorf("error querying student ID: %s", err)
	}

	return studentID, nil
}

func (r *PostgreSQLRepository) GetTeacherID(teacherEmail string) (int, error) {
	query := "SELECT teacher_id FROM TEACHER WHERE teacher_email = $1"

	var teacherID int
	err := r.Db.QueryRow(query, teacherEmail).Scan(&teacherID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("teacher with email %s not found", teacherEmail)
		}
		return 0, fmt.Errorf("error querying teacherID: %s", err)
	}

	return teacherID, nil
}

func (r *PostgreSQLRepository) PopulateTables() error {
	query := `
		INSERT INTO STUDENT (student_email)
		VALUES ('studentjon@gmail.com'),
		       ('studenthon@gmail.com'),
		       ('studentmay@gmail.com'),
		       ('studentagnes@gmail.com'),
		       ('studentmiche@gmail.com'),
		       ('studentbob@gmail.com'),
		       ('studentbad@gmail.com'),
		       ('studentmary@gmail.com');
		
		INSERT INTO TEACHER (teacher_email)
		VALUES ('teacherken@gmail.com'),
		       ('teacherjoe@gmail.com'),
		       ('teachermax@gmail.com');
	 `

	_, err := r.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert into tables: %w", err)
	}

	return nil
}

func (r *PostgreSQLRepository) ClearTables() error {
	query := `
		DELETE FROM REGISTRATION;
		DELETE FROM SUSPENSION;
		DELETE FROM STUDENT;
		DELETE FROM TEACHER;
		`

	_, err := r.Db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
