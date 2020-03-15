package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type SocietySuite struct {
	suite.Suite
	service *userService
	e       *echo.Echo
	db      *pg.DB
}

func (s *SocietySuite) SetupSuite() {
	var err error

	db := pg.Connect(&pg.Options{
		User:     "goo",
		Password: "goo",
		Database: "goo",
		Addr:     "localhost:5432",
	})

	s.service = CreateService(db)

	s.e = echo.New()
}

func (s *SocietySuite) TestGetIncidents() {
	req := httptest.NewRequest(echo.GET, "/incidents", nil)
	q := req.URL.Query()
	q.Add("from", "30")
	q.Add("ip", "54.194.80.144")
	req.URL.RawQuery = q.Encode()

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.Set("owner", "2")

	s.NoError(s.service.GetIncidents(c))
}

func (s *SocietySuite) TestGetIncidentsCount() {
	req := httptest.NewRequest(echo.GET, "/incidents/info", nil)
	q := req.URL.Query()
	q.Add("from", "30")
	req.URL.RawQuery = q.Encode()

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.Set("owner", "2")

	s.NoError(s.service.GetIncidentsCount(c))
}

func (s *SocietySuite) TearDownSuite() {

}

func (s *SocietySuite) SetupTest() {

}

func TestIncidentsServiceSuite(t *testing.T) {
	suite.Run(t, &SocietySuite{})
}
