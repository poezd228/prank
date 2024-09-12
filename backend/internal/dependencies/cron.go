package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/service/cron"
)

func (d *Dependencies) Cron() cron.Cron {
	if d.cron == nil {
		d.cron = cron.NewCron(*d.cfg.Server, d.TimeManager(), d.logger)
	}
	return d.cron

}
