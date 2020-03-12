package user

import (
	"github.com/go-pg/pg/v9"
)

type userAccess struct {
	db *pg.DB
}
