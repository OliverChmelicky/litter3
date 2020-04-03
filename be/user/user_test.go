package user

import (
	"encoding/json"
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
		err             string
	}{
		{
			heinrich:  &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@Janovutbr.cz", CreatedAt: time.Now()},
			peterAsks: &User{FirstName: "Novy", LastName: "Member"},
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

		req := httptest.NewRequest(echo.POST, "/users/friend", strings.NewReader(string(bytes)))

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
	}
}

//
//func (s *UserSuite) Test_RequestFriendshipExistingFriendship() {
//	candidates := []struct {
//		heinrich        *User
//		society         *Society
//		peterAsks       *User
//		applicationForm *UserGroupRequest
//		finalApplicant  *Applicant
//		err             string
//	}{
//		{
//			heinrich:  &User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "ja@TestApplyFormMembershipExistingMember.com", CreatedAt: time.Now()},
//			society:   &Society{Name: "TestApplyFormMembershipExistingMember"},
//			peterAsks: &User{FirstName: "Novy", LastName: "Member", Email: "blbost@peterAsks.com"},
//		},
//	}
//
//	var err error
//	for i, _ := range candidates {
//		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
//		s.Nil(err)
//		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
//		s.Nil(err)
//
//		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].heinrich.Id)
//		s.Nil(err)
//
//		peterAsks := &Member{UserId: candidates[i].peterAsks.Id, SocietyId: candidates[i].society.Id, Permission: membership("member")}
//		err := s.service.UserAccess.Db.Insert(peterAsks)
//		s.Nil(err)
//
//		testExistence := new(Member)
//		err = s.db.Model(testExistence).Where("user_id = ?", candidates[i].peterAsks.Id).Select()
//		if err != nil {
//			s.Nil(err) //end test
//		}
//		if errors.Is(err, pg.ErrNoRows) {
//			fmt.Println("Should be found something")
//			s.Error(nil) //throw error in test
//		}
//
//		//for filling request structure
//		candidates[i].applicationForm = &UserGroupRequest{UserId: candidates[i].peterAsks.Id, SocietyId: candidates[i].society.Id}
//		candidates[i].finalApplicant = &Applicant{UserId: candidates[i].peterAsks.Id, SocietyId: candidates[i].society.Id}
//	}
//
//	for _, candidate := range candidates {
//		bytes, err := json.Marshal(candidate.applicationForm)
//		s.Nil(err)
//
//		req := httptest.NewRequest(echo.POST, "/societies/apply", strings.NewReader(string(bytes)))
//
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//		rec := httptest.NewRecorder()
//		c := s.e.NewContext(req, rec)
//		c.Set("userId", candidate.peterAsks.Id)
//
//		s.Nil(s.service.ApplyForMembership(c))
//
//		resp := &Applicant{}
//		err = json.Unmarshal(rec.Body.Bytes(), resp)
//		s.NotNil(err)
//
//		s.EqualValues("User is already a member", rec.Body.String())
//	}
//}
//
//func (s *UserSuite) Test_DismissFriendshipRequest() {
//	candidates := []struct {
//		heinrich    *User
//		society     *Society
//		peterAsks   *User
//		application *Applicant
//	}{
//		{
//			heinrich:  &User{Id: "2", FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
//			society:   &Society{Name: "More members than one"},
//			peterAsks: &User{FirstName: "Hello", LastName: "Flowup"},
//		},
//	}
//
//	var err error
//	for i, _ := range candidates {
//		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
//		s.Nil(err)
//		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
//		s.Nil(err)
//
//		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].heinrich.Id)
//		s.Nil(err)
//
//		//for filling request structure
//		candidates[i].application = &Applicant{UserId: candidates[i].peterAsks.Id, SocietyId: candidates[i].society.Id}
//		_, err = s.service.UserAccess.AddApplicant(candidates[i].application)
//		s.Nil(err)
//	}
//
//	for _, cand := range candidates {
//		req := httptest.NewRequest(echo.DELETE, "/societies/dismiss/"+cand.society.Id+"/"+cand.peterAsks.Id, nil)
//
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//		rec := httptest.NewRecorder()
//		c := s.e.NewContext(req, rec)
//
//		c.Set("userId", cand.heinrich.Id)
//
//		c.SetPath("/societies/dismiss/:societyId/:userId")
//		c.SetParamNames("societyId", "userId")
//		c.SetParamValues(cand.society.Id, cand.peterAsks.Id)
//
//		s.NoError(s.service.DismissApplicant(c))
//
//		s.EqualValues("", rec.Body.String())
//	}
//
//}
//
//func (s *UserSuite) Test_AddFriend() {
//	candidates := []struct {
//		heinrich    *User
//		society     *Society
//		peterAsks   *User
//		application *Applicant
//		response    *Member
//	}{
//		{
//			heinrich:  &User{FirstName: "John", LastName: "Modest", Email: "Ja@Janovutbr.cz"},
//			society:   &Society{Name: "More members than one"},
//			peterAsks: &User{FirstName: "Hello", LastName: "Flowup"},
//			response:  &Member{Permission: "member"},
//		},
//	}
//
//	var err error
//	for i := range candidates {
//		candidates[i].heinrich, err = s.service.UserAccess.CreateUser(candidates[i].heinrich)
//		s.Nil(err)
//		candidates[i].peterAsks, err = s.service.UserAccess.CreateUser(candidates[i].peterAsks)
//		s.Nil(err)
//
//		candidates[i].society, err = s.service.UserAccess.CreateSocietyWithAdmin(candidates[i].society, candidates[i].heinrich.Id)
//		s.Nil(err)
//		candidates[i].response.UserId = candidates[i].peterAsks.Id
//		candidates[i].response.SocietyId = candidates[i].society.Id
//
//		//for filling request structure
//		candidates[i].application = &Applicant{UserId: candidates[i].peterAsks.Id, SocietyId: candidates[i].society.Id}
//		_, err = s.service.UserAccess.AddApplicant(candidates[i].application)
//		s.Nil(err)
//	}
//
//	for _, candidate := range candidates {
//		req := httptest.NewRequest(echo.POST, "/societies/"+candidate.application.SocietyId+"/"+candidate.application.UserId+"/approve", nil)
//
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//		rec := httptest.NewRecorder()
//		c := s.e.NewContext(req, rec)
//
//		c.Set("userId", candidate.peterAsks.Id)
//		c.SetParamNames("societyId", "userId")
//		c.SetParamValues(candidate.society.Id, candidate.peterAsks.Id)
//
//		s.Nil(s.service.AcceptApplicant(c))
//
//		resp := &Member{}
//		err = json.Unmarshal(rec.Body.Bytes(), resp)
//		s.Nil(err)
//
//		s.EqualValues(candidate.response, resp)
//	}
//}

func (s *UserSuite) TearDownSuite() {
	s.db.Close()
}

func (s *UserSuite) SetupTest() {
	s.Nil(s.db.DropTable((*User)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Society)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*FriendRequest)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Friends)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))

	s.Nil(s.db.CreateTable((*User)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Society)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*FriendRequest)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Friends)(nil), &orm.CreateTableOptions{IfNotExists: true}))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &UserSuite{})
}
