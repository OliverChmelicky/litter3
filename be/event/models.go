package event

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
	Publc       bool
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
	Permission       eventPermission
	AsSociety        bool
	SocietyId        string
}

type Event struct {
	tableName    struct{} `pg:"events"json:"-"`
	Id           string   `pg:",pk"`
	Date         time.Time
	Description  string
	Publc        bool      `pg:",use_zero"`
	CreatedAt    time.Time `pg:"default:now()"`
	TrashIds     []string  `pg:"-"`
	UsersIds     []string  `pg:"-"`
	SocietiesIds []string  `pg:"-"`
}

var _ pg.BeforeInsertHook = (*Event)(nil)

func (u *Event) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.Id = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	return ctx, nil
}

type eventPermission string

//    'creator',
//    'editor',
//    'viewer'

type EventSociety struct {
	tableName  struct{} `pg:"events_societies"json:"-"`
	SocietyId  string   `pg:",pk"`
	EventId    string   `pg:",pk"`
	Permission eventPermission
}

type EventUser struct {
	tableName  struct{} `pg:"events_users"json:"-"`
	UserId     string   `pg:",pk"`
	EventId    string   `pg:",pk"`
	Permission eventPermission
}

type EventTrash struct {
	tableName struct{} `pg:"events_trash"json:"-"`
	EventId   string   `pg:",pk"`
	TrashId   string   `pg:",pk"`
}
