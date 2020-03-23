package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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
		log.Error("PostgresSQL is down")
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
		s.Nil(err)

		s.EqualValues(candidate.societyUpdate, resp)

	}

}

func (s *SocietySuite) TestMembershipApplication_Apply_Remove() {
	candidates := []struct {
		admin           *User
		society         *Society
		newMember       *User
		applicationForm *UserGroupRequest
		finalApplicant  *Applicant
		err             string
	}{
		{
			admin:     &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:   &Society{Name: "Dake meno"},
			newMember: &User{FirstName: "Novy", LastName: "Member"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.userAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.userAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		//for filling request structure
		candidates[i].applicationForm = &UserGroupRequest{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		candidates[i].finalApplicant = &Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.applicationForm)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/societies/apply", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.newMember.Id)

		s.NoError(s.service.ApplyForMembership(c))

		resp := &Applicant{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.finalApplicant.CreatedAt = resp.CreatedAt
		s.EqualValues(candidate.finalApplicant, resp)

		//remove application
		req = httptest.NewRequest(echo.DELETE, "/society/"+candidate.society.Id, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = s.e.NewContext(req, rec)

		c.SetParamNames("societyId")
		c.SetParamValues(candidate.society.Id)
		c.Set("userId", candidate.newMember.Id)

		s.NoError(s.service.RemoveApplicationForMembership(c))
		s.EqualValues("", rec.Body.String())
	}
}

////TODO test ApplyForMembership ked uz je clenom society
//func (s *Society) TestApplyFormMembershipExistingMember() {
//
//}

func (s *SocietySuite) TestApproveMember_RemoveApplication() {
	candidates := []struct {
		admin       *User
		society     *Society
		newMember   *User
		application *Applicant
	}{
		{
			admin:     &User{Id: "2", FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			society:   &Society{Name: "More members than one"},
			newMember: &User{FirstName: "Hello", LastName: "Flowup"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.userAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.userAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		//for filling request structure
		candidates[i].application = &Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.userAccess.AddApplicant(candidates[i].application)
		s.Nil(err)
	}

	for _, cand := range candidates {
		req := httptest.NewRequest(echo.PUT, "/societies/dismiss/"+cand.society.Id+"/"+cand.newMember.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", cand.admin.Id)

		c.SetPath("/societies/dismiss/:societyId/:userId")
		c.SetParamNames("societyId", "userId")
		c.SetParamValues(cand.society.Id, cand.newMember.Id)

		s.NoError(s.service.DismissApplicant(c))
		fmt.Println(rec.Code)
		fmt.Println(rec.Body.String())

		s.EqualValues("", rec.Body.String())
	}

}

func (s *SocietySuite) TestApproveMember_AddMember() {

}

//uz mam v skupine admina a membera
//changeUserRights to admin
//remove admin by another admin
//test remove member

//test delete society //delete user --> complicated pockaj si na odpoved veduceho

func (s *SocietySuite) TearDownSuite() {
	s.db.Close()
}

func (s *SocietySuite) SetupTest() {
	s.Nil(s.db.DropTable((*User)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Society)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Applicant)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Member)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))

	s.Nil(s.db.CreateTable((*User)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Society)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Applicant)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Member)(nil), &orm.CreateTableOptions{IfNotExists: true}))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &SocietySuite{})
}
