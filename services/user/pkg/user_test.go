package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	service *userService
	e       *echo.Echo
	db      *pg.DB
}

func (s *UserSuite) SetupSuite() {
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

	s.service = CreateService(db, nil)
	s.db = db

	s.e = echo.New()
}

// FRIENDSHIP
func (s *UserSuite) Test_ApplyForFriendship_RemoveRequest_AllByUser() {
	candidates := []struct {
		heinrich        *models.User
		peterAsks       *models.User
		applicationForm *models.IdMessage
		friendship      *models.FriendRequest
		err             *custom_errors.ErrorModel
	}{
		{
			heinrich:  &models.User{Id: "1", FirstName: "Heinrich", LastName: "Herrer", Uid: "7", Email: "Heinrich@Herrer.tibet", CreatedAt: time.Now()},
			peterAsks: &models.User{FirstName: "Novy", LastName: "Member", Uid: "5", Email: "Ja@Peter.cz"},
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrConflict},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)
		candidates[i].applicationForm = &models.IdMessage{Id: candidates[i].heinrich.Id}

		//filling correct answer
		if strings.Compare(candidates[i].heinrich.Id, candidates[i].peterAsks.Id) == -1 {
			candidates[i].friendship = &models.FriendRequest{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		} else {
			candidates[i].friendship = &models.FriendRequest{User1Id: candidates[i].peterAsks.Id, User2Id: candidates[i].heinrich.Id}
		}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.applicationForm)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/users/friend/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.ApplyForFriendshipById(c))

		resp := &models.FriendRequest{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.friendship.CreatedAt = resp.CreatedAt
		correctIdOrderFriendRequest(candidate.friendship, resp)
		s.EqualValues(candidate.friendship, resp)

		//try againg, should throw an error
		bytes, err = json.Marshal(candidate.applicationForm)
		s.Nil(err)

		req = httptest.NewRequest(echo.POST, "/users/friend/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = s.e.NewContext(req, rec)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.ApplyForFriendshipById(c))

		respErr := &custom_errors.ErrorModel{}
		err = json.Unmarshal(rec.Body.Bytes(), respErr)
		s.Nil(err)
		candidate.err.Message = respErr.Message
		s.EqualValues(candidate.err, respErr)

		//remove application
		req = httptest.NewRequest(echo.DELETE, "/users/friend/deny/"+candidate.peterAsks.Id, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = s.e.NewContext(req, rec)

		c.SetParamNames("notWanted")
		c.SetParamValues(candidate.heinrich.Id)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.RemoveApplicationForFriendship(c))
		s.EqualValues("", rec.Body.String())

		rq := &models.FriendRequest{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User2Id, candidate.friendship.User1Id).Select()
		s.EqualValues(pg.ErrNoRows, err)
	}
}

func (s *UserSuite) Test_RequestFriendshipExistingFriendship() {
	candidates := []struct {
		heinrich        *models.User
		peterAsks       *models.User
		applicationForm *models.IdMessage
		err             *custom_errors.ErrorModel
	}{
		{
			heinrich:  &models.User{Id: "1", FirstName: "Heinrich", LastName: "Herrer", Uid: "5", Email: "ja@TestApplyFormMembershipExistingMember.com", CreatedAt: time.Now()},
			peterAsks: &models.User{FirstName: "Novy", LastName: "Member", Uid: "7", Email: "blbost@peterAsks.com"},
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrConflict, Message: "YOU ARE FIENDS ALREADY"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)
		candidates[i].applicationForm = &models.IdMessage{Id: candidates[i].heinrich.Id}

		//creating friendship
		existingFriendship := &models.Friends{User1Id: candidates[i].peterAsks.Id, User2Id: candidates[i].heinrich.Id}

		err := s.service.UserAccess.Db.Insert(existingFriendship)
		s.Nil(err)

		testExistence := new(models.Friends)
		err = s.db.Model(testExistence).Where("user1_id = ? and user2_id = ?", existingFriendship.User1Id, existingFriendship.User2Id).Select()
		if err != nil {
			s.Nil(err) //end test
		}
		if errors.Is(err, pg.ErrNoRows) {
			s.Error(nil) //throw error in test
		}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.applicationForm)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/users/friend/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.peterAsks.Id)

		s.Nil(s.service.ApplyForFriendshipById(c))

		resp := &custom_errors.ErrorModel{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.err, resp)
	}
}

func (s *UserSuite) Test_AddFriend() {
	candidates := []struct {
		heinrich  *models.User
		peterAsks *models.User
		request   *models.FriendRequest
		response  *models.Friends
	}{
		{
			heinrich:  &models.User{FirstName: "John", LastName: "Modest", Uid: "7", Email: "On@Janovutbr.com"},
			peterAsks: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "5", Email: "TY@Janovutbr.cz"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)

		//for filling request structure
		candidates[i].request = &models.FriendRequest{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		_, err = s.service.UserAccess.AddFriendshipRequest(candidates[i].request)
		s.Nil(err)
		//fir filling final answer
		candidates[i].response = &models.Friends{User1Id: candidates[i].peterAsks.Id, User2Id: candidates[i].heinrich.Id}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.request)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/users/friend/accept", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", candidate.heinrich.Id)

		c.SetParamNames("wantedUser")
		c.SetParamValues(candidate.peterAsks.Id)

		s.Nil(s.service.AcceptFriendship(c))

		resp := &models.Friends{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.response.CreatedAt = resp.CreatedAt
		s.EqualValues(candidate.response, resp)
	}
}

func (s *UserSuite) Test_RemoveFriend() {
	candidates := []struct {
		heinrich   *models.User
		peterAsks  *models.User
		friendship *models.Friends
	}{
		{
			heinrich:  &models.User{FirstName: "John", LastName: "Modest", Uid: "6", Email: "Ja@Janovutbr.italy"},
			peterAsks: &models.User{FirstName: "Hello", LastName: "Flowup", Uid: "5", Email: "Ja@Milan.cz"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)

		//create friendship
		candidates[i].friendship = &models.Friends{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		err = s.db.Insert(candidates[i].friendship)
		s.Nil(err)
	}

	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.DELETE, "/users/friends/remove/"+candidate.heinrich.Id, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetParamNames("notWanted")
		c.SetParamValues(candidate.heinrich.Id)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.RemoveFriend(c))
		s.EqualValues(http.StatusOK, rec.Code)

		rq := &models.Friends{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User1Id, candidate.friendship.User2Id).Select()
		s.EqualValues(pg.ErrNoRows, err)
	}
}

// commented because Firebase will be nil
// func (s *SocietySuite) Test_DeleteUser() {
// 	candidates := []struct {
// 		admin                                     *models.User
// 		society                                   *models.Society
// 		memberAndEventAttendantAndFriendRequester *models.User
// 		applicantAndEventAttendantAndFriend       *models.User
// 		event                                     *models.Event
// 		trash                                     *models.Trash //check trash-event
// 		collection                                *models.Collection
// 	}{
// 		{
// 			admin:   &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "3", Email: "Ja@kamo.com", CreatedAt: time.Now()},
// 			society: &models.Society{Name: "Dake meno"},
// 			memberAndEventAttendantAndFriendRequester: &models.User{FirstName: "Hello", LastName: "Friend", Uid: "8", Email: "Ja@Janovutbr.com"},
// 			applicantAndEventAttendantAndFriend:       &models.User{FirstName: "Hello", LastName: "Fero", Uid: "9", Email: "Ja@Maria.cz"},
// 			event:                                     &models.Event{Id: "123", Date: time.Now(), CreatedAt: time.Now()},
// 			trash:                                     &models.Trash{Id: "321", CreatedAt: time.Now(), Location: models.Point{0, 0}},
// 			collection:                                &models.Collection{Id: "111", CreatedAt: time.Now(), TrashId: "321", EventId: "123", Weight: 32},
// 		},
// 	}

// 	//create user, his society add member and applicant
// 	//create friend and friend requester
// 	//create trash, event, add two attendants there + society is admin
// 	//create collection with eventId and create trash-event relation
// 	for i, _ := range candidates {
// 		var err error
// 		//first part
// 		candidates[i].admin, err = s.service.UserAccess.CreateUser(candidates[i].admin)
// 		s.Nil(err)
// 		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].admin.Id)
// 		s.Nil(err)
// 		candidates[i].memberAndEventAttendantAndFriendRequester, err = s.service.UserAccess.CreateUser(candidates[i].memberAndEventAttendantAndFriendRequester)
// 		s.Nil(err)
// 		candidates[i].applicantAndEventAttendantAndFriend, err = s.service.UserAccess.CreateUser(candidates[i].applicantAndEventAttendantAndFriend)
// 		s.Nil(err)

// 		application := &models.Applicant{UserId: candidates[i].memberAndEventAttendantAndFriendRequester.Id, SocietyId: candidates[i].society.Id}
// 		_, err = s.service.UserAccess.AddApplicant(application)
// 		s.Nil(err)
// 		_, err = s.service.UserAccess.AcceptApplicant(application.UserId, application.SocietyId)
// 		s.Nil(err)

// 		application = &models.Applicant{UserId: candidates[i].applicantAndEventAttendantAndFriend.Id, SocietyId: candidates[i].society.Id}
// 		_, err = s.service.UserAccess.AddApplicant(application)
// 		s.Nil(err)

// 		//second part
// 		err = s.db.Insert(&models.Friends{User1Id: candidates[i].admin.Id, User2Id: candidates[i].applicantAndEventAttendantAndFriend.Id})
// 		s.Nil(err)
// 		err = s.db.Insert(&models.FriendRequest{User1Id: candidates[i].admin.Id, User2Id: candidates[i].memberAndEventAttendantAndFriendRequester.Id})
// 		s.Nil(err)

// 		//third part
// 		err = s.db.Insert(candidates[i].trash)
// 		s.Nil(err)
// 		err = s.db.Insert(candidates[i].event)
// 		s.Nil(err)
// 		err = s.db.Insert(&models.EventUser{UserId: candidates[i].applicantAndEventAttendantAndFriend.Id, EventId: candidates[i].event.Id, Permission: "viewer"})
// 		s.Nil(err)
// 		err = s.db.Insert(&models.EventUser{UserId: candidates[i].memberAndEventAttendantAndFriendRequester.Id, EventId: candidates[i].event.Id, Permission: "viewer"})
// 		s.Nil(err)
// 		err = s.db.Insert(&models.EventSociety{SocietyId: candidates[i].society.Id, EventId: candidates[i].event.Id, Permission: "creator"})
// 		s.Nil(err)

// 		//fourth part
// 		candidates[i].collection.EventId = candidates[i].event.Id
// 		candidates[i].collection.TrashId = candidates[i].trash.Id

// 		err = s.db.Insert(candidates[i].collection)
// 		s.Nil(err)
// 		err = s.db.Insert(&models.EventTrash{TrashId: candidates[i].trash.Id, EventId: candidates[i].event.Id})
// 		s.Nil(err)
// 	}

// 	for _, candidate := range candidates {
// 		//check user, his society add member and applicant
// 		err := s.db.Model(&models.User{}).Where("id = ?", candidate.admin.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.Society{}).Where("id = ?", candidate.society.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.Applicant{}).Where("user_id = ?", candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.Member{}).Where("user_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.Nil(err)

// 		//check friend and friend requester
// 		err = s.db.Model(&models.Friends{}).Where("user1_id = ? or user2_id = ?", candidate.applicantAndEventAttendantAndFriend.Id, candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.FriendRequest{}).Where("user1_id = ? or user2_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id, candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.Nil(err)

// 		//create trash, event, add two attendants there + society is admin
// 		err = s.db.Model(&models.Trash{}).Where("id = ?", candidate.trash.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.Event{}).Where("id = ?", candidate.event.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.EventSociety{}).Where("society_id = ?", candidate.society.Id).Select()
// 		s.Nil(err)

// 		//check collection with eventId and create trash-event relation
// 		err = s.db.Model(&models.Collection{}).Where("id = ?", candidate.collection.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.EventTrash{}).Where("trash_id = ?", candidate.trash.Id).Select()
// 		s.Nil(err)

// 		req := httptest.NewRequest(echo.POST, "/users/"+candidate.admin.Id, nil)

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := s.e.NewContext(req, rec)

// 		c.Set("userId", candidate.admin.Id)

// 		s.NoError(s.service.DeleteUser(c))
// 		s.Nil(err)

// 		s.db.AddQueryHook(middlewareService.DbMiddleware{})
// 		s.EqualValues(200, rec.Code)

// 		//check user, his society add member and applicant
// 		err = s.db.Model(&models.User{}).Where("id = ?", candidate.admin.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.Society{}).Where("id = ?", candidate.society.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.Applicant{}).Where("user_id = ?", candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.Member{}).Where("user_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.NotNil(err)

// 		//check friend and friend requester
// 		err = s.db.Model(&models.Friends{}).Where("user1_id = ? or user2_id = ?", candidate.applicantAndEventAttendantAndFriend.Id, candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.FriendRequest{}).Where("user1_id = ? or user2_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id, candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.NotNil(err)

// 		//create trash, event, add two attendants there + society is admin
// 		err = s.db.Model(&models.Trash{}).Where("id = ?", candidate.trash.Id).Select()
// 		s.Nil(err)
// 		err = s.db.Model(&models.Event{}).Where("id = ?", candidate.event.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.applicantAndEventAttendantAndFriend.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.EventUser{}).Where("user_id = ?", candidate.memberAndEventAttendantAndFriendRequester.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.EventSociety{}).Where("society_id = ?", candidate.society.Id).Select()
// 		s.NotNil(err)

// 		//check collection with eventId and create trash-event relation
// 		err = s.db.Model(&models.Collection{}).Where("id = ?", candidate.collection.Id).Select()
// 		s.NotNil(err)
// 		err = s.db.Model(&models.EventTrash{}).Where("trash_id = ?", candidate.trash.Id).Select()
// 		s.NotNil(err)
// 	}
// }

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
}

func (s *UserSuite) SetupTest() {
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

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &UserSuite{})
}

func correctIdOrderFriend(expected *models.Friends, tested *models.Friends) {
	order := strings.Compare(expected.User1Id, expected.User2Id)
	if order == 1 {
		tmp := expected.User1Id
		expected.User1Id = expected.User2Id
		expected.User2Id = tmp
	}
	order = strings.Compare(tested.User1Id, tested.User2Id)
	if order == 1 {
		tmp := tested.User1Id
		tested.User1Id = tested.User2Id
		tested.User2Id = tmp
	}
}
func correctIdOrderFriendRequest(expected *models.FriendRequest, tested *models.FriendRequest) {
	order := strings.Compare(expected.User1Id, expected.User2Id)
	if order == 1 {
		tmp := expected.User1Id
		expected.User1Id = expected.User2Id
		expected.User2Id = tmp
	}
	order = strings.Compare(tested.User1Id, tested.User2Id)
	if order == 1 {
		tmp := tested.User1Id
		tested.User1Id = tested.User2Id
		tested.User2Id = tmp
	}
}
