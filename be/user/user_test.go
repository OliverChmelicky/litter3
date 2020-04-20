package user

import (
	"encoding/json"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

//
//
//
//	FRIENDSHIP
//
//
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
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrApplyForFriendship},
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

		s.NoError(s.service.ApplyForFriendship(c))

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

		s.NoError(s.service.ApplyForFriendship(c))

		respErr := &custom_errors.ErrorModel{}
		err = json.Unmarshal(rec.Body.Bytes(), respErr)
		s.Nil(err)
		candidate.err.Message = respErr.Message
		s.EqualValues(candidate.err, respErr)

		//remove application
		req = httptest.NewRequest(echo.DELETE, "/users/remove/friend/"+candidate.peterAsks.Id, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = s.e.NewContext(req, rec)

		c.SetParamNames("unfriendId")
		c.SetParamValues(candidate.heinrich.Id)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.RemoveApplicationForFriendship(c))
		s.EqualValues("", rec.Body.String())

		rq := &models.FriendRequest{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User1Id, candidate.friendship.User2Id).Select()
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

		s.Nil(s.service.ApplyForFriendship(c))

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
		candidates[i].response = &models.Friends{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.request)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/users/friend/accept", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", candidate.heinrich.Id)

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

		c.SetParamNames("unfriendId")
		c.SetParamValues(candidate.heinrich.Id)
		c.Set("userId", candidate.peterAsks.Id)

		s.NoError(s.service.RemoveFriend(c))
		s.EqualValues(http.StatusOK, rec.Code)

		rq := &models.Friends{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User1Id, candidate.friendship.User2Id).Select()
		s.EqualValues(pg.ErrNoRows, err)
	}
}

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
