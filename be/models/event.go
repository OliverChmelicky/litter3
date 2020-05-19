package models

import (
	"context"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EventRequest struct {
	Id          string
	UserId      string
	SocietyId   string
	AsSociety   bool
	Description string
	Date        time.Time
	Trash       []string
}

type EventPickerRequest struct {
	PickerId  string
	EventId   string
	AsSociety bool
}

type GetEventsRequest struct {
	PickerId string
	Paging   int
}

type EventPermissionRequest struct {
	ChangingRightsTo string
	EventId          string
	Permission       EventPermission
	AsSociety        bool
	SocietyId        string
}

type CreateEvent struct {
	tableName   struct{} `pg:"events"json:"-"`
	Id          string   `pg:",pk"`
	Date        time.Time
	Description string
	CreatedAt   time.Time `pg:"default:now()"`
	TrashIds    []string  `pg:"-"`
}

var _ pg.BeforeInsertHook = (*CreateEvent)(nil)

func (u *CreateEvent) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.Id = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	return ctx, nil
}

type Event struct {
	tableName    struct{} `pg:"events"json:"-"`
	Id           string   `pg:",pk"`
	Date         time.Time
	Description  string
	CreatedAt    time.Time `pg:"default:now()"`
	Trash        []*Trash  `pg:"many2many:events_trash"`
	UsersIds     []*EventUser
	SocietiesIds []*EventSociety
}

type EventPermission string

//    'creator',
//    'editor',
//    'viewer'

type EventSociety struct {
	tableName  struct{} `pg:"events_societies"json:"-"`
	SocietyId  string   `pg:",pk"`
	EventId    string   `pg:",pk"`
	Permission EventPermission
}

type EventUser struct {
	tableName  struct{} `pg:"events_users"json:"-"`
	UserId     string   `pg:",pk"`
	EventId    string   `pg:",pk"`
	Permission EventPermission
}

type EventTrash struct {
	tableName struct{} `pg:"events_trash"json:"-"`
	EventId   string   `pg:",pk"`
	TrashId   string   `pg:",pk"`
}

type EventPagingAnsw struct {
	Events []Event
	Paging Paging
}
