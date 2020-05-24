package models

import (
	"context"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	tableName   struct{} `pg:"users"json:"-"`
	Id          string   `pg:",pk"`
	FirstName   string
	LastName    string
	Email       string
	Uid         string
	Avatar      string
	Societies   []Society    `pg:"many2many:societies_members"`
	Collections []Collection `pg:"many2many:users_collections"`
	Events      []EventUser  `pg:"many2many:events_users"`
	CreatedAt   time.Time    `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*User)(nil)

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.Id = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Society struct {
	tableName     struct{} `pg:"societies"json:"-"`
	Id            string   `pg:",pk"`
	Name          string
	Avatar        string
	Description   string
	Users         []User `pg:"many2many:societies_members"`
	Applicants    []User `pg:"many2many:societies_applicants"`
	MemberRights  []Member
	ApplicantsIds []Applicant
	CreatedAt     time.Time `pg:"default:now()"`
}

type SocietyAnswSimple struct {
	Id          string
	Name        string
	Avatar      string
	Description string
	UsersNumb   int
	CreatedAt   time.Time
}

type SocietyPagingAnsw struct {
	Societies []SocietyAnswSimple
	Paging    Paging
}

var _ pg.BeforeInsertHook = (*Society)(nil)

func (u *Society) BeforeInsert(ctx context.Context) (context.Context, error) {
	if u.Id == "" {
		u.Id = uuid.NewV4().String()
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Membership string

type IdMessage struct {
	Id string
}

type IdsMessage struct {
	Ids []string
}

type EmailMessage struct {
	Email string
}

type UserGroupRequest struct {
	UserId    string
	SocietyId string
}

//middlewares are down
type Member struct {
	tableName  struct{} `pg:"societies_members"json:"-"`
	UserId     string   `pg:",pk"`
	SocietyId  string   `pg:",pk"`
	Permission Membership
	CreatedAt  time.Time `pg:"default:now()"`
}

type Applicant struct {
	tableName struct{}  `pg:"societies_applicants"json:"-"`
	UserId    string    `pg:",pk"`
	SocietyId string    `pg:",pk"`
	CreatedAt time.Time `pg:"default:now()"`
}

type Friends struct {
	tableName struct{}  `pg:"friends"json:"-"`
	User1Id   string    `pg:",pk"`
	User2Id   string    `pg:",pk"`
	CreatedAt time.Time `pg:"default:now()"`
}

type FriendRequest struct {
	tableName struct{}  `pg:"friend_requests"json:"-"`
	User1Id   string    `pg:",pk"`
	User2Id   string    `pg:",pk"`
	CreatedAt time.Time `pg:"default:now()"`
}

var _ pg.BeforeInsertHook = (*Member)(nil)

func (u *Member) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Applicant)(nil)

func (u *Applicant) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Friends)(nil)

func (u *Friends) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*Friends)(nil)

func (u *Friends) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeDeleteHook = (*Friends)(nil)

func (u *Friends) BeforeDelete(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeDeleteHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeDelete(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
