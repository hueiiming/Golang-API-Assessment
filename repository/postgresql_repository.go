package repository

import "Golang-API-Assessment/types"

type Repository interface {
	GetCommonStudents() *types.CommonStudents
	GetNotification() *types.Notification
}

type PostgreSQLRepository struct{}

func NewPostgreSQLRepository() *PostgreSQLRepository {
	return &PostgreSQLRepository{}
}

func (r *PostgreSQLRepository) GetCommonStudents() *types.CommonStudents {
	return &types.CommonStudents{
		Students: []string{
			"test@gmail.com",
			"test2@gmail.com",
		},
	}
}

func (r *PostgreSQLRepository) GetNotification() *types.Notification {
	return &types.Notification{}
}
