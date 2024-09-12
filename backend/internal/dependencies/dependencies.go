package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/db"
	"github.com/mvd-inc/anibliss/internal/handler"
	"github.com/mvd-inc/anibliss/internal/handler/middleware"
	"github.com/mvd-inc/anibliss/internal/repository/auth"
	jwt2 "github.com/mvd-inc/anibliss/internal/repository/jwt"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
	"github.com/mvd-inc/anibliss/internal/repository/users"
	auth2 "github.com/mvd-inc/anibliss/internal/service/auth"
	"github.com/mvd-inc/anibliss/internal/service/cron"
	"github.com/mvd-inc/anibliss/internal/service/jwt"
	users2 "github.com/mvd-inc/anibliss/internal/service/users"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/mvd-inc/anibliss/pkg/time_manager"
)

type Dependencies struct {
	cfg    *config.Config
	logger logger.Logger

	startFuncs []func()
	stopFuncs  []func()
	//db
	DBClient *db.PostgresClient
	//pkg

	timeManager time_manager.TimeManager
	//handler
	middleware  middleware.Middleware
	mainHandler handler.Handler

	//repo

	transactionRepo transactions.Repository
	usersRepo       users.Repository
	authRepo        auth.Repository
	jwtRepo         jwt2.Repository
	// service
	usersService users2.Service
	authService  auth2.Service
	jwtService   jwt.Service

	// cron
	cron cron.Cron
}

func NewDependencies(
	cfg *config.Config,
	logger logger.Logger,
	options ...Option,
) *Dependencies {

	res := &Dependencies{
		cfg:    cfg,
		logger: logger,
	}

	for _, o := range options {
		o.Apply(res)
	}

	return res
}

func (d *Dependencies) Start() {
	for _, f := range d.startFuncs {
		f()
	}
}

func (d *Dependencies) Stop() {
	for i := len(d.stopFuncs) - 1; i >= 0; i-- {
		d.stopFuncs[i]()
	}
}
