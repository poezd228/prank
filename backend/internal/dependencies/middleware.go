package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/handler/middleware"
)

func (d *Dependencies) Middleware() middleware.Middleware {
	if d.middleware == nil {
		d.middleware = middleware.NewMiddleware(*d.cfg.Handler,
			d.AuthService(), d.logger, d.cfg.Handler.QueueSize, d.UsersService())
	}
	return d.middleware

}
