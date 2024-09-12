package cron

import (
	"context"
	"time"

	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/mvd-inc/anibliss/pkg/time_manager"
)

//go:generate mockgen -source cron.go -destination ../../mocks/service/cron/cron.go

type Cron interface {
	RunJobs()
	Stop()
}

type cron struct {
	cfg         config.ServerConfig
	timeManager time_manager.TimeManager
	closeCh     chan any
	ctxMain     context.Context
	ctxCancel   context.CancelFunc
	l           logger.Logger
}

func NewCron(
	cfg config.ServerConfig,
	timeManager time_manager.TimeManager,
	l logger.Logger) Cron {
	ctx, cancel := context.WithCancel(context.Background())
	return &cron{
		cfg:         cfg,
		timeManager: timeManager,
		closeCh:     make(chan any),
		ctxMain:     ctx,
		ctxCancel:   cancel,
		l:           l,
	}
}

func (c *cron) RunJobs() {
	go c.runJobs()
}

func (c *cron) runJobs() {
	for {
		// время следующего планового срабатывания (по мск часовому поясу, то есть utc+3)
		nextTime := c.nextTime()
		now := c.timeManager.Now()
		select {
		case <-time.After(nextTime.Sub(now)):
			break
		case <-c.closeCh:
			return
		}
		c.runIteration(nextTime)

	}
}

func (c *cron) runIteration(nextTime time.Time) {

}

func (c *cron) Stop() {
	c.ctxCancel()
	c.closeCh <- struct{}{}
}

func (c *cron) nextTime() time.Time {
	now := c.timeManager.Now()
	next := now.Add(30 * time.Second)
	res := c.truncate(next, 30*time.Second)
	return res
}

func (c *cron) truncate(t time.Time, duration time.Duration) time.Time {
	timestamp := t.UnixMilli()
	durationMs := duration.Milliseconds()
	timestamp = durationMs * (timestamp / durationMs)
	return c.timeManager.MillisecondsToTime(timestamp)
}
