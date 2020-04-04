package user

import (
	"encoding/json"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
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

	s.service = CreateService(db)
	s.db = db

	s.e = echo.New()
}

func (s *UserSuite) Test_ApplyForFriendship_RemoveRequest_AllByUser() {
	candidates := []struct {
		heinrich        *User
		peterAsks       *User
		applicationForm *IdMessage
		friendship      *FriendRequest
		err             *custom_errors.ErrorModel
	}{
		{
			heinrich:  &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			peterAsks: &User{FirstName: "Novy", LastName: "Member"},
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrApplyForFriendship},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)
		candidates[i].applicationForm = &IdMessage{Id: candidates[i].heinrich.Id}

		//filling correct answer
		if strings.Compare(candidates[i].heinrich.Id, candidates[i].peterAsks.Id) == -1 {
			candidates[i].friendship = &FriendRequest{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		} else {
			candidates[i].friendship = &FriendRequest{User1Id: candidates[i].peterAsks.Id, User2Id: candidates[i].heinrich.Id}
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

		resp := &FriendRequest{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.friendship.CreatedAt = resp.CreatedAt
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

		rq := &FriendRequest{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User1Id, candidate.friendship.User2Id).Select()
		s.EqualValues(pg.ErrNoRows, err)
	}
}

func (s *UserSuite) Test_RequestFriendshipExistingFriendship() {
	candidates := []struct {
		heinrich        *User
		peterAsks       *User
		applicationForm *IdMessage
		err             *custom_errors.ErrorModel
	}{
		{
			heinrich:  &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "ja@TestApplyFormMembershipExistingMember.com", CreatedAt: time.Now()},
			peterAsks: &User{FirstName: "Novy", LastName: "Member", Email: "blbost@peterAsks.com"},
			err:       &custom_errors.ErrorModel{ErrorType: custom_errors.ErrConflict, Message: "YOU ARE FIENDS ALREADY"},
		},
	}

	var err error
	for i, _ := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)
		candidates[i].applicationForm = &IdMessage{Id: candidates[i].heinrich.Id}

		//creating friendship
		existingFriendship := &Friends{User1Id: candidates[i].peterAsks.Id, User2Id: candidates[i].heinrich.Id}
		correctIdOrderFriend(existingFriendship)

		err := s.service.UserAccess.Db.Insert(existingFriendship)
		s.Nil(err)

		testExistence := new(Friends)
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
		heinrich  *User
		peterAsks *User
		request   *FriendRequest
		response  *Friends
	}{
		{
			heinrich:  &User{FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			peterAsks: &User{FirstName: "Hello", LastName: "Flowup"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)

		//for filling request structure
		candidates[i].request = &FriendRequest{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		_, err = s.service.UserAccess.AddFriendshipRequest(candidates[i].request)
		s.Nil(err)
		//fir filling final answer
		candidates[i].response = &Friends{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		correctIdOrderFriend(candidates[i].response)
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

		resp := &Friends{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidate.response.CreatedAt = resp.CreatedAt
		s.EqualValues(candidate.response, resp)
	}
}

func (s *UserSuite) Test_RemoveFriend() {
	candidates := []struct {
		heinrich   *User
		peterAsks  *User
		friendship *Friends
	}{
		{
			heinrich:  &User{FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
			peterAsks: &User{FirstName: "Hello", LastName: "Flowup"},
		},
	}

	var err error
	for i := range candidates {
		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
		s.Nil(err)
		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
		s.Nil(err)

		//create friendship
		candidates[i].friendship = &Friends{User1Id: candidates[i].heinrich.Id, User2Id: candidates[i].peterAsks.Id}
		correctIdOrderFriend(candidates[i].friendship)
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

		rq := &Friends{}
		err = s.db.Model(rq).Where("user1_id = ? and user2_id = ?", candidate.friendship.User1Id, candidate.friendship.User2Id).Select()
		s.EqualValues(pg.ErrNoRows, err)
	}
}

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
}

func (s *UserSuite) SetupTest() {
	s.Nil(s.db.DropTable((*User)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Society)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*FriendRequest)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Friends)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Applicant)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Member)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))

	s.Nil(s.db.CreateTable((*User)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Society)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*FriendRequest)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Friends)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Applicant)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Member)(nil), &orm.CreateTableOptions{IfNotExists: true}))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &UserSuite{})
}
