package user

import (
	"context"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	tableName struct{} `pg:"users"json:"-"`
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*User)(nil)

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.Id = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Society struct {
	tableName struct{} `pg:"societies"json:"-"`
	Id        string
	Name      string
	CreatedAt time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*Society)(nil)

func (u *Society) BeforeInsert(ctx context.Context) (context.Context, error) {
	if u.Id == "" {
		u.Id = uuid.NewV4().String()
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Member struct {
	tableName  struct{} `pg:"societies_members"json:"-"`
	UserId     string
	SocietyId  string
	Permission membership
}

type Applicant struct {
	tableName struct{} `pg:"societies_applicants"json:"-"`
	UserId    string
	SocietyId string
	CreatedAt time.Time `pg:"default:now()"`
}

type membership string

type UserGroupRequest struct {
	UserId    string
	SocietyId string
}
