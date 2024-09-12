package integration_tests

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mvd-inc/anibliss/internal/app"
	"github.com/mvd-inc/anibliss/internal/config"
	"github.com/mvd-inc/anibliss/internal/db"
	"github.com/mvd-inc/anibliss/internal/dependencies"
	"github.com/mvd-inc/anibliss/migrations"
	"github.com/mvd-inc/anibliss/pkg/logger"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	suite.Suite
	postgreSQL *PostgreSQLContainer
	pgClient   *db.PostgresClient

	cfg    *config.Config
	logger logger.Logger

	deps    *dependencies.Dependencies
	app     app.App
	handler http.Handler
}

func (s *TestSuite) SetupSuite() {
	var err error
	s.cfg = config.Init("../configs/dev")
	s.Require().NoError(err)
	s.logger, err = logger.NewLogger()
	s.Require().NoError(err)
}

func (s *TestSuite) SetupTest() {
	ctx := context.Background()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)
	s.postgreSQL = psqlContainer

	psqlPort, err := strconv.Atoi(psqlContainer.MappedPort)
	s.Require().NoError(err)

	s.cfg.Postgres.Port = psqlPort

	for i := 0; i < 20; i++ {
		s.pgClient, err = db.NewPostgresClient(context.Background(), s.postgreSQL.GetDSN())
		if err != nil && i != 19 {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		s.Require().NoError(err)
		break
	}
	err = s.pgClient.DB.Ping(ctx)
	s.Require().NoError(err)

	err = migrations.Migrate(s.postgreSQL.GetDSN())
	s.Require().NoError(err)

	err = s.setupData(s.pgClient)
	s.Require().NoError(err)

	if s.deps == nil {
		s.deps = dependencies.NewDependencies(s.cfg, s.logger)
	}
	s.app = app.NewApp(s.deps, s.logger)
	s.handler = s.app.GetHandler()
	s.app.Start()
}

func (s *TestSuite) TearDownTest() {
	ctx := context.Background()
	if s.app != nil {
		s.app.Stop()
	}

	if s.postgreSQL != nil {
		s.Require().NoError(s.postgreSQL.Terminate(ctx))
	}
}

func (s *TestSuite) setupData(pgClient *db.PostgresClient) error {

	//	pgClient.DB.Exec(context.Background(), `
	//	INSERT INTO tags (
	//		id,
	//		name
	//	)
	//	VALUES (DEFAULT,'pigs'),(DEFAULT,'lalka'),(DEFAULT,'o4ko')
	//	`)
	//
	//	pgClient.DB.Exec(context.Background(), `
	//	INSERT INTO geo (
	//		locale,
	//		price
	//	)
	//	VALUES ('en', 2000000000),('ru', 3000000000)
	//	`)
	//
	//	return nil
	return nil
}

type PostgreSQLContainer struct {
	testcontainers.Container
	MappedPort string
	Host       string
}

func (p *PostgreSQLContainer) GetDSN() string {
	return fmt.Sprintf("postgres://anibliss:1337@127.0.0.1:%s/anibliss", p.MappedPort)
}

func NewPostgreSQLContainer(ctx context.Context) (*PostgreSQLContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		Name:         "anibliss-postgres-test",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "tgquest",
			"POSTGRES_DB":       "tgquest",
			"POSTGRES_PASSWORD": "1337",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	started, err := checkPostgreSQLStarted(postgresContainer)
	if err != nil {
		return nil, err
	}
	if !started {
		return nil, fmt.Errorf("container not started")
	}
	return &PostgreSQLContainer{
		Container:  postgresContainer,
		MappedPort: mappedPort.Port(),
		Host:       host,
	}, nil
}

func checkPostgreSQLStarted(c testcontainers.Container) (bool, error) {
	for i := 0; i < 20; i++ {
		code, _, err := c.Exec(context.Background(), []string{"pg_isready", "-d", "tgquest", "-U", "tgquest"})
		if err != nil {
			return false, err
		}
		if code == 0 {
			return true, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false, nil
}
