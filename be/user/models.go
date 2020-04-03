package user

import (
	"context"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type User struct {
	tableName struct{} `pg:"users"json:"-"`
	Id        string   `pg:",pk"`
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
	Id        string   `pg:",pk"`
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

type membership string

type IdMessage struct {
	Id string
}

type UserGroupRequest struct {
	UserId    string
	SocietyId string
}

type Member struct {
	tableName  struct{} `pg:"societies_members"json:"-"`
	UserId     string   `pg:",pk"`
	SocietyId  string   `pg:",pk"`
	Permission membership
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

var _ pg.BeforeInsertHook = (*Friends)(nil)

func (u *Friends) BeforeInsert(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendsStruct(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*Friends)(nil)

func (u *Friends) BeforeUpdate(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendsStruct(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeDeleteHook = (*Friends)(nil)

func (u *Friends) BeforeDelete(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendsStruct(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeInsert(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendRequest(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeUpdate(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendRequest(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

var _ pg.BeforeDeleteHook = (*FriendRequest)(nil)

func (u *FriendRequest) BeforeDelete(ctx context.Context) (context.Context, error) {
	cmp := strings.Compare(u.User1Id, u.User2Id)
	if cmp == 1 {
		swapUsersFriendRequest(u)
	}
	u.CreatedAt = time.Now()
	return ctx, nil
}

func swapUsersFriendsStruct(u *Friends) {
	tmp := u.User1Id
	u.User1Id = u.User2Id
	u.User2Id = tmp
}

func swapUsersFriendRequest(u *FriendRequest) {
	tmp := u.User1Id
	u.User1Id = u.User2Id
	u.User2Id = tmp
}
