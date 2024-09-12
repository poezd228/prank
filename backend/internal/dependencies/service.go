package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/service/auth"
	"github.com/mvd-inc/anibliss/internal/service/jwt"
	"github.com/mvd-inc/anibliss/internal/service/users"
)

func (d *Dependencies) AuthService() auth.Service {
	if d.authService == nil {
		d.authService = auth.NewService(*d.cfg, d.TransactionRepo(), d.JwtService(), d.JwtRepo())
	}
	return d.authService

}

func (d *Dependencies) JwtService() jwt.Service {
	if d.jwtService == nil {
		d.jwtService = jwt.NewService(d.TransactionRepo(), d.JwtRepo(), d.TimeManager(), d.logger, d.cfg.Auth)
	}
	return d.jwtService

}

func (d *Dependencies) UsersService() users.Service {
	if d.usersService == nil {
		d.usersService = users.NewService(*d.cfg, d.TransactionRepo(), d.UsersRepo())
	}
	return d.usersService

}
