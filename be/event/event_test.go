package event

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type TrashSuite struct {
	suite.Suite
	service     *EventService
	userAccess  *user.UserAccess
	trashAccess *trash.TrashAccess
	e           *echo.Echo
	db          *pg.DB
}

func (s *TrashSuite) SetupSuite() {
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

	s.userAccess = &user.UserAccess{Db: db}
	s.trashAccess = &trash.TrashAccess{Db: db}

	s.service = CreateService(db, &user.UserAccess{Db: db}, &trash.TrashAccess{Db: db})
	s.db = db

	s.e = echo.New()
}

//create event --> hard
func (s *TrashSuite) Test_CreateTrash_User() {
	candidates := []struct {
		creatorUser  *user.User
		eventRequest *EventRequest
		trash        []trash.Trash
		event        *Event
	}{
		{
			creatorUser:  &user.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			eventRequest: &EventRequest{UserId: "1", AsSociety: false, Date: time.Now(), Publc: true},
			trash:        []trash.Trash{{Id: "1", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			event:        &Event{Date: time.Now(), Publc: true},
		},
		{
			creatorUser:  &user.User{Id: "2", FirstName: "Damian", LastName: "Zelenina", Email: "On@friend.com", CreatedAt: time.Now()},
			eventRequest: &EventRequest{UserId: "2", AsSociety: false, Date: time.Now(), Publc: true},
			trash: []trash.Trash{
				{Id: "9", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Id: "2", Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			event: &Event{Date: time.Now(), Publc: true},
		},
	}

	for i, _ := range candidates {
		usr, err := s.userAccess.CreateUser(candidates[i].creatorUser)
		s.Nil(err)
		candidates[i].creatorUser = usr
		candidates[i].event.UsersIds = append(candidates[i].event.UsersIds, usr.Id)

		for _, x := range candidates[i].trash {
			newTrash, err := s.trashAccess.CreateTrash(&x)
			s.Nil(err)
			candidates[i].eventRequest.Trash = append(candidates[i].eventRequest.Trash, newTrash.Id)
			candidates[i].event.TrashIds = append(candidates[i].event.TrashIds, newTrash.Id)
		}

		bytes, err := json.Marshal(candidates[i].eventRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/event/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creatorUser.Id)

		s.NoError(s.service.CreateEvent(c))

		resp := &Event{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		fmt.Println(rec.Body.String())

		candidates[i].event.Id = resp.Id
		candidates[i].event.Date = resp.Date
		candidates[i].event.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].event, resp)
	}

	//TODO
	//test update trash
	//for _, candidate := range candidates {
	//	bytes, err := json.Marshal(candidate.updating)
	//	s.Nil(err)
	//
	//	req := httptest.NewRequest(echo.PUT, "/trash/update", strings.NewReader(string(bytes)))
	//
	//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//	rec := httptest.NewRecorder()
	//	c := s.e.NewContext(req, rec)
	//	c.Set("userId", candidate.creator.Id)
	//
	//	s.NoError(s.service.UpdateTrash(c))
	//
	//	resp := &Trash{}
	//	err = json.Unmarshal(rec.Body.Bytes(), resp)
	//	s.Nil(err)
	//
	//	s.EqualValues(candidate.updating, resp)
	//
	//}
}

func (s *TrashSuite) Test_CreateTrash_Society() {

}

//getEvent
func (s *TrashSuite) Test_GetEvent() {

}

//attend event
//don`t attend
//preved prava
//get society events
//get user events
//create collection from event
//delete event --> hard

func (s *TrashSuite) SetupTest() {
	var tableInfo []struct {
		Table string
	}
	query := `SELECT table_name "table"
				FROM information_schema.tables WHERE table_schema='public'
					AND table_type='BASE TABLE' AND table_name!= 'gopg_migrations';`
	_, err := s.db.Query(&tableInfo, query)
	if err != nil {
		log.Error(err)
		return
	}

	truncateQueries := make([]string, len(tableInfo))

	for i, info := range tableInfo {
		if info.Table == "spatial_ref_sys" { //postgis extension
			continue
		}
		truncateQueries[i] = "TRUNCATE " + info.Table + " CASCADE;"
	}

	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		for _, query := range truncateQueries {
			_, err = tx.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &TrashSuite{})
}
