package auth

import (
	"context"

	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	jwt2 "github.com/mvd-inc/anibliss/internal/repository/jwt"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
	"github.com/mvd-inc/anibliss/internal/service/jwt"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
)

type Service interface {
	JwtAuth(ctx context.Context, purpose domain.AuthPurpose,
		token string) (acc domain.Account, err errors.ServiceError)
	AccAuth(ctx context.Context, acc domain.Account) (domain.TokenResponse, errors.ServiceError)
	JwtRefresh(ctx context.Context, acc domain.Account) domain.TokenResponse
}
type service struct {
	cfg             config.Config
	transactionRepo transactions.Repository
	jwtService      jwt.Service
	jwtRepo         jwt2.Repository
}

func NewService(cfg config.Config, transactionRepo transactions.Repository, jwtService jwt.Service, jwtRepo jwt2.Repository) Service {
	return &service{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		jwtService:      jwtService,
		jwtRepo:         jwtRepo,
	}

}
func (s *service) JwtAuth(ctx context.Context,
	purpose domain.AuthPurpose,
	token string) (acc domain.Account, err errors.ServiceError) {

	account, _, e := s.jwtService.Auth(ctx, token, purpose)
	if e != nil {

		return domain.Account{}, e
	}

	return account, nil

}
func (s *service) AccAuth(ctx context.Context, acc domain.Account) (domain.TokenResponse, errors.ServiceError) {
	var role domain.Role
	role = domain.Role(acc.Role)
	_, accessToken, refreshToken, err := s.jwtService.CreateTokens(ctx, role, acc.Id)
	if err != nil {
		return domain.TokenResponse{}, err
	}
	response := domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil

}
func (s *service) JwtRefresh(ctx context.Context, acc domain.Account) domain.TokenResponse {
	var role domain.Role
	role = domain.Role(acc.Role)
	_, accessToken, refreshToken, _ := s.jwtService.ReCreateTokens(ctx, role, acc.Id)
	response := domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response

}
