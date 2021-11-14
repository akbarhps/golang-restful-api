package repository

import (
	"context"
	"database/sql"
	"go-api/domain"
	"go-api/helper"
)

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (*userRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user *domain.User) {
	query := "INSERT INTO users(id, full_name, username, email, password, created_at) VALUES(?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, user.Id, user.DisplayName, user.Username, user.Email, user.Password, user.CreatedAt)
	helper.PanicIfError(err)
}

func (*userRepositoryImpl) Find(ctx context.Context, tx *sql.Tx, user *domain.User) []domain.User {
	query := "SELECT id, full_name, username, email, password, created_at FROM users WHERE id = ? OR email = ? OR username = ? OR full_name LIKE ?"
	rows, err := tx.QueryContext(ctx, query, user.Id, user.Email, user.Username, user.DisplayName)
	helper.PanicIfError(err)
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var rowUser domain.User

		err = rows.Scan(&rowUser.Id, &rowUser.DisplayName, &rowUser.Username, &rowUser.Email, &rowUser.Password, &rowUser.CreatedAt)
		helper.PanicIfError(err)

		users = append(users, rowUser)
	}

	return users
}

func (*userRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *domain.User) {
	query := "UPDATE users SET full_name = ?, username = ?, email = ?, password = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, user.DisplayName, user.Username, user.Email, user.Password, user.Id)
	helper.PanicIfError(err)
}

func (*userRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user *domain.User) {
	query := "DELETE FROM users WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, user.Id)
	helper.PanicIfError(err)
}

func (*userRepositoryImpl) DeleteAll(ctx context.Context, tx *sql.Tx) {
	query := "DELETE FROM users"
	_, err := tx.ExecContext(ctx, query)
	helper.PanicIfError(err)
}
