package repository

import (
	"Golang-API-Assessment/types"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestPostgreSQLRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQLRepository{
		db: db,
	}

	// Registration SQL query was prepared once then executed twice in the loop
	mock.ExpectPrepare("INSERT INTO").ExpectExec().
		WithArgs("teacherken@gmail.com", "studentjon@gmail.com").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO").
		WithArgs("teacherken@gmail.com", "studenthon@gmail.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	request := types.RegisterRequest{
		Teacher: "teacherken@gmail.com",
		Students: []string{
			"studentjon@gmail.com",
			"studenthon@gmail.com",
		},
	}

	err = repo.Registration(request)
	if err != nil {
		t.Errorf("Failed to insert into REGISTRATION TABLE: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed expectations: %s", err)
	}
}
