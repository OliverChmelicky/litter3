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
	user := &UserModel{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", Created: 538}
	newSociety := &SocietyModel{Name: "Dake meno"}
	numberOfAdmins := 1

	user, err := s.service.userAccess.CreateUser(user)
	s.NoError(err)

	bytes, _ := json.Marshal(newSociety)

	req := httptest.NewRequest(echo.POST, "/societies/new", strings.NewReader(string(bytes)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)
	c.Set("userId", user.Id)

	s.NoError(s.service.CreateSociety(c))

	resp := &SocietyModel{}
	err = json.Unmarshal(rec.Body.Bytes(), resp)
	s.Nil(err)
	fmt.Println(resp)

	newSociety.Id = resp.Id
	newSociety.Created = resp.Created
	s.EqualValues(newSociety, resp)

	isAdmin, numAdmins, err := s.service.isUserSocietyAdmin(user.Id, resp.Id)
	s.Nil(err)
	s.True(isAdmin)
	s.EqualValues(numberOfAdmins, numAdmins)

	//test update group

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
