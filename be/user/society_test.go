package user

import (
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"strings"
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

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Error("PostgreSQL is down")
	}

	s.service = CreateService(db)
	s.db = db

	s.e = echo.New()
}

func (s *SocietySuite) TestCreateSociety() {
	user := &UserModel{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", Created: 538}

	user, err := s.service.userAccess.CreateUser(user)
	s.NoError(err)

	newSociety := SocietyModel{Name: "Dake meno"}
	bytes, _ := json.Marshal(newSociety)

	req := httptest.NewRequest(echo.POST, "/societies/new", strings.NewReader(string(bytes)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.Set("userId", user.Id)

	s.NoError(s.service.CreateSociety(c))
}

//func (s *SocietySuite) TestGetIncidentsCount() {
//	req := httptest.NewRequest(echo.GET, "/incidents/info", nil)
//	q := req.URL.Query()
//	q.Add("from", "30")
//	req.URL.RawQuery = q.Encode()
//
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := s.e.NewContext(req, rec)
//	c.Set("owner", "2")
//
//	s.NoError(s.service.GetIncidentsCount(c))
//}

func (s *SocietySuite) TearDownSuite() {
	s.db.Close()
}

func (s *SocietySuite) SetupTest() {

}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &SocietySuite{})
}
