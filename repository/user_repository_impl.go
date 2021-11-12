package repository

import (
	"context"
	"database/sql"
	"go-api/domain"
)

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (repository *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users(id, full_name, username, email, password, created_at) VALUES(?, ?, ?, ?, ?, ?)"
	_, err := repository.DB.ExecContext(ctx, query, user.Id, user.FullName, user.Username, user.Email, user.Password, user.CreatedAt)
	return err
}

func (repository *userRepositoryImpl) Find(ctx context.Context, user *domain.User) ([]domain.User, error) {
	query := "SELECT id, full_name, username, email, password, created_at FROM users WHERE id = ? OR email = ? OR username = ? OR full_name LIKE ?"
	rows, err := repository.DB.QueryContext(ctx, query, user.Id, user.Email, user.Username, user.FullName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var rowUser domain.User
		err = rows.Scan(&rowUser.Id, &rowUser.FullName, &rowUser.Username, &rowUser.Email, &rowUser.Password, &rowUser.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, rowUser)
	}

	return users, nil
}

func (repository *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	query := "UPDATE users SET full_name = ?, username = ?, email = ?, password = ? WHERE id = ?"
	_, err := repository.DB.ExecContext(ctx, query, user.FullName, user.Username, user.Email, user.Password, user.Id)
	return err
}

func (repository *userRepositoryImpl) Delete(ctx context.Context, user *domain.User) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := repository.DB.ExecContext(ctx, query, user.Id)
	return err
}

func (repository *userRepositoryImpl) DeleteAll(ctx context.Context) error {
	query := "DELETE FROM users"
	_, err := repository.DB.ExecContext(ctx, query)
	return err
}
