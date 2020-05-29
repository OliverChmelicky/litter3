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

	s.service = CreateService(db, nil, nil)
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
			user:          &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", Uid: "123", CreatedAt: time.Now()},
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
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Motyl.cz", Uid: "123", CreatedAt: time.Now()},
			society:   &models.Society{Name: "Dake meno"},
			newMember: &models.User{FirstName: "Novy", LastName: "Member", Uid: "321", Email: "dakto@novy.cz"},
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
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "ja@TestApplyFormMembershipExistingMember.com", Uid: "987", CreatedAt: time.Now()},
			society:   &models.Society{Name: "TestApplyFormMembershipExistingMember"},
			newMember: &models.User{FirstName: "Novy", LastName: "Member", Uid: "5547", Email: "blbost@newMember.com"},
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
			admin:     &models.User{Id: "2", FirstName: "John", LastName: "Modest", Uid: "987654", Email: "Ja@Janovutbr.cz"},
			society:   &models.Society{Name: "More members than one"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "564", Email: "me@mew.cz"},
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
			admin:     &models.User{FirstName: "John", LastName: "Modest", Uid: "45", Email: "Ja@Janovutbr.cz"},
			society:   &models.Society{Name: "More members than one"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "2", Email: "Ja@Janovutbr.com"},
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
		newMembership []models.Member
		err           string
	}{
		{
			admin:         &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "9", Email: "Ja@Herrer.cz", CreatedAt: time.Now()},
			society:       &models.Society{Name: "Dake meno"},
			friend:        &models.User{FirstName: "Novy", LastName: "Member", Uid: "3", Email: "Peter@me.cz"},
			oldMembership: &models.Member{Permission: models.Membership("member")},
			newMembership: []models.Member{{Permission: models.Membership("admin")}},
		},
		{
			admin:         &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "7", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			society:       &models.Society{Name: "Dake meno"},
			friend:        &models.User{FirstName: "Novy", LastName: "Member", Uid: "1", Email: "me@me.cz"},
			oldMembership: &models.Member{Permission: models.Membership("admin")},
			newMembership: []models.Member{{Permission: models.Membership("member")}},
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

		candidates[i].newMembership[0].UserId = candidates[i].friend.Id
		candidates[i].newMembership[0].SocietyId = candidates[i].society.Id
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

		var resp []models.Member
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		candidate.newMembership[0].CreatedAt = resp[0].CreatedAt
		s.EqualValues(candidate.newMembership, resp)
	}
}

func (s *SocietySuite) Test_GetSocietyMembers() {
	candidates := []struct {
		admin     *models.User
		society   *models.Society
		newMember *models.User
	}{
		{
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "3", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			society:   &models.Society{Name: "Dake meno"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "8", Email: "Ja@Janovutbr.com"},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		application := &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(application)
		s.Nil(err)

		_, err = s.service.UserAccess.AcceptApplicant(application.UserId, application.SocietyId)
		s.Nil(err)
	}

	for _, candidate := range candidates {
		members, err := s.service.UserAccess.GetSociety(candidate.society.Id)
		s.Nil(err)

		s.EqualValues(2, len(members.Users))
	}
}

func (s *SocietySuite) Test_GetSocietyAdmins() {
	candidates := []struct {
		admin     *models.User
		society   *models.Society
		newMember *models.User
		relation  models.Member
	}{
		{
			admin:     &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "3", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			society:   &models.Society{Name: "Dake meno"},
			newMember: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "8", Email: "Ja@Janovutbr.com"},
			relation:  models.Member{},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)
		candidates[i].newMember, err = s.service.UserAccess.CreateUser(candidates[i].newMember)
		s.Nil(err)

		candidates[i].relation.Permission = models.Membership("admin")
		candidates[i].relation.UserId = candidates[i].admin.Id
		candidates[i].relation.SocietyId = candidates[i].society.Id

		application := &models.Applicant{UserId: candidates[i].newMember.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(application)
		s.Nil(err)

		_, err = s.service.UserAccess.AcceptApplicant(application.UserId, application.SocietyId)
		s.Nil(err)
	}

	for i, candidate := range candidates {
		req := httptest.NewRequest(echo.POST, "/societies/"+candidate.society.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetParamNames("societyId")
		c.SetParamValues(candidate.society.Id)

		s.NoError(s.service.GetSocietyAdmins(c))

		var resp []string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		s.EqualValues(candidates[i].relation.UserId, resp[0])
	}
}

//TODO remove admin by another admin
//TODO remove member by admin

func (s *SocietySuite) Test_DeleteSociety() {
	candidates := []struct {
		admin                      *models.User
		society                    *models.Society
		memberAndEventAttendant    *models.User
		applicantAndEventAttendant *models.User
		event                      *models.Event
		trash                      *models.Trash //check trash-event
		collection                 *models.Collection
	}{
		{
			admin:                      &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "3", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			society:                    &models.Society{Name: "Dake meno"},
			memberAndEventAttendant:    &models.User{FirstName: "Hello", LastName: "Friend", Uid: "8", Email: "Ja@Janovutbr.com"},
			applicantAndEventAttendant: &models.User{FirstName: "Hello", LastName: "Fero", Uid: "9", Email: "Ja@Maria.cz"},
			event:                      &models.Event{Id: "123", Date: time.Now(), CreatedAt: time.Now()},
			trash:                      &models.Trash{Id: "321", CreatedAt: time.Now(), Location: models.Point{0, 0}},
			collection:                 &models.Collection{Id: "111", CreatedAt: time.Now(), TrashId: "321", EventId: "123", Weight: 32},
		},
	}

	//create user, his society add member and applicant
	//create trash, event, add two attendants there + society is admin
	//create collection with eventId and create trash-event relation
	for i, _ := range candidates {
		var err error
		//first part
		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
		s.Nil(err)
		candidates[i].memberAndEventAttendant, err = s.service.UserAccess.CreateUser(candidates[i].memberAndEventAttendant)
		s.Nil(err)
		candidates[i].applicantAndEventAttendant, err = s.service.UserAccess.CreateUser(candidates[i].applicantAndEventAttendant)
		s.Nil(err)

		application := &models.Applicant{UserId: candidates[i].memberAndEventAttendant.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(application)
		s.Nil(err)
		_, err = s.service.UserAccess.AcceptApplicant(application.UserId, application.SocietyId)
		s.Nil(err)

		application = &models.Applicant{UserId: candidates[i].applicantAndEventAttendant.Id, SocietyId: candidates[i].society.Id}
		_, err = s.service.UserAccess.AddApplicant(application)
		s.Nil(err)

		//second part
		err = s.db.Insert(candidates[i].trash)
		s.Nil(err)
		err = s.db.Insert(candidates[i].event)
		s.Nil(err)
		err = s.db.Insert(&models.EventUser{UserId: candidates[i].applicantAndEventAttendant.Id, EventId: candidates[i].event.Id, Permission: "viewer"})
		s.Nil(err)
		err = s.db.Insert(&models.EventUser{UserId: candidates[i].memberAndEventAttendant.Id, EventId: candidates[i].event.Id, Permission: "viewer"})
		s.Nil(err)
		err = s.db.Insert(&models.EventSociety{SocietyId: candidates[i].society.Id, EventId: candidates[i].event.Id, Permission: "creator"})
		s.Nil(err)

		//third part
		candidates[i].collection.EventId = candidates[i].event.Id
		candidates[i].collection.TrashId = candidates[i].trash.Id

		err = s.db.Insert(candidates[i].collection)
		s.Nil(err)
		err = s.db.Insert(&models.EventTrash{TrashId: candidates[i].trash.Id, EventId: candidates[i].event.Id})
		s.Nil(err)
	}

	//create user, his society add member and applicant
	//create trash, event, add two attendants there + society is admin
	//create collection with eventId and create trash-event relation
	for _, candidate := range candidates {
		err := s.db.Model(&models.Society{}).Where("id = ?", candidate.society.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.Applicant{}).Where("user_id = ?", candidate.applicantAndEventAttendant.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.Member{}).Where("user_id = ?", candidate.memberAndEventAttendant.Id).Select()
		s.Nil(err)

		err = s.db.Model(&models.Trash{}).Where("id = ?", candidate.trash.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.Event{}).Where("id = ?", candidate.event.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.memberAndEventAttendant.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.applicantAndEventAttendant.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.EventSociety{}).Where("society_id = ?", candidate.society.Id).Select()
		s.Nil(err)

		err = s.db.Model(&models.Collection{}).Where("id = ?", candidate.collection.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.EventTrash{}).Where("trash_id = ?", candidate.trash.Id).Select()
		s.Nil(err)

		req := httptest.NewRequest(echo.DELETE, "/societies/"+candidate.society.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetParamNames("societyId")
		c.SetParamValues(candidate.society.Id)

		c.Set("userId", candidate.admin.Id)

		s.NoError(s.service.DeleteSociety(c))

		s.EqualValues(200, rec.Code)

		err = s.db.Model(&models.Society{}).Where("id = ?", candidate.society.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.Applicant{}).Where("user_id = ?", candidate.applicantAndEventAttendant.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.Member{}).Where("user_id = ?", candidate.memberAndEventAttendant).Select()
		s.NotNil(err)

		err = s.db.Model(&models.Trash{}).Where("id = ?", candidate.trash.Id).Select()
		s.Nil(err)
		err = s.db.Model(&models.Event{}).Where("id = ?", candidate.event.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.memberAndEventAttendant.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.applicantAndEventAttendant.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.EventSociety{}).Where("society_id = ?", candidate.society.Id).Select()
		s.NotNil(err)

		err = s.db.Model(&models.Collection{}).Where("id = ?", candidate.collection.Id).Select()
		s.NotNil(err)
		err = s.db.Model(&models.EventTrash{}).Where("trash_id = ?", candidate.trash.Id).Select()
		s.NotNil(err)
	}
}

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
