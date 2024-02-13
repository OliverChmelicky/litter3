package trash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
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
		creator  *models.User
		trash    *models.Trash
		updating *models.Trash
	}{
		{
			creator:  &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "4f6f", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:    &models.Trash{Location: models.Point{20, 30}, Cleaned: false, Size: models.Size("bag"), Accessibility: models.Accessibility("car"), TrashType: models.TrashType(1)},
			updating: &models.Trash{Location: models.Point{99, 69}, Cleaned: true, Size: models.Size("bag"), Accessibility: models.Accessibility("unknown"), TrashType: models.TrashType(1)},
		},
	}

	for i, _ := range candidates {
		user, err := s.userAccess.CreateUser(candidates[i].creator)
		candidates[i].creator = user
		candidates[i].trash.FinderId = user.Id
		s.NoError(err)

		bytes, err := json.Marshal(candidates[i].trash)
		s.Nil(err)

		req := httptest.NewRequest(echo.POST, "/trash/new", strings.NewReader(string(bytes)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creator.Id)

		s.NoError(s.service.CreateTrash(c))

		resp := &models.Trash{}
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

		resp := &models.Trash{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		s.EqualValues(candidate.updating, resp)

	}
}

// func (s *TrashSuite) Test_GetAround() {
// 	candidates := []struct {
// 		creator           *models.User
// 		trash             *models.Trash
// 		rangeRequest      *models.RangeRequest
// 		collectionRequest *models.CreateCollectionRandomRequest
// 		randomImg         *models.TrashImage
// 	}{
// 		{
// 			creator:           &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "Velikonoce", Email: "Ja@kamo.com", CreatedAt: time.Now()},
// 			trash:             &models.Trash{Location: models.Point{20, 30}, Cleaned: false, Size: models.Size("bag"), Accessibility: models.Accessibility("car"), TrashType: models.TrashType(1)},
// 			rangeRequest:      &models.RangeRequest{Location: models.Point{20, 29.99}, Radius: 5000.0},
// 			collectionRequest: &models.CreateCollectionRandomRequest{Weight: 32},
// 			randomImg:         &models.TrashImage{Url: "dasd"},
// 		},
// 	}

// 	for i, _ := range candidates {
// 		var err error
// 		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
// 		s.Nil(err)
// 		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
// 		s.Nil(err)
// 		candidates[i].collectionRequest.TrashId = candidates[i].trash.Id
// 		collection, err := s.service.CreateCollectionRandom(candidates[i].collectionRequest, candidates[i].creator.Id)
// 		s.Nil(err)
// 		candidates[i].randomImg.TrashId = candidates[i].trash.Id
// 		err = s.db.Insert(candidates[i].randomImg)
// 		s.Nil(err)

// 		bytes, err := json.Marshal(candidates[i].rangeRequest)
// 		s.Nil(err)
// 		queryParams :=
// 			"?lng=" + fmt.Sprintf("%f", candidates[i].rangeRequest.Location[0]) +
// 				"&lat=" + fmt.Sprintf("%f", candidates[i].rangeRequest.Location[1]) +
// 				"&radius=" + fmt.Sprintf("%f", candidates[i].rangeRequest.Radius)

// 		req := httptest.NewRequest(echo.POST, "/trash/get/range"+queryParams, strings.NewReader(string(bytes)))

// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := s.e.NewContext(req, rec)

// 		s.NoError(s.service.GetTrashInRange(c))

// 		var resp []models.Trash
// 		err = json.Unmarshal(rec.Body.Bytes(), &resp)
// 		s.Nil(err)

// 		fmt.Println("Len of resp: ", len(resp))

// 		s.EqualValues(candidates[i].trash.Location, resp[0].Location)
// 		collection.CreatedAt = resp[0].Collections[0].CreatedAt
// 		//		s.EqualValues(*collection, resp[0].Collections[0])
// 		//		s.EqualValues(*candidates[i].randomImg, resp[0].Images[0])
// 	}
// }

func (s *TrashSuite) Test_CreateCommentOnTrash() {
	candidates := []struct {
		creator         *models.User
		trash           *models.Trash
		commentRequest  *models.TrashCommentRequest
		updatingRequest *models.TrashCommentRequest
		actualComment   *models.TrashComment
		trashComments   []models.TrashComment
	}{
		{
			creator:         &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "8848", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:           &models.Trash{Location: models.Point{20, 30}},
			commentRequest:  &models.TrashCommentRequest{Message: "prva message"},
			actualComment:   &models.TrashComment{Message: "prva message"},
			updatingRequest: &models.TrashCommentRequest{Message: "DRUHA message"},
			trashComments:   []models.TrashComment{{Message: "DRUHA message"}},
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

		resp := &models.TrashComment{}
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

		resp := &models.TrashComment{}
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

		s.NoError(s.service.GetTrashCommentsByTrashId(c))

		var resp []models.TrashComment
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)
		arr := candidate.trashComments

		s.EqualValues(len(candidate.trashComments), len(resp))
		for i, comment := range resp {
			s.EqualValues(arr[i].Message, comment.Message)
		}
	}

	//test delete trashComent
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
		creator    *models.User
		trash      *models.Trash
		request    *models.CreateCollectionRandomRequest
		collection *models.Collection
	}{
		{
			creator:    &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "sds", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:      &models.Trash{Location: models.Point{20, 30}},
			request:    &models.CreateCollectionRandomRequest{Weight: 369.7, CleanedTrash: true},
			collection: &models.Collection{CleanedTrash: true, Weight: 369.7},
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
		candidates[i].trash.Cleaned = true

		resp := &models.Collection{}
		err = json.Unmarshal(rec.Body.Bytes(), resp)
		s.Nil(err)

		candidates[i].collection.Id = resp.Id
		candidates[i].collection.CreatedAt = resp.CreatedAt
		s.EqualValues(candidates[i].collection, resp)

		trash, err := s.service.GetTrash(candidates[i].trash.Id)
		s.Nil(err)
		s.EqualValues(candidates[i].trash.Cleaned, trash.Cleaned)
	}
}

func (s *TrashSuite) Test_DeleteCollection() {
	candidates := []struct {
		creator    *models.User
		trash      *models.Trash
		request    *models.CreateCollectionRandomRequest
		collection *models.Collection
		friends    []*models.User
	}{
		{
			creator:    &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "sds", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:      &models.Trash{Location: models.Point{20, 30}},
			request:    &models.CreateCollectionRandomRequest{Weight: 369.7, CleanedTrash: true},
			collection: &models.Collection{},
		},
		{
			creator:    &models.User{FirstName: "Miro", LastName: "Motyka", Uid: "sdsw", Email: "Ja@kamo.comsa"},
			trash:      &models.Trash{Location: models.Point{20, 30}},
			request:    &models.CreateCollectionRandomRequest{Weight: 369.7, CleanedTrash: true},
			collection: &models.Collection{Users: []models.User{}, CreatedAt: time.Now()},
			friends: []*models.User{
				{FirstName: "Niekto", LastName: "Novy", Uid: "me", Email: "Ja@kamo.in"},
			},
		},
	}

	for i, cand := range candidates {
		var err error
		var friendsIds []string

		for j, friend := range cand.friends {
			candidates[i].friends[j], err = s.userAccess.CreateUser(friend)
			s.Nil(err)
			friendsIds = append(friendsIds, candidates[i].friends[j].Id)
		}
		candidates[i].request.Friends = friendsIds

		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)
		candidates[i].request.TrashId = candidates[i].trash.Id
		candidates[i].collection, err = s.service.TrashAccess.CreateCollectionRandom(candidates[i].request, candidates[i].creator.Id)
		s.Nil(err)

		req := httptest.NewRequest(echo.DELETE, "/trash/collection/"+candidates[i].collection.Id, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := s.e.NewContext(req, rec)
		c.Set("userId", candidates[i].creator.Id)
		c.SetParamNames("collectionId")
		c.SetParamValues(candidates[i].collection.Id)

		s.NoError(s.service.DeleteCollectionFromUser(c))

		s.EqualValues(http.StatusOK, rec.Code)

		if len(candidates[i].friends) == 0 {
			collection, err := s.service.TrashAccess.GetCollection(candidates[i].collection.Id)
			s.NotNil(err)
			s.EqualValues(&models.Collection{}, collection)
		} else {
			collection, err := s.service.TrashAccess.GetCollection(candidates[i].collection.Id)
			s.Nil(err)
			candidates[i].collection.CreatedAt = collection.CreatedAt
			candidates[i].collection.Users = []models.User{
				{
					Id:        collection.Users[0].Id,
					FirstName: candidates[i].friends[0].FirstName,
					LastName:  candidates[i].friends[0].LastName,
					Email:     candidates[i].friends[0].Email,
					Uid:       candidates[i].friends[0].Uid,
					Avatar:    candidates[i].friends[0].Avatar,
					CreatedAt: collection.Users[0].CreatedAt,
				},
			}
			s.EqualValues(candidates[i].collection, collection)
		}
	}
}

func (s *TrashSuite) Test_DeleteTrashWithComments() {
	candidates := []struct {
		creator *models.User
		trash   *models.Trash
		comment *models.TrashComment
	}{
		{
			creator: &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "6ads", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:   &models.Trash{Location: models.Point{20, 30}},
			comment: &models.TrashComment{Message: "prva message"},
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)

		candidates[i].comment.TrashId = candidates[i].trash.Id
		candidates[i].comment.UserId = candidates[i].creator.Id
		candidates[i].comment, err = s.service.TrashAccess.CreateTrashComment(candidates[i].comment)
		s.Nil(err)

		//TODO add image
	}

	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.DELETE, "/trash/"+candidate.trash.Id, nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.SetParamNames("trashId")
		c.SetParamValues(candidate.trash.Id)

		s.NoError(s.service.DeleteTrash(c))

		s.EqualValues(http.StatusOK, rec.Code)

		err := s.db.Select(candidate.trash)
		s.EqualValues(pg.ErrNoRows, err)
	}
}

func (s *TrashSuite) Test_DeleteTrashWithCollection() {
	candidates := []struct {
		creator    *models.User
		trash      *models.Trash
		collection *models.CreateCollectionRandomRequest
		error      error
		code       int
	}{
		{
			creator:    &models.User{Id: "1", FirstName: "Jano", LastName: "Motyka", Uid: "6ads", Email: "Ja@kamo.com", CreatedAt: time.Now()},
			trash:      &models.Trash{Location: models.Point{20, 30}},
			collection: &models.CreateCollectionRandomRequest{Weight: 694},
			error:      fmt.Errorf("Error trash has some collections already "),
			code:       http.StatusInternalServerError,
		},
	}

	for i, _ := range candidates {
		var err error
		candidates[i].creator, err = s.userAccess.CreateUser(candidates[i].creator)
		s.Nil(err)
		candidates[i].trash, err = s.service.TrashAccess.CreateTrash(candidates[i].trash)
		s.Nil(err)

		candidates[i].collection.TrashId = candidates[i].trash.Id
		_, err = s.service.TrashAccess.CreateCollectionRandom(candidates[i].collection, candidates[i].creator.Id)
		s.Nil(err)

		//TODO add image
	}

	for _, candidate := range candidates {
		req := httptest.NewRequest(echo.PUT, "/trash/collection/update", nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.SetParamNames("trashId")
		c.SetParamValues(candidate.trash.Id)

		s.NoError(s.service.DeleteTrash(c))

		var resp custom_errors.ErrorModel
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Nil(err)

		s.EqualValues(candidate.code, rec.Code)
		s.EqualValues(candidate.error.Error(), resp.Message)
	}
}

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

func fillActualComment(comment *models.TrashComment, resp *models.TrashComment, creatorId, trashId string) {
	comment.UserId = creatorId
	comment.TrashId = trashId
	comment.CreatedAt = resp.CreatedAt
	//comment.UpdatedAt = resp.UpdatedAt
	comment.Id = resp.Id
}
