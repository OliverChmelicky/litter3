package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func (s *SocietySuite) TestCRU_Society() {
	candidates := []struct {
		user          *User
		society       *Society
		numOfAdmins   int
		societyUpdate *Society
	}{
		{
			user:          &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:       &Society{Name: "Dake meno"},
			numOfAdmins:   1,
			societyUpdate: &Society{Name: "Nove menicko ako v restauracii"},
		},
	}

	for i, _ := range candidates {
		user, err := s.service.userAccess.CreateUser(candidates[i].user)
		candidates[i].user = user
		s.NoError(err)

		bytes, err := json.Marshal(candidates[i].society)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/societies/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].user.Id)

		s.NoError(s.service.CreateSociety(c))

		resp := &Society{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].society.Id = resp.Id
		candidates[i].society.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].society, resp)
		candidates[i].societyUpdate.Id = resp.Id
		candidates[i].societyUpdate.CreatedAt = resp.CreatedAt

		//oprav acces admina

		isAdmin, numAdmins, err := s.service.isUserSocietyAdmin(candidates[i].user.Id, resp.Id)
		s.Nil(err)
		s.True(isAdmin)
		s.EqualValues(candidates[i].numOfAdmins, numAdmins)
	}

	//test update group
	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.societyUpdate)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/societies/update", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.user.Id)

		s.NoError(s.service.UpdateSociety(c))

		resp := &Society{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		fmt.Println(err)
		s.Nil(err)

		s.EqualValues(candidate.societyUpdate, resp)

	}

}

//add admin
//remove admin by another admin
//test remove member

//test ApplyForMembership nie je clenom
////test ApplyForMembership je uz clenom
////test removeApplication for membership

//test add member
//test dismiss applicant

//test delete society --> complicated

func (s *SocietySuite) TearDownSuite() {
	s.db.Close()
}

func (s *SocietySuite) SetupTest() {

}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &SocietySuite{})
}
