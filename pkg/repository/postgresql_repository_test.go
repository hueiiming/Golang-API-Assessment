package repository

import (
	"Golang-API-Assessment/pkg/types"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"testing"
)

func TestPostgreSQLRepository_Registration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	// Registration SQL query was prepared once then executed twice in the loop
	mock.ExpectPrepare("INSERT INTO REGISTRATION").ExpectExec().
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO REGISTRATION").
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	teacherID := 1
	studentIDs := []int{1, 2}
	err = repo.Registration(teacherID, studentIDs)
	if err != nil {
		t.Errorf("Failed to insert into REGISTRATION TABLE: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}

func TestPostgreSQLRepository_GetCommonStudents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	// Register student to teacher
	mock.ExpectPrepare("INSERT INTO REGISTRATION").ExpectExec().
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	teacherID := 1
	studentIDs := []int{2}
	err = repo.Registration(teacherID, studentIDs)

	if err != nil {
		t.Errorf("Failed to insert into REGISTRATION TABLE: %s", err)
	}

	// Get common students
	rows := sqlmock.NewRows([]string{"student_email"}).
		AddRow("studentjon@gmail.com")

	pqTeachers := pq.StringArray([]string{"teacherken@gmail.com"})

	mock.ExpectPrepare("SELECT s.student_email FROM REGISTRATION r").ExpectQuery().
		WithArgs(pqTeachers).
		WillReturnRows(rows)

	request := []string{"teacherken@gmail.com"}

	_, err = repo.GetCommonStudents(request)
	if err != nil {
		t.Errorf("Failed to get student_email from REGISTRATION TABLE: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}

func TestPostgreSQLRepository_Suspension(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	mock.ExpectPrepare("INSERT INTO SUSPENSION").ExpectExec().
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	studentID := 3

	err = repo.Suspension(studentID)
	if err != nil {
		t.Errorf("Failed to insert into SUSPENSION TABLE: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}

func TestPostgreSQLRepository_GetNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	// Register student to teacher
	mock.ExpectPrepare("INSERT INTO REGISTRATION").ExpectExec().
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	teacherID := 1
	studentIDs := []int{2}
	err = repo.Registration(teacherID, studentIDs)

	if err != nil {
		t.Errorf("Failed to insert into REGISTRATION TABLE: %s", err)
	}

	// Suspend student
	mock.ExpectPrepare("INSERT INTO SUSPENSION").ExpectExec().
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	studentID := 3

	err = repo.Suspension(studentID)
	if err != nil {
		t.Errorf("Failed to insert into SUSPENSION TABLE: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}

	// Get notifications
	rows := sqlmock.NewRows([]string{"student_email"}).
		AddRow("studentjon@gmail.com")
	pqEmails := pq.StringArray([]string{"studentagnes@gmail.com", "studentmiche@gmail.com"})

	mock.ExpectPrepare("SELECT DISTINCT s.student_email FROM STUDENT s LEFT JOIN REGISTRATION r").ExpectQuery().
		WithArgs("teacherken@gmail.com", pqEmails).
		WillReturnRows(rows)

	notifRequest := &types.NotificationRequest{
		Teacher:      "teacherken@gmail.com",
		Notification: "Hello students! @studentagnes@gmail.com @studentmiche@gmail.com",
	}

	_, err = repo.GetNotification(notifRequest)
	if err != nil {
		t.Errorf("Failed to get student_email from REGISTRATION AND SUSPENSION table: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}

func TestPostgreSQLRepository_GetStudentID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	rows := sqlmock.NewRows([]string{"student_id"}).
		AddRow(1)

	studentEmail := "studentjon@gmail.com"
	mock.ExpectQuery("^SELECT student_id FROM STUDENT*").
		WithArgs(studentEmail).
		WillReturnRows(rows)

	_, err = repo.GetStudentID(studentEmail)

	if err != nil {
		t.Errorf("Failed to get student_id from STUDENT table: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}

func TestPostgreSQLRepository_GetTeacherID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		Db: db,
	}

	rows := sqlmock.NewRows([]string{"teacher_id"}).
		AddRow(1)

	teacherEmail := "teacherken@gmail.com"
	mock.ExpectQuery("^SELECT teacher_id FROM TEACHER*").
		WithArgs(teacherEmail).
		WillReturnRows(rows)

	_, err = repo.GetTeacherID(teacherEmail)

	if err != nil {
		t.Errorf("Failed to get teacher_id from TEACHER table: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}
