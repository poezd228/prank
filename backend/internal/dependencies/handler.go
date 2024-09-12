package dependencies

import "github.com/mvd-inc/anibliss/internal/handler"

func (d *Dependencies) Handler() handler.Handler {
	if d.mainHandler == nil {
		d.mainHandler = handler.NewHandler(*d.cfg.Handler, d.Middleware(), d.logger, d.UsersService(), d.AuthService())
	}
	return d.mainHandler

}
