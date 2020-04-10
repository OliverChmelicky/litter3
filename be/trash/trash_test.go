package trash

import (
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
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
	service    *trashService
	userAccess *user.UserAccess
	e          *echo.Echo
	db         *pg.DB
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

	s.service = CreateService(db)
	s.db = db
	s.userAccess = &user.UserAccess{Db: db}

	s.e = echo.New()
}

func (s *TrashSuite) Test_CreateTrash() {
	candidates := []struct {
		creator  *user.User
		trash    *Trash
		updating *Trash
	}{
		{
			creator:  &user.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:    &Trash{Location: Point{20, 30}, Cleaned: false, Size: Size("bag"), Accessibility: Accessibility("car"), TrashType: TrashType("organic")},
			updating: &Trash{Location: Point{99, 69}, Cleaned: true, Size: Size("bag"), Accessibility: Accessibility("unknown"), TrashType: TrashType("organic")},
		},
	}

	for i, _ := range candidates {
		user, err := s.userAccess.CreateUser(candidates[i].creator)
		candidates[i].creator = user
		s.NoError(err)

		bytes, err := json.Marshal(candidates[i].trash)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/trash/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creator.Id)

		s.NoError(s.service.CreateTrash(c))

		resp := &Trash{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].trash.Id = resp.Id
		candidates[i].trash.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].trash, resp)
		candidates[i].updating.Id = resp.Id
		candidates[i].updating.CreatedAt = resp.CreatedAt
	}

	//test update trash
	for _, candidate := range candidates {
		bytes, err := json.Marshal(candidate.updating)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/trash/update", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.creator.Id)

		s.NoError(s.service.UpdateTrash(c))

		resp := &Trash{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.updating, resp)

	}
}

func (s *TrashSuite) Test_GetAround() {
	candidates := []struct {
		creator      *user.User
		trash        *Trash
		rangeRequest *RangeRequest
	}{
		{
			creator:      &user.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:        &Trash{Location: Point{20, 30}, Cleaned: false, Size: Size("bag"), Accessibility: Accessibility("car"), TrashType: TrashType("organic")},
			rangeRequest: &RangeRequest{Location: Point{20, 29.99}, Radius: 5000.0},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)

		bytes, err := json.Marshal(candidates[i].rangeRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/trash/get/range", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		s.NoError(s.service.GetTrashInRange(c))

		var resp []Trash
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		s.EqualValues(candidates[i].trash.Location, resp[0].Location)
	}
}

func (s *TrashSuite) Test_CreateCommentOnTrash() {
	candidates := []struct {
		creator         *user.User
		trash           *Trash
		commentRequest  *TrashCommentRequest
		updatingRequest *TrashCommentRequest
		actualComment   *TrashComment
		trashComments   []TrashComment
	}{
		{
			creator:         &user.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:           &Trash{Location: Point{20, 30}},
			commentRequest:  &TrashCommentRequest{Message: "prva message"},
			actualComment:   &TrashComment{Message: "prva message"},
			updatingRequest: &TrashCommentRequest{Message: "DRUHA message"},
			trashComments:   []TrashComment{{Message: "DRUHA message"}},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)

		candidates[i].commentRequest.Id = candidates[i].trash.Id

		bytes, err := json.Marshal(candidates[i].commentRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/trash/comment/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creator.Id)

		s.NoError(s.service.CreateTrashComment(c))

		resp := &TrashComment{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		fillActualComment(candidates[i].actualComment, resp, candidates[i].creator.Id, candidates[i].trash.Id)
		s.EqualValues(candidates[i].actualComment, resp)
	}

	//test update trash
	for i, candidate := range candidates {
		candidates[i].updatingRequest.Id = candidates[i].actualComment.Id
		bytes, err := json.Marshal(candidate.updatingRequest)
		s.Nil(err)

		req := httptest.NewRequest(echo.PUT, "/trash/comment/update", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidate.creator.Id)

		s.NoError(s.service.UpdateTrashComment(c))

		resp := &TrashComment{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		fillActualComment(candidates[i].actualComment, resp, candidates[i].creator.Id, candidates[i].trash.Id)
		candidates[i].actualComment.Message = candidate.updatingRequest.Message
		s.EqualValues(candidates[i].actualComment, resp)
		candidates[i].updatingRequest.Id = resp.Id
	}

	//test get trash
	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.GET, "/trash/comment/", nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.SetParamNames("trashId")
		c.SetParamValues(candidate.trash.Id)

		s.NoError(s.service.GetTrashComments(c))

		var resp []TrashComment
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)
		arr := candidate.trashComments

		s.EqualValues(len(candidate.trashComments), len(resp))
		for i, comment := range resp {
			s.EqualValues(arr[i].Message, comment.Message)
		}
	}

	//test delete trash
	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.DELETE, "/trash/comment/", nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)

		c.Set("userId", candidate.creator.Id)

		c.SetParamNames("commentId")
		c.SetParamValues(candidate.actualComment.Id)

		s.NoError(s.service.DeleteTrashComment(c))

		s.EqualValues("", rec.Body.String())
	}
}

func (s *TrashSuite) Test_CreateCollectionRandom_GetCollection() {
	candidates := []struct {
		creator    *user.User
		trash      *Trash
		request    *CreateCollectionRandomRequest
		collection *Collection
	}{
		{
			creator:    &user.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:      &Trash{Location: Point{20, 30}},
			request:    &CreateCollectionRandomRequest{Weight: 369.7, CleanedTrash: false},
			collection: &Collection{CleanedTrash: false, Weight: 369.7},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)

		candidates[i].request.TrashId = candidates[i].trash.Id

		candidates[i].collection.TrashId = candidates[i].trash.Id

		bytes, err := json.Marshal(candidates[i].request)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/trash/collection/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creator.Id)

		s.NoError(s.service.CreateCollection(c))

		resp := &Collection{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].collection.Id = resp.Id
		candidates[i].collection.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].collection, resp)
	}
}

//TODO test get collections of user

func (s *TrashSuite) SetupTest() {
	referencerTables := []string{
		"users",
		"trash",
		"collections",
		"users_collections",
		"trash_comments",
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

func fillActualComment(comment *TrashComment, resp *TrashComment, creatorId, trashId string) {
	comment.UserId = creatorId
	comment.TrashId = trashId
	comment.CreatedAt = resp.CreatedAt
	comment.UpdatedAt = resp.UpdatedAt
	comment.Id = resp.Id
}
