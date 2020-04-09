package event

import (
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			updatingEvent: &EventRequest{Date: time.Now(), Publc: true},
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
			updatingEvent: &EventRequest{Date: time.Now(), Publc: true},
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

		for _, x := range candidates[i].updatingTrash {
			newTrash, err := s.trashAccess.CreateTrash(&x)
			s.Nil(err)
			candidates[i].updatingEvent.Trash = append(candidates[i].updatingEvent.Trash, newTrash.Id)
		}

		candidates[i].eventRequest.UserId = usr.Id
		event, err := s.service.eventAccess.CreateEvent(candidates[i].eventRequest)
		s.Nil(err)

		candidates[i].event.Id = event.Id
		candidates[i].event.UsersIds = event.UsersIds
		candidates[i].event.TrashIds = event.TrashIds
	}

	//UPDATE
	for i, candidate := range candidates {
		updatingRequest := &EventRequest{
			Id:        candidate.event.Id,
			UserId:    candidate.updatingEvent.UserId,
			SocietyId: candidate.updatingEvent.SocietyId,
			AsSociety: candidate.updatingEvent.AsSociety,
			Date:      candidate.updatingEvent.Date,
			Publc:     candidate.updatingEvent.Publc,
			Trash:     candidate.updatingEvent.Trash,
		}
		bytes, err := json.Marshal(updatingRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/events/"+candidate.event.Id, strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.creatorUser.Id)

		s.NoError(s.service.UpdateEvent(c))

		resp := &Event{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(len(candidates[i].updatingEvent.Trash), len(resp.TrashIds))

		candidates[i].event.Date = resp.Date
		candidates[i].event.CreatedAt = resp.CreatedAt
		candidates[i].event.TrashIds = resp.TrashIds
		s.EqualValues(candidates[i].event, resp)
	}

	//TODO try ORM GET
	for i, candidate := range candidates {
		req := httptest.NewRequest(echo.GET, "/events/"+candidate.event.Id, nil)

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
		s.EqualValues(len(candidates[i].event.TrashIds), len(resp.TrashIds))

		candidates[i].event.TrashIds = resp.TrashIds
		s.EqualValues(candidates[i].event, resp)
	}
}

func (s *TrashSuite) Test_CreateTrashUser_AttendEvent_CannotAttend() {
	candidates := []struct {
		admin         *user.User
		eventRequest  *EventRequest
		wantsToAttend *user.User
		trash         []trash.Trash

		attendRequest EventPickerRequest
		event         *Event
	}{
		{
			admin:         &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			eventRequest:  &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash:         []trash.Trash{{Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			wantsToAttend: &user.User{Email: "attends@first.com", FirstName: "joshua", LastName: "Bosh"},
			event:         &Event{Date: time.Now(), Publc: true},
		},
		{
			admin:        &user.User{Email: "ja@me.cpe", FirstName: "Big", LastName: "Rocky"},
			eventRequest: &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash: []trash.Trash{
				{Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			wantsToAttend: &user.User{Email: "attends@second.com", FirstName: "joshua", LastName: "Bosh"},
			event:         &Event{Date: time.Now(), Publc: true},
		},
	}

	for i, _ := range candidates {
		admin, err := s.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].admin = admin

		attendee, err := s.userAccess.CreateUser(candidates[i].wantsToAttend)
		s.Nil(err)
		candidates[i].wantsToAttend = attendee

		for _, tr := range candidates[i].trash {
			newTrash, err := s.trashAccess.CreateTrash(&tr)
			s.Nil(err)
			candidates[i].eventRequest.Trash = append(candidates[i].eventRequest.Trash, newTrash.Id)
			candidates[i].event.TrashIds = append(candidates[i].event.TrashIds, newTrash.Id)
		}

		candidates[i].eventRequest.UserId = admin.Id
		candidates[i].event.UsersIds = append(candidates[i].event.UsersIds, admin.Id)
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

		candidates[i].attendRequest.EventId = resp.Id
		candidates[i].attendRequest.PickerId = attendee.Id
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(&candidate.attendRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/event/attend", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.wantsToAttend.Id)

		s.NoError(s.service.AttendEvent(c))
		s.EqualValues(http.StatusCreated, rec.Code)
		//Chcek if record exists
		reality := &EventUser{EventId: candidate.event.Id, UserId: candidate.wantsToAttend.Id}
		err = s.db.Select(reality)
		s.Nil(err)

		expected := &EventUser{EventId: candidate.event.Id, UserId: candidate.wantsToAttend.Id, Permission: eventPermission("viewer"), CreatedAt: reality.CreatedAt}
		s.EqualValues(expected, reality)
		s.Nil(err)

	}

	for _, candidate := range candidates {
		queryParams := "picker=" + candidate.attendRequest.PickerId + "&event=" + candidate.attendRequest.EventId + "&asSociety=" + strconv.FormatBool(candidate.attendRequest.AsSociety)

		req := httptest.NewRequest(echo.DELETE, "/event/cannot/attend?"+queryParams, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.wantsToAttend.Id)

		s.NoError(s.service.CannotAttendEvent(c))
		s.EqualValues(http.StatusOK, rec.Code)

		reality := &EventUser{EventId: candidate.event.Id, UserId: candidate.wantsToAttend.Id}
		err := s.db.Select(reality)
		s.EqualValues(pg.ErrNoRows, err)
	}

}

func (s *TrashSuite) Test_GetSocietyEvents() {
	candidates := []struct {
		admin   *user.User
		society *user.Society

		event   *Event
		request *EventPickerRequest
	}{
		{
			admin:   &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			society: &user.Society{Name: "Olala"},
			event:   &Event{Date: time.Now(), Publc: true, Description: "this is my first description"},
			request: &EventPickerRequest{AsSociety: true},
		},
		{
			admin:   &user.User{Email: "ja@me.cpe", FirstName: "Big", LastName: "Rocky"},
			society: &user.Society{Name: "HAHA"},
			event:   &Event{Date: time.Now(), Publc: true, Description: "this is cool"},
			request: &EventPickerRequest{AsSociety: true},
		},
	}

	for i, _ := range candidates {
		admin, err := s.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].admin = admin

		society, err := s.userAccess.CreateSocietyWithAdmin(candidates[i].society, admin.Id)
		s.Nil(err)
		candidates[i].society = society

		event, err := s.service.eventAccess.CreateEvent(&EventRequest{UserId: admin.Id, AsSociety: true, SocietyId: society.Id, Date: candidates[i].event.Date, Publc: true})
		candidates[i].event = event
		s.Nil(err)

		candidates[i].request.EventId = event.Id
		candidates[i].request.PickerId = society.Id
	}

	for _, candidate := range candidates {
		queryParams := "picker=" + candidate.request.PickerId + "&event=" + candidate.request.EventId + "&asSociety=" + strconv.FormatBool(candidate.request.AsSociety)

		req := httptest.NewRequest(echo.GET, "/events/society?"+queryParams, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		s.NoError(s.service.GetSocietyEvents(c))
		s.EqualValues(http.StatusOK, rec.Code)

		resp := []Event{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		candidate.event.CreatedAt = resp[0].CreatedAt
		candidate.event.SocietiesIds = nil
		candidate.event.Date = resp[0].Date
		s.EqualValues(*candidate.event, resp[0])
	}
}

func (s *TrashSuite) Test_GetUserEvents() {
	candidates := []struct {
		admin   *user.User
		event   *Event
		request *EventPickerRequest
	}{
		{
			admin:   &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			event:   &Event{Date: time.Now(), Publc: true, Description: "this is my first description"},
			request: &EventPickerRequest{AsSociety: true},
		},
		{
			admin:   &user.User{Email: "ja@me.cpe", FirstName: "Big", LastName: "Rocky"},
			event:   &Event{Date: time.Now(), Publc: true, Description: "this is cool"},
			request: &EventPickerRequest{AsSociety: true},
		},
	}

	for i, _ := range candidates {
		admin, err := s.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].admin = admin

		event, err := s.service.eventAccess.CreateEvent(&EventRequest{UserId: admin.Id, Date: candidates[i].event.Date, Publc: true})
		candidates[i].event = event
		s.Nil(err)

		candidates[i].request.EventId = event.Id
		candidates[i].request.PickerId = admin.Id
	}

	for _, candidate := range candidates {
		queryParams := "picker=" + candidate.request.PickerId

		req := httptest.NewRequest(echo.GET, "/events/user?"+queryParams, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		s.NoError(s.service.GetUserEvents(c))
		s.EqualValues(http.StatusOK, rec.Code)

		resp := []Event{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		candidate.event.CreatedAt = resp[0].CreatedAt
		candidate.event.UsersIds = nil
		candidate.event.Date = resp[0].Date
		s.EqualValues(*candidate.event, resp[0])
	}
}

func (s *TrashSuite) Test_CreateCollectionFromEvents() {
	candidates := []struct {
		admin        *user.User
		eventRequest *EventRequest
		trash        []trash.Trash

		eventId          string
		requestOrganized *trash.CreateCollectionOrganizedRequest
		collections      []trash.CreateCollectionRequest
	}{
		{
			admin:        &user.User{Email: "ja@me.cpg", FirstName: "joshua", LastName: "Bosh"},
			eventRequest: &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash:        []trash.Trash{{Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")}},
			collections:  []trash.CreateCollectionRequest{{CleanedTrash: false, Weight: 622.642}},
		},
		{
			admin:        &user.User{Email: "ja@me.cpe", FirstName: "Big", LastName: "Rocky"},
			eventRequest: &EventRequest{AsSociety: false, Date: time.Now(), Publc: true},
			trash: []trash.Trash{
				{Location: trash.Point{20, 30}, Cleaned: false, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
				{Location: trash.Point{50, 16}, Cleaned: true, Size: trash.Size("bag"), Accessibility: trash.Accessibility("car"), TrashType: trash.TrashType("organic")},
			},
			collections: []trash.CreateCollectionRequest{{CleanedTrash: false, Weight: 622.31}, {CleanedTrash: false, Weight: 63.74}},
		},
	}

	for i, _ := range candidates {
		admin, err := s.userAccess.CreateUser(candidates[i].admin)
		s.Nil(err)
		candidates[i].admin = admin

		//I need to have same amount of trash and collection
		for j, tr := range candidates[i].trash {
			newTrash, err := s.trashAccess.CreateTrash(&tr)
			s.Nil(err)
			candidates[i].eventRequest.Trash = append(candidates[i].eventRequest.Trash, newTrash.Id)
			candidates[i].collections[j].TrashId = newTrash.Id
		}

		candidates[i].eventRequest.UserId = admin.Id
		event, err := s.service.eventAccess.CreateEvent(candidates[i].eventRequest)
		s.Nil(err)
		candidates[i].eventId = event.Id
		candidates[i].requestOrganized = &trash.CreateCollectionOrganizedRequest{
			EventId:     event.Id,
			Collections: candidates[i].collections,
		}
	}

	for _, candidate := range candidates {
		bytes, err := json.Marshal(&candidate.requestOrganized)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/event/collections/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.admin.Id)

		s.NoError(s.service.CreateCollections(c))
		s.EqualValues(http.StatusCreated, rec.Code)

		var newCollections []trash.Collection
		err = s.db.Model(&newCollections).Where("event_id = ?", candidate.eventId).Select()
		s.Nil(err)

		s.EqualValues(len(candidate.collections), len(newCollections))

	}

}

//TODO ako to je s updatom tejto kolekcie?  --> update je napicu...ale bude sa dat

//TODO delete event

//preved prava nova featura a nebude to
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
