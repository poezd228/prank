package jwt

import (
	"context"
	"fmt"

	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
)

type Repository interface {
	//DropAllTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, id int64) error
	DropTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, id int64) error
	FindNumberTX(ctx context.Context, tx transactions.Transaction, id int64) (int64, error)
	AddTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token domain.Token) (domain.Token, error)
	CheckTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token domain.Token) (domain.Token, error)
	DropOldTokens(ctx context.Context, tx transactions.Transaction, timestamp int64) error
}
type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) DropOldTokens(ctx context.Context, tx transactions.Transaction, timestamp int64) error {

	query := `DELETE FROM tokens WHERE expires_at < $1`
	_, err := tx.Txm().Exec(ctx, query, timestamp)
	if err != nil {
		return err
	}
	return nil

}

func (r *repository) DropTokensTX(ctx context.Context, tx transactions.Transaction, role domain.Role, id int64) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := tx.Txm().Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil

}
func (r *repository) AddTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token domain.Token) (domain.Token, error) {
	query := `
		INSERT INTO tokens (id, number, purpose, secret, expires_at)
		VALUES($1, $2, $3, $4, $5)
	`

	res, err := tx.Txm().Exec(ctx, query, token.Id, token.Number, token.Purpose, token.Secret, token.ExpiresAt)
	if err != nil {
		return domain.Token{}, err
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected != 1 {
		return domain.Token{}, err
	}

	return token, err
}

func (r *repository) CheckTokenTX(ctx context.Context, tx transactions.Transaction, role domain.Role, token domain.Token) (domain.Token, error) {
	var token1 domain.Token
	query := `
		SELECT id, number, purpose, secret, expires_at
		FROM tokens
		WHERE id=$1
			AND number= $2
			AND purpose= $3
			AND secret= $4
	`
	rows, err := tx.Txm().Query(ctx, query, token.Id, token.Number, token.Purpose, token.Secret)

	if err != nil {
		return domain.Token{}, err
	}

	if !rows.Next() {
		return domain.Token{}, errors.TokenDoesNotExist
	}

	err = rows.Scan(&token1.Id, &token1.Number, &token1.Purpose, &token1.Secret, &token1.ExpiresAt)

	if err != nil {

		return domain.Token{}, errors.TokenDoesNotExist
	}

	return token1, nil
}

func (r *repository) FindNumberTX(ctx context.Context, tx transactions.Transaction, id int64) (int64, error) {
	var numbers []int64
	query := `
		SELECT number 
		FROM tokens 
		WHERE id=$1 AND purpose=0 
		ORDER BY number
	`
	rows, err := tx.Txm().Query(ctx, query, id)

	for rows.Next() {
		var number int64
		if err := rows.Scan(&number); err != nil {
			return 0, err
		}
		numbers = append(numbers, number)
	}
	//err := res.Scan(numbers)
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return r.findNumbers(numbers)
}

func (r *repository) findNumbers(numbers []int64) (int64, error) {

	if len(numbers) == 0 {
		return 0, nil
	}
	if numbers[len(numbers)-1] == int64(len(numbers)-1) {
		return int64(len(numbers)), nil
	}
	for i, n := range numbers {
		if n != int64(i) {
			return int64(i), nil
		}
	}

	return 0, errors.ServiceErr{Err: nil}
}
