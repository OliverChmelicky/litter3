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

func (s *SocietySuite) Test_CRUsociety() {
	candidates := []struct {
		user          *User
		society       *Society
		numOfAdmins   int
		societyUpdate *Society
	}{
		{
			user:          &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
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

func (s *SocietySuite) Test_ApplyFormMembership_RemoveApplication_AllByUser() {
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

		req := httptest.NewRequest(echo.POST, "/societies/apply", strings.NewReader(string(bytes)))

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

func (s *SocietySuite) Test_ApplyFormMembershipExistingMember() {
	candidates := []struct {
		admin           *User
		society         *Society
		newMember       *User
		applicationForm *UserGroupRequest
		finalApplicant  *Applicant
		err             string
	}{
		{
			admin:     &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "ja@TestApplyFormMembershipExistingMember.com", CreatedAt: time.Now()},
			society:   &Society{Name: "TestApplyFormMembershipExistingMember"},
			newMember: &User{FirstName: "Novy", LastName: "Member", Email: "blbost@newMember.com"},
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

		newMember := &Member{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id, Permission: membership("member")}
		err := s.service.userAccess.db.Insert(newMember)
		s.Nil(err)

		testExistence := new(Member)
		err = s.db.Model(testExistence).Where("user_id = ?", candidates[i].newMember.Id).Select()
		if err != nil {
			s.Nil(err) //end test
		}
		if err == pg.ErrNoRows {
			fmt.Println("Should be found something")
			s.Error(nil) //throw error in test
		}

		//for filling request structure
		candidates[i].applicationForm = &UserGroupRequest{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		candidates[i].finalApplicant = &Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.applicationForm)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/societies/apply", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.newMember.Id)

		s.Nil(s.service.ApplyForMembership(c))

		resp := &Applicant{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.NotNil(err)

		s.EqualValues("User is already a member", rec.Body.String())
	}
}

func (s *SocietySuite) Test_DismissApplicant() {
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
		req := httptest.NewRequest(echo.DELETE, "/societies/dismiss/"+cand.society.Id+"/"+cand.newMember.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", cand.admin.Id)

		c.SetPath("/societies/dismiss/:societyId/:userId")
		c.SetParamNames("societyId", "userId")
		c.SetParamValues(cand.society.Id, cand.newMember.Id)

		s.NoError(s.service.DismissApplicant(c))

		s.EqualValues("", rec.Body.String())
	}

}

func (s *SocietySuite) Test_AddMember() {
	candidates := []struct {
		admin       *User
		society     *Society
		newMember   *User
		application *Applicant
		response    *Member
	}{
		{
			admin:     &User{FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			society:   &Society{Name: "More members than one"},
			newMember: &User{FirstName: "Hello", LastName: "Flowup"},
			response:  &Member{Permission: "member"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].admin, err = s.service.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.userAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.userAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)
		candidates[i].response.UserId = candidates[i].newMember.Id
		candidates[i].response.SocietyId = candidates[i].society.Id

		//for filling request structure
		candidates[i].application = &Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.userAccess.AddApplicant(candidates[i].application)
		s.Nil(err)
	}

	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.POST, "/societies/"+candidate.application.SocietyId+"/"+candidate.application.UserId+"/approve", nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", candidate.newMember.Id)
		c.SetParamNames("societyId", "userId")
		c.SetParamValues(candidate.society.Id, candidate.newMember.Id)

		s.Nil(s.service.AcceptApplicant(c))

		resp := &Member{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.response, resp)
	}
}

//TODO
//uz mam v skupine admina a membera
//changeUserRights to admin
func (s *SocietySuite) Test_ChangeRights() {
	candidates := []struct {
		admin         *User
		society       *Society
		friend        *User
		oldMembership *Member
		newMembership *Member
		err           string
	}{
		{
			admin:         &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:       &Society{Name: "Dake meno"},
			friend:        &User{FirstName: "Novy", LastName: "Member"},
			oldMembership: &Member{Permission: membership("member")},
			newMembership: &Member{Permission: membership("admin")},
		},
		{
			admin:         &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:       &Society{Name: "Dake meno"},
			friend:        &User{FirstName: "Novy", LastName: "Member"},
			oldMembership: &Member{Permission: membership("admin")},
			newMembership: &Member{Permission: membership("member")},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].friend, err = s.service.userAccess.CreateUser(candidates[i].friend)
		s.Nil(err)

		candidates[i].society, err = s.service.userAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		candidates[i].oldMembership.UserId = candidates[i].friend.Id
		candidates[i].oldMembership.SocietyId = candidates[i].society.Id
		err = s.db.Insert(candidates[i].oldMembership)
		s.Nil(err)
		err := s.db.Select(candidates[i].oldMembership)
		if err != nil {
			fmt.Println(err)
		}

		candidates[i].newMembership.UserId = candidates[i].friend.Id
		candidates[i].newMembership.SocietyId = candidates[i].society.Id
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.newMembership)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/society/members/update", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.admin.Id)

		s.NoError(s.service.ChangeMemberRights(c))

		resp := &Member{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.newMembership, resp)
	}
}

//remove admin by another admin
//remove member by admin
//test remove sam seba(som admin) a som posledny admin v society

//test delete society //delete user --> complicated pockaj si na odpoved veduceho
//ale delete society mozem spravit dvomi sposobmi, odide posledny admin, alebo to admin zrusi priamo

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
