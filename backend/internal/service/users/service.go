package users

import (
	"context"

	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
	"github.com/mvd-inc/anibliss/internal/repository/users"
	"github.com/mvd-inc/anibliss/pkg/utils"
)

type Service interface {
	GetUserLP(ctx context.Context, login string, password string) (domain.Account, errors.ServiceError)
	GetUser(ctx context.Context, account domain.Account) (domain.Account, errors.ServiceError)
	ChangeUserPass(ctx context.Context, oldPass, newPass, login string) errors.ServiceError
	CreateUser(ctx context.Context, login string, password string) errors.ServiceError
}
type service struct {
	cfg             config.Config
	transactionRepo transactions.Repository
	usersRepo       users.Repository
}

func NewService(cfg config.Config, transactionRepo transactions.Repository, usersRepo users.Repository) Service {
	return &service{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		usersRepo:       usersRepo,
	}

}

func (s *service) CreateUser(ctx context.Context, login string, password string) errors.ServiceError {
	//e := validators.ValidateUserPasswordAndLogin(login, password)
	//if e != nil{
	//	return e
	//}
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return errors.DatabaseError(err)
	}
	defer tx.Rollback(context.Background())
	registered, err := s.usersRepo.CheckUser(ctx, tx, login)
	if err != nil {
		return errors.CreateUserErr(err)
	}
	if !registered {
		return errors.AlreadyRegistered()
	}
	err = s.usersRepo.CreateUser(ctx, tx, login, utils.HashSha256(password))
	if err != nil {
		return errors.CreateUserErr(err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return errors.DatabaseError(err)
	}
	return nil

}
func (s *service) ChangeUserPass(ctx context.Context, oldPass, newPass, login string) errors.ServiceError {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return errors.DatabaseError(err)
	}
	err = s.usersRepo.ChangePassword(ctx, tx, oldPass, newPass, login)
	if err != nil {
		return errors.WrongCredentials()
	}
	return nil

}

func (s *service) GetUser(ctx context.Context, account domain.Account) (domain.Account, errors.ServiceError) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return domain.Account{}, errors.ParseFailed(err)
	}
	defer tx.Rollback(context.Background())
	acc, err := s.usersRepo.GetUserById(ctx, tx, account)
	if err != nil {
		return domain.Account{}, errors.UserNotFoundById()
	}

	return acc, nil

}

func (s *service) GetUserLP(ctx context.Context, login string, password string) (domain.Account, errors.ServiceError) {
	tx, err := s.transactionRepo.StartTransaction(ctx)
	if err != nil {
		return domain.Account{}, errors.ParseFailed(err)
	}
	defer tx.Rollback(context.Background())
	acc, err := s.usersRepo.GetUserByLogPass(ctx, tx, login, utils.HashSha256(password))
	if err != nil {
		return domain.Account{}, errors.WrongCredentials()
	}

	return acc, nil

}
