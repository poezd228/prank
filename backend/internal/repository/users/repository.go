package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
)

//go:generate mockgen -source repository.go -destination ../../mocks/users/users.go

type Repository interface {
	GetUserById(ctx context.Context, tx transactions.Transaction, acc domain.Account) (domain.Account, error)
	GetUserByLogPass(ctx context.Context, tx transactions.Transaction, login string, password string) (domain.Account, error)
	ChangePassword(ctx context.Context, tx transactions.Transaction, oldPass, newPass, login string) error
	CreateUser(ctx context.Context, tx transactions.Transaction, login string, password string) error
	CheckUser(ctx context.Context, tx transactions.Transaction, login string) (bool, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) ChangePassword(ctx context.Context, tx transactions.Transaction, oldPass, newPass, login string) error {
	query := `UPDATE users SET password = $1 WHERE login = $2 AND password = $3`
	row, err := tx.Txm().Exec(ctx, query, newPass, login, oldPass)

	if err != nil || row.RowsAffected() < 1 {
		return fmt.Errorf("Wrong login or password ")
	}
	return nil

}

func (r *repository) GetUserById(ctx context.Context, tx transactions.Transaction, acc domain.Account) (domain.Account, error) {
	account := domain.Account{}
	query := `SELECT id, login, role FROM users WHERE id = $1`
	row := tx.Txm().QueryRow(ctx, query, acc.Id)
	err := row.Scan(&account.Id, &account.Login, &account.Role)
	if err != nil {
		return domain.Account{}, err
	}

	return account, nil

}

func (r *repository) GetUserByLogPass(ctx context.Context, tx transactions.Transaction, login string, password string) (domain.Account, error) {
	acc := domain.Account{}
	query := `SELECT id, login, role FROM users WHERE login = $1 AND password = $2`
	row := tx.Txm().QueryRow(ctx, query, login, password)

	err := row.Scan(&acc.Id, &acc.Login, &acc.Role)
	if err != nil {
		return domain.Account{}, err
	}

	return acc, nil

}

func (r *repository) CheckUser(ctx context.Context, tx transactions.Transaction, login string) (bool, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE login=$1`
	err := tx.Txm().QueryRow(ctx, query, login).Scan(&count)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("users/CheckUser %w", err)
	}
	if count == 0 {
		return true, nil

	} else {
		return false, nil
	}

}
func (r *repository) CreateUser(ctx context.Context, tx transactions.Transaction, login string, password string) error {
	query := `INSERT INTO users (id, login, password, role) values (DEFAULT,$1,$2, $3)`
	_, err := tx.Txm().Exec(ctx, query, login, password, "user")
	if err != nil {
		return fmt.Errorf("users/CreateUser %w", err)
	}
	return nil

}
