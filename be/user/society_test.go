package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
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
		user          *models.User
		society       *models.Society
		numOfAdmins   int
		societyUpdate *models.Society
	}{
		{
			user:          &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			society:       &models.Society{Name: "Dake meno"},
			numOfAdmins:   1,
			societyUpdate: &models.Society{Name: "Nove menicko ako v restauracii"},
		},
	}

	for i, _ := range candidates {
		user, err := s.service.UserAccess.CreateUser(candidates[i].user)
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

		resp := &models.Society{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].society.Id = resp.Id
		candidates[i].society.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].society, resp)
		candidates[i].societyUpdate.Id = resp.Id
		candidates[i].societyUpdate.CreatedAt = resp.CreatedAt

		//oprav acces admina

		isAdmin, numAdmins, err := s.service.UserAccess.IsUserSocietyAdmin(candidates[i].user.Id, resp.Id)
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

		resp := &models.Society{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.societyUpdate, resp)

	}

}

func (s *SocietySuite) Test_ApplyForMembership_RemoveApplication_AllByUser() {
	candidates := []struct {
		admin           *models.User
		society         *models.Society
		newMember       *models.User
		applicationForm *models.IdMessage
		finalApplicant  *models.Applicant
		err             string
	}{
		{
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Motyl.cz", CreatedAt: time.Now()},
			society:   &models.Society{Name: "Dake meno"},
			newMember: &models.User{FirstName: "Novy", LastName: "Member", Email: "dakto@novy.cz"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		//for filling request structure
		candidates[i].applicationForm = &models.IdMessage{Id: candidates[i].society.Id}
		candidates[i].finalApplicant = &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
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

		resp := &models.Applicant{}
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
		admin           *models.User
		society         *models.Society
		newMember       *models.User
		applicationForm *models.IdMessage
		finalApplicant  *models.Applicant
		err             *custom_errors.ErrorModel
	}{
		{
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "ja@TestApplyFormMembershipExistingMember.com", CreatedAt: time.Now()},
			society:   &models.Society{Name: "TestApplyFormMembershipExistingMember"},
			newMember: &models.User{FirstName: "Novy", LastName: "Member", Email: "blbost@newMember.com"},
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrConflict, Message: "User is already a member"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		newMember := &models.Member{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id, Permission: models.Membership("member")}
		err := s.service.UserAccess.Db.Insert(newMember)
		s.Nil(err)

		testExistence := new(models.Member)
		err = s.db.Model(testExistence).Where("user_id = ?", candidates[i].newMember.Id).Select()
		if err != nil {
			s.Nil(err) //end test
		}
		if errors.Is(err, pg.ErrNoRows) {
			fmt.Println("Should be found something")
			s.Error(nil) //throw error in test
		}

		//for filling request structure
		candidates[i].applicationForm = &models.IdMessage{Id: candidates[i].society.Id}
		candidates[i].finalApplicant = &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
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

		resp := &custom_errors.ErrorModel{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)

		s.EqualValues(candidate.err, resp)
	}
}

func (s *SocietySuite) Test_DismissApplicant() {
	candidates := []struct {
		admin       *models.User
		society     *models.Society
		newMember   *models.User
		application *models.Applicant
	}{
		{
			admin:     &models.User{Id: "2", FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			society:   &models.Society{Name: "More members than one"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Email: "me@mew.cz"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		//for filling request structure
		candidates[i].application = &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(candidates[i].application)
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
		admin     *models.User
		society   *models.Society
		newMember *models.User
		request   *models.UserGroupRequest
		response  *models.Member
	}{
		{
			admin:     &models.User{FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			society:   &models.Society{Name: "More members than one"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Email: "Ja@Janovutbr.com"},
			response:  &models.Member{Permission: "member"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)
		candidates[i].response.UserId = candidates[i].newMember.Id
		candidates[i].response.SocietyId = candidates[i].society.Id

		//preparing request
		candidates[i].request = &models.UserGroupRequest{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}

		//preparing db
		application := &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(application)
		s.Nil(err)
	}

	for _, candidate := range candidates {
		body, err := json.Marshal(candidate.request)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/societies/approve", strings.NewReader(string(body)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", candidate.admin.Id)

		s.Nil(s.service.AcceptApplicant(c))

		resp := &models.Member{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.response.CreatedAt = resp.CreatedAt
		s.EqualValues(candidate.response, resp)
	}
}

func (s *SocietySuite) Test_ChangeRights() {
	candidates := []struct {
		admin         *models.User
		society       *models.Society
		friend        *models.User
		oldMembership *models.Member
		newMembership *models.Member
		err           string
	}{
		{
			admin:         &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Herrer.cz", CreatedAt: time.Now()},
			society:       &models.Society{Name: "Dake meno"},
			friend:        &models.User{FirstName: "Novy", LastName: "Member", Email: "Peter@me.cz"},
			oldMembership: &models.Member{Permission: models.Membership("member")},
			newMembership: &models.Member{Permission: models.Membership("admin")},
		},
		{
			admin:         &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:       &models.Society{Name: "Dake meno"},
			friend:        &models.User{FirstName: "Novy", LastName: "Member", Email: "me@me.cz"},
			oldMembership: &models.Member{Permission: models.Membership("admin")},
			newMembership: &models.Member{Permission: models.Membership("member")},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].friend, err = s.service.UserAccess.CreateUser(candidates[i].friend)
		s.Nil(err)

		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)

		candidates[i].oldMembership.UserId = candidates[i].friend.Id
		candidates[i].oldMembership.SocietyId = candidates[i].society.Id
		err = s.db.Insert(candidates[i].oldMembership)
		s.Nil(err)
		err := s.db.Select(candidates[i].oldMembership)
		if err != nil {
			log.Error(err)
			s.Nil(err)
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

		resp := &models.Member{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.newMembership.CreatedAt = resp.CreatedAt
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
	referencerTables := []string{
		"users",
		"societies",
		"societies_members",
		"societies_applicants",
		"events_users",
		"friends",
		"friend_requests",
	}
	referencerTableQueries := make([]string, len(referencerTables))
	for i, table := range referencerTables {
		if table == "spatial_ref_sys" { //postgis extension
			continue
		}
		referencerTableQueries[i] = "TRUNCATE " + table + " CASCADE;"
	}

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, query := range referencerTableQueries {
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Error(err)
	}
}

func TestSocietyServiceSuite(t *testing.T) {
	suite.Run(t, &SocietySuite{})
}
