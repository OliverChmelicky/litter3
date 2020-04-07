package event

import (
	"encoding/json"
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

		candidates[i].event.Id = resp.Id
		candidates[i].event.Date = resp.Date
		candidates[i].event.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].event, resp)
	}

}

func (s *TrashSuite) Test_CreateTrash_Society() {
	candidates := []struct {
		admin          *user.User
		creatorSociety *user.Society
		eventRequest   *EventRequest
		trash          []trash.Trash
		event          *Event
	}{
		{
			admin:          &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			creatorSociety: &user.Society{Name: "Original", CreatedAt: time.Now()},
			eventRequest:   &EventRequest{AsSociety: true, Date: time.Now(), Publc: true},
			trash:          []trash.Trash{{Id: "1", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			event:          &Event{Date: time.Now(), Publc: true},
		},
		{
			admin:          &user.User{Email: "ja@me.cpe", FirstName: "Big", LastName: "Rocky"},
			creatorSociety: &user.Society{Name: "company", CreatedAt: time.Now()},
			eventRequest:   &EventRequest{AsSociety: true, Date: time.Now(), Publc: true},
			trash: []trash.Trash{
				{Id: "9", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Id: "2", Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			event: &Event{Date: time.Now(), Publc: true},
		},
	}

	for i, _ := range candidates {
		admin, err := s.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		society, err := s.userAccess.CreateSocietyWithAdmin(candidates[i].creatorSociety, admin.Id)
		s.Nil(err)
		candidates[i].creatorSociety = society
		candidates[i].eventRequest.SocietyId = society.Id
		candidates[i].event.SocietiesIds = append(candidates[i].event.SocietiesIds, society.Id)

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
		c.Set("userId", candidates[i].admin.Id)

		s.NoError(s.service.CreateEvent(c))

		resp := &Event{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].event.Id = resp.Id
		candidates[i].event.Date = resp.Date
		candidates[i].event.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].event, resp)
	}
}

//getEvent
func (s *TrashSuite) Test_GetEvent_UpdateEvent() {
	candidates := []struct {
		creatorUser   *user.User
		eventRequest  *EventRequest
		trash         []trash.Trash
		updatingTrash []trash.Trash
		event         *Event
		updatingEvent *EventRequest
	}{
		{
			creatorUser:  &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			eventRequest: &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash:        []trash.Trash{{Id: "1", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			updatingTrash: []trash.Trash{
				{Id: "9", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Id: "2", Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			event:         &Event{Date: time.Now(), Publc: true},
			updatingEvent: &EventRequest{Date: time.Now(), Publc: false},
		},
		{
			creatorUser:  &user.User{Email: "ja@he.cpe", FirstName: "Big", LastName: "Rocky"},
			eventRequest: &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash: []trash.Trash{
				{Id: "9", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Id: "2", Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			updatingTrash: []trash.Trash{{Id: "1", Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			event:         &Event{Date: time.Now(), Publc: true},
			updatingEvent: &EventRequest{Date: time.Now(), Publc: false},
		},
	}

	//create
	for i := range candidates {
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

		candidates[i].eventRequest.UserId = usr.Id
		event, err := s.service.eventAccess.CreateEvent(candidates[i].eventRequest)
		s.Nil(err)

		candidates[i].event.Id = event.Id
		candidates[i].event.UsersIds = event.UsersIds
		candidates[i].event.TrashIds = event.TrashIds
	}
	//TODO update

	//get
	for i, candidate := range candidates {
		req := httptest.NewRequest(echo.GET, "/events"+candidate.event.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetParamNames("eventId")
		c.SetParamValues(candidate.event.Id)

		s.NoError(s.service.GetEvent(c))

		resp := &Event{}
		err := json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].event.Date = resp.Date
		candidates[i].event.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].event, resp)
	}
}

//attend event
//don`t attend
//preved prava
//get society events
//get user events
//create collection from event
//delete event --> hard

func (s *TrashSuite) SetupTest() {
	referencerTables := []string{
		"events_societies",
		"events_trash",
		"events_users",
		"societies_members",
		"users",
		"events",
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
	suite.Run(t, &TrashSuite{})
}
