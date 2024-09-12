package integration_tests

import (
	"github.com/golang/mock/gomock"
	"github.com/mvd-inc/anibliss/internal/db"
	"github.com/mvd-inc/anibliss/internal/dependencies"
	mock_cron "github.com/mvd-inc/anibliss/internal/mocks/service/cron"
	mock_time_manager "github.com/mvd-inc/anibliss/mocks/pkg/time_manager"
)

type LoginSuite struct {
	TestSuite
	ctrl            *gomock.Controller
	mockTimeManager *mock_time_manager.MockTimeManager
	mockCron        *mock_cron.MockCron
}

func (s *LoginSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
}

func (s *LoginSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockTimeManager = mock_time_manager.NewMockTimeManager(s.ctrl)
	s.mockCron = mock_cron.NewMockCron(s.ctrl)
	s.mockCron.EXPECT().RunJobs()
	s.mockCron.EXPECT().Stop()
	s.deps = dependencies.NewDependencies(
		s.cfg,
		s.logger,
		dependencies.WithTimeManager(s.mockTimeManager),
		dependencies.WithCron(s.mockCron),
	)
	s.TestSuite.SetupTest()
	err := s.setupData(s.pgClient)
	s.Require().NoError(err)
}

func (s *LoginSuite) setupData(pgClient *db.PostgresClient) error {
	return nil
}

func (s *LoginSuite) TestLogin() {
	testConte

}
