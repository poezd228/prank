package app

import (
	"net/http"

	"github.com/mvd-inc/anibliss/internal/dependencies"
	"github.com/mvd-inc/anibliss/internal/service/cron"
	"github.com/mvd-inc/anibliss/pkg/logger"
)

type App interface {
	Start() error
	Stop()
	GetHandler() http.Handler
}

type app struct {
	deps   *dependencies.Dependencies
	logger logger.Logger
	cron   cron.Cron
}

func NewApp(
	deps *dependencies.Dependencies,
	logger logger.Logger,
) App {
	return &app{
		deps:   deps,
		logger: logger,
	}
}

func (a *app) Start() error {
	a.cron = a.deps.Cron()

	a.deps.Start()
	a.cron.RunJobs()

	return nil
}

func (a *app) Stop() {

	if a.cron != nil {
		a.cron.Stop()
	}

	a.deps.Stop()
}

func (a *app) GetHandler() http.Handler {
	mux := http.NewServeMux()

	a.deps.Handler().FillHandlers(mux)

	return a.deps.Middleware().RateLimitMiddleware(mux)
}
