package dependencies

import (
	"context"

	"github.com/mvd-inc/anibliss/internal/db"
)

func (d *Dependencies) DbClient() *db.PostgresClient {
	if d.DBClient == nil {
		var err error
		if d.cfg.Postgres.ReadOnlyAllowed {
			d.DBClient, err = db.NewPostgresqlClientWithReadWriteSplit(
				context.Background(),
				d.cfg.Postgres.PgReadOnlySource(),
				d.cfg.Postgres.PgSource(),
				d.cfg.Postgres.CertLoc,
			)
			if err != nil {
				d.logger.Panic("failed to init db client ", err, d.cfg.Postgres.PgSource(), d.cfg.Postgres.PgReadOnlySource())
			}
		} else {
			d.DBClient, err = db.NewPostgresClient(context.Background(), d.cfg.Postgres.PgSource())
			if err != nil {
				d.logger.Panic("failed to init db client ", err, d.cfg.Postgres.PgSource(), d.cfg.Postgres.PgReadOnlySource())
			}
		}
		d.stopFuncs = append(d.stopFuncs, d.DBClient.Close)
	}
	return d.DBClient
}
