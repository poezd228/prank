package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/repository/auth"
	"github.com/mvd-inc/anibliss/internal/repository/jwt"
	"github.com/mvd-inc/anibliss/internal/repository/transactions"
	"github.com/mvd-inc/anibliss/internal/repository/users"
)

func (d *Dependencies) TransactionRepo() transactions.Repository {
	if d.transactionRepo == nil {
		d.transactionRepo = transactions.NewTxRepository(d.DbClient())
	}
	return d.transactionRepo

}
func (d *Dependencies) JwtRepo() jwt.Repository {
	if d.jwtRepo == nil {
		d.jwtRepo = jwt.NewRepository()
	}
	return d.jwtRepo

}

func (d *Dependencies) UsersRepo() users.Repository {
	if d.usersRepo == nil {
		d.usersRepo = users.NewRepository()

	}
	return d.usersRepo

}
func (d *Dependencies) AuthRepo() auth.Repository {
	if d.authRepo == nil {
		return auth.NewRepository()
	}
	return d.authRepo

}
