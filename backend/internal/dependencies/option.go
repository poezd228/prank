package dependencies

import (
	"github.com/mvd-inc/anibliss/internal/service/cron"
	"github.com/mvd-inc/anibliss/pkg/time_manager"
)

type Option interface {
	Apply(d *Dependencies)
}

type timeManagerOption struct {
	t time_manager.TimeManager
}

func (o *timeManagerOption) Apply(d *Dependencies) {
	d.timeManager = o.t
}

func WithTimeManager(t time_manager.TimeManager) Option {
	return &timeManagerOption{
		t: t,
	}
}

type cronOption struct {
	c cron.Cron
}

func (o *cronOption) Apply(d *Dependencies) {
	d.cron = o.c
}

func WithCron(c cron.Cron) Option {
	return &cronOption{
		c: c,
	}
}
