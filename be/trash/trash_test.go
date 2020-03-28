package trash

import (
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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
			trash:    &Trash{Location: Point{20, 30}},
			updating: &Trash{Location: Point{99, 69}},
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
		log.Info("a trash je")
		log.Info(resp)
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

func (s *TrashSuite) SetupTest() {
	s.Nil(s.db.DropTable((*user.User)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))
	s.Nil(s.db.DropTable((*Trash)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))

	s.Nil(s.db.CreateTable((*user.User)(nil), &orm.CreateTableOptions{IfNotExists: true}))
	s.Nil(s.db.CreateTable((*Trash)(nil), &orm.CreateTableOptions{IfNotExists: true}))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &TrashSuite{})
}
