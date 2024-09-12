package config

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type (
	AppMode string
	Config  struct {
		Server   *ServerConfig
		Handler  *HandlerConfig
		Postgres *PostgresConfig
		Auth     *AuthConfig
	}
	ServerConfig struct {
		Port int
	}
	HandlerConfig struct {
		RequestTimeout time.Duration
		SwaggerHost    string
		QueueSize      int
	}

	PostgresConfig struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     int

		ReadOnlyHost string
		ReadOnlyPort int

		ReadOnlyAllowed bool

		SSL     string
		CertLoc string
	}
	AuthConfig struct {
		Key                 string
		AccessTokenTimeout  time.Duration
		RefreshTokenTimeout time.Duration
		AuthTimeout         time.Duration
	}
)

func Init(configPath string) *Config {
	jsonCfg := viper.New()
	jsonCfg.AddConfigPath(filepath.Dir(configPath))
	jsonCfg.SetConfigName(filepath.Base(configPath))
	if err := jsonCfg.ReadInConfig(); err != nil {
		log.Panic("json Config init error")
	}
	envCfg := viper.New()
	envCfg.SetConfigFile(".env")
	if err := envCfg.ReadInConfig(); err != nil {
		log.Panic("env Config init error")
	}

	server := &ServerConfig{
		Port: jsonCfg.GetInt("server.port")}
	handler := &HandlerConfig{
		RequestTimeout: time.Millisecond * jsonCfg.GetDuration("handler.request_timeout"),
		SwaggerHost:    jsonCfg.GetString("handler.swagger_host"),
		QueueSize:      jsonCfg.GetInt("handler.queue_size")}
	dbCfg := &PostgresConfig{
		Host:         envCfg.GetString("POSTGRES_HOST"),
		User:         envCfg.GetString("POSTGRES_USER"),
		Password:     envCfg.GetString("POSTGRES_PASSWORD"),
		DBName:       envCfg.GetString("POSTGRES_DB"),
		Port:         envCfg.GetInt("POSTGRES_PORT"),
		SSL:          envCfg.GetString("POSTGRES_SSLMODE"),
		ReadOnlyHost: envCfg.GetString("POSTGRES_HOST"),
		ReadOnlyPort: envCfg.GetInt("POSTGRES_PORT"),
	}
	authCfg := &AuthConfig{
		Key:                 jsonCfg.GetString("auth.key"),
		AccessTokenTimeout:  jsonCfg.GetDuration("auth.access_token_timeout"),
		RefreshTokenTimeout: jsonCfg.GetDuration("auth.refresh_token_timeout"),
		AuthTimeout:         jsonCfg.GetDuration("auth.auth_timeout"),
	}
	return &Config{Server: server, Handler: handler, Postgres: dbCfg, Auth: authCfg}

}

func (p *PostgresConfig) PgSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}

func (p *PostgresConfig) PgReadOnlySource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		p.ReadOnlyHost, p.ReadOnlyPort, p.User, p.Password, p.DBName)
}
