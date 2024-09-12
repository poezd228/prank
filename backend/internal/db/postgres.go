package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type PostgresClient struct {
	DB              *pgxpool.Pool
	ReadOnlyDB      *pgxpool.Pool
	ReadOnlyAllowed bool
}

func prepareCfg(cfg *pgxpool.Config) {
	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Second
	if _, ok := cfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"]; !ok {
		cfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = fmt.Sprintf("%d", (20 * time.Second).Milliseconds())
	}
	if _, ok := cfg.ConnConfig.RuntimeParams["statement_timeout"]; !ok {
		cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprintf("%d", (20 * time.Second).Milliseconds())
	}
}

func NewPostgresClient(ctx context.Context, DSN string) (*PostgresClient, error) {
	cfg, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	prepareCfg(cfg)

	client, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the postgresql database: %w", err)
	}
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error while ping to postgres")
	}

	return &PostgresClient{
		DB:              client,
		ReadOnlyDB:      client,
		ReadOnlyAllowed: false,
	}, nil
}

func NewPostgresqlClientWithReadWriteSplit(ctx context.Context, readOnlyDSN, readWriteDSN, rootLoc string) (*PostgresClient, error) {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(rootLoc)
	if err != nil {
		return nil, fmt.Errorf("error reading cert: %w", err)
	}
	rootCertPool.AppendCertsFromPEM(pem)

	readOnlyCfg, err := pgxpool.ParseConfig(readOnlyDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}
	readOnlyCfg.ConnConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}
	prepareCfg(readOnlyCfg)

	readOnlyClient, err := pgxpool.NewWithConfig(ctx, readOnlyCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the postgresql database: %w", err)
	}
	if err := readOnlyClient.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error while ping to postgres: %w", err)
	}

	readWriteCfg, err := pgxpool.ParseConfig(readWriteDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}
	readWriteCfg.ConnConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}
	prepareCfg(readWriteCfg)

	readWriteClient, err := pgxpool.NewWithConfig(ctx, readWriteCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the postgresql database: %w", err)
	}
	if err := readWriteClient.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error while ping to postgres")
	}

	return &PostgresClient{
		DB:              readWriteClient,
		ReadOnlyDB:      readOnlyClient,
		ReadOnlyAllowed: false,
	}, nil
}

func (c *PostgresClient) Close() {
	c.DB.Close()
	if c.ReadOnlyAllowed {
		c.ReadOnlyDB.Close()
	}
}
