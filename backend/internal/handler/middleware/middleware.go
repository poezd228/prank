package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/domain"
	"github.com/mvd-inc/anibliss/internal/errors"
	"github.com/mvd-inc/anibliss/internal/handler/writers"
	"github.com/mvd-inc/anibliss/internal/service/auth"
	"github.com/mvd-inc/anibliss/internal/service/users"
	"github.com/mvd-inc/anibliss/pkg/logger"
)

type Middleware interface {
	AuthMiddleware(next func(acc domain.Account, w http.ResponseWriter, r *http.Request)) http.Handler
	RateLimitMiddleware(handler http.Handler) http.Handler
	RefreshMiddleware(next func(acc domain.Account, w http.ResponseWriter, r *http.Request)) http.Handler
}

type middleware struct {
	cfg         config.HandlerConfig
	logger      logger.Logger
	q           chan any
	userService users.Service
	authService auth.Service
}

func NewMiddleware(
	cfg config.HandlerConfig,
	authService auth.Service,
	logger logger.Logger,
	queueSize int,
	usersService users.Service,

) Middleware {
	return &middleware{
		cfg:         cfg,
		authService: authService,
		logger:      logger,
		q:           make(chan any, queueSize),
		userService: usersService,
	}
}

const (
	cookieName        = "access-token"
	refreshCookieName = "refresh-token"
)

func (m *middleware) AuthMiddleware(next func(acc domain.Account, w http.ResponseWriter, r *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var (
			accessCookie *http.Cookie
			err          error
		)
		accessCookie, err = r.Cookie(cookieName)
		if err != nil {
			writers.WriteErrorResponseWithErrorLog(w, m.logger, errors.MissingCookies(err), "ru")
			return

		}
		token := accessCookie.Value

		acc, e := m.authService.JwtAuth(ctx, domain.PurposeAccess, token)

		if e != nil {
			m.logger.Error(e)
			writers.WriteErrorResponseWithErrorLog(w, m.logger, e, "ru")
			return
		}
		m.logger.Info(acc)
		account, err := m.userService.GetUser(ctx, acc)

		next(account, w, r)

	})

}

func (m *middleware) RefreshMiddleware(next func(acc domain.Account, w http.ResponseWriter, r *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var (
			refreshCookie *http.Cookie
			err           error
		)

		refreshCookie, err = r.Cookie(refreshCookieName)
		if err != nil {
			writers.WriteErrorResponseWithErrorLog(w, m.logger, errors.MissingCookies(err), "ru")
			return

		}
		token := refreshCookie.Value
		acc, e := m.authService.JwtAuth(ctx, domain.PurposeRefresh, token)
		if e != nil {
			writers.WriteErrorResponseWithErrorLog(w, m.logger, e, "ru")
			return
		}
		account, err := m.userService.GetUser(ctx, acc)
		next(account, w, r)

	})

}
func (m *middleware) releaseWorker() {
	<-m.q
}
func (m *middleware) RateLimitMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := m.acquireWorker(r.Context()); e != nil {
			fmt.Println("Rate_limit")
			return
		}
		defer m.releaseWorker()
		handler.ServeHTTP(w, r)
	})
}
func (m *middleware) acquireWorker(ctx context.Context) errors.ServiceError {
	select {
	case <-ctx.Done():
		return errors.ErrorWithCode(429, fmt.Errorf("too many request"), map[string]string{
			"en": "Server overloaded. Please, try again later",
			"ru": "Сервер перегружен. Пожалуйста, попробуйте позже",
		})
	case m.q <- struct{}{}:
		return nil
	}
}
