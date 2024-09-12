package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	jwt2 "github.com/mvd-inc/anibliss/internal/repository/jwt"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/mvd-inc/anibliss/pkg/time_manager"
)

const (
	Alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
)

type (
	Service interface {
		CreateTokens(ctx context.Context, role domain.Role, id int64) (int64, domain.TokenInfo, domain.TokenInfo, errors.ServiceError)
		Auth(ctx context.Context, token string, purpose domain.AuthPurpose) (domain.Account, int64, errors.ServiceError)
		ReCreateTokens(ctx context.Context, role domain.Role, id int64) (int64, domain.TokenInfo, domain.TokenInfo, errors.ServiceError)
		DropOldTokens(ctx context.Context, timestamp int64) errors.ServiceError
	}
	service struct {
		tx transactions.Repository

		repo jwt2.Repository

		timeManager time_manager.TimeManager
		//randomAdapter random_util.Adapter

		l logger.Logger

		key                                                  string
		accessTokenTimeout, refreshTokenTimeout, authTimeout time.Duration
		lock                                                 *sync.RWMutex
	}
)

func NewService(tx transactions.Repository, repo jwt2.Repository, timeAdapter time_manager.TimeManager, l logger.Logger, cfg *config.AuthConfig) Service {
	return &service{
		tx:                  tx,
		repo:                repo,
		timeManager:         timeAdapter,
		l:                   l,
		key:                 cfg.Key,
		accessTokenTimeout:  cfg.AccessTokenTimeout,
		refreshTokenTimeout: cfg.RefreshTokenTimeout,
		authTimeout:         cfg.AuthTimeout,
		lock:                &sync.RWMutex{},
	}

}

func (s *service) Auth(
	ctx context.Context, token string, purpose domain.AuthPurpose,
) (domain.Account, int64, errors.ServiceError) {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx, e := s.tx.StartTransaction(ctx)
	if e != nil {
		return domain.Account{}, 0, errors.DatabaseError(e)
	}
	defer tx.Rollback(context.Background())

	t, err := s.parseToken(token)

	if err != nil {
		return domain.Account{}, 0, err
	}

	acc, number, err := s.checkToken(ctx, tx, t, purpose)
	if err != nil {
		return domain.Account{}, 0, errors.CheckTokenFailed(err)
	}
	e = tx.Commit(ctx)
	if e != nil {
		return domain.Account{}, 0, err
	}

	return acc, number, nil
}

func (s *service) DropOldTokens(ctx context.Context, timestamp int64) errors.ServiceError {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.tx.StartTransaction(ctx)
	if err != nil {
		return errors.DatabaseError(err)

	}
	defer tx.Rollback(context.Background())
	err = s.repo.DropOldTokens(ctx, tx, timestamp)
	if err != nil {
		return errors.DropOldTokensError(err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return errors.DatabaseError(err)
	}
	return nil
}

func (s *service) ReCreateTokens(ctx context.Context, role domain.Role, id int64) (int64, domain.TokenInfo, domain.TokenInfo, errors.ServiceError) {

	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.tx.StartTransaction(ctx)
	if err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}
	defer tx.Rollback(context.Background())
	err = s.repo.DropTokensTX(ctx, tx, role, id)
	if err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DropTokensError(err)

	}
	number, accessToken, refreshToken, e := s.createTokens(ctx, tx, role, id)
	if e != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, e
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}

	return number, accessToken, refreshToken, nil
}

func (s *service) CreateTokens(ctx context.Context, role domain.Role, id int64) (int64, domain.TokenInfo, domain.TokenInfo, errors.ServiceError) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tx, err := s.tx.StartTransaction(ctx)
	if err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}
	defer tx.Rollback(context.Background())

	number, accessToken, refreshToken, e := s.createTokens(ctx, tx, role, id)
	if e != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, e
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}

	return number, accessToken, refreshToken, nil
}

func (s *service) createTokens(
	ctx context.Context, tx transactions.Transaction, role domain.Role, id int64,
) (int64, domain.TokenInfo, domain.TokenInfo, errors.ServiceError) {
	now := s.timeManager.Now()

	tx, err := s.tx.StartTransaction(ctx)
	if err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}
	defer tx.Rollback(context.Background())

	number, err := s.repo.FindNumberTX(ctx, tx, id)
	if err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.FindNumberError(err)
	}
	accessExpiresAt, refreshExpiresAt := now.Add(s.accessTokenTimeout), now.Add(s.refreshTokenTimeout)

	accessTokenHash, e := s.generateTokenHash(ctx, tx, role, id, number, domain.PurposeAccess, accessExpiresAt)
	if e != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, e
	}

	refreshTokenHash, e := s.generateTokenHash(ctx, tx, role, id, number, domain.PurposeRefresh, refreshExpiresAt)
	if e != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, e
	}

	accessToken := domain.TokenInfo{
		Token:     accessTokenHash,
		ExpiresAt: accessExpiresAt.UnixNano() / 1e+6,
	}
	refreshToken := domain.TokenInfo{
		Token:     refreshTokenHash,
		ExpiresAt: refreshExpiresAt.UnixNano() / 1e+6,
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, domain.TokenInfo{}, domain.TokenInfo{}, errors.DatabaseError(err)
	}

	return number, accessToken, refreshToken, nil
}

func (s *service) generateTokenHash(
	ctx context.Context, tx transactions.Transaction, role domain.Role, id, number int64, purpose domain.AuthPurpose,
	expire time.Time,
) (string, errors.ServiceError) {
	secret := s.generateSecret(role, id, number, purpose)
	tokenToAdd := domain.Token{
		Id:        id,
		Number:    number,
		Purpose:   int(purpose),
		Secret:    secret,
		ExpiresAt: expire.UnixNano() / 1e+6,
	}
	if _, err := s.repo.AddTokenTX(ctx, tx, role, tokenToAdd); err != nil {
		return "", errors.TokenAddFailed(err)
	}
	claims := jwt.MapClaims{
		"id":      id,
		"role":    role,
		"purpose": purpose,
		"secret":  secret,
		"exp":     expire.Unix(),
		"number":  number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, e := token.SignedString([]byte(s.key))
	if e != nil {
		return "", errors.SignStringErr(e)
	}
	return res, nil
}

func (s *service) parseToken(token string) (*jwt.Token, errors.ServiceError) {
	res, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ParseTokenFailed(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(s.key), nil
	})
	if err != nil && res == nil {
		return nil, errors.ParseTokenFailed(err)
	}
	return res, nil
}

func (s *service) checkToken(
	ctx context.Context, tx transactions.Transaction, token *jwt.Token, purpose domain.AuthPurpose,
) (domain.Account, int64, errors.ServiceError) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return domain.Account{}, 0, errors.InvalidTokenClaims
	}

	if !claims.VerifyExpiresAt(s.timeManager.Now().Unix(), true) {
		return domain.Account{}, 0, errors.AuthExpiredToken
	}

	if realPurpose, ok := claims["purpose"].(float64); !ok {

		return domain.Account{}, 0, errors.AuthInvalidToken
	} else if domain.AuthPurpose(realPurpose) != purpose {

		return domain.Account{}, 0, errors.AuthInvalidTokenPurpose

	}

	if !token.Valid {

		return domain.Account{}, 0, errors.AuthInvalidToken
	}

	id, err := s.parseTokenIntClaim(claims, "id")
	if err != nil {
		return domain.Account{}, 0, err
	}

	number, err := s.parseTokenIntClaim(claims, "number")
	if err != nil {

		return domain.Account{}, 0, err
	}

	role, err := s.parseTokenStringClaim(claims, "role")
	if err != nil {

		return domain.Account{}, 0, err
	}

	secret, err := s.parseTokenStringClaim(claims, "secret")
	if err != nil {

		return domain.Account{}, 0, err
	}

	tokenModel := domain.Token{
		Id:      id,
		Number:  number,
		Purpose: int(purpose),
		Secret:  secret,
	}
	s.l.Info(role)
	s.l.Info(tokenModel)

	if _, e := s.repo.CheckTokenTX(ctx, tx, domain.Role(role), tokenModel); e != nil {
		return domain.Account{}, 0, errors.AuthInvalidToken
	}

	return domain.Account{
		Role: role,
		Id:   id,
	}, number, nil
}

func (s *service) parseTokenIntClaim(claims jwt.MapClaims, key string) (int64, errors.ServiceError) {
	if parsedValue, ok := claims[key].(float64); !ok {
		return 0, errors.AuthInvalidToken
	} else {
		return int64(parsedValue), nil
	}
}
func (s *service) generateSecret(role domain.Role, id, number int64, purpose domain.AuthPurpose) string {
	toHashElems := []string{
		fmt.Sprintf("%s", role),
		fmt.Sprintf("%d", id),
		fmt.Sprintf("%d", number),
		fmt.Sprintf("%d", purpose),
		fmt.Sprintf("%d", s.timeManager.Now().UnixNano()),
		RandomString(20),
	}
	toHash := strings.Join(toHashElems, "_")
	hash := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hash[:])
}
func (s *service) parseTokenStringClaim(claims jwt.MapClaims, key string) (string, errors.ServiceError) {
	if stringValue, ok := claims[key].(string); !ok {
		return "", errors.AuthInvalidToken
	} else {
		return stringValue, nil
	}
}

func RandomString(length int) string {
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return string(res)
}
