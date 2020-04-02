package event

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TrashSuite struct {
	suite.Suite
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

	//s.service = CreateService(db)
	s.db = db
	s.userAccess = &user.UserAccess{Db: db}

	s.e = echo.New()
}

//create event --> hard
//create collection from event
//get society events
//attend event
//don`t attend
//preved prava
//delete event --> hard

func (s *TrashSuite) SetupTest() {
	s.Nil(s.db.DropTable((*user.User)(nil), &orm.DropTableOptions{IfExists: true, Cascade: true}))

	s.Nil(s.db.CreateTable((*user.User)(nil), &orm.CreateTableOptions{IfNotExists: true}))
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &TrashSuite{})
}
