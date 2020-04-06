package event

import "time"

type EventRequest struct {
}

type Event struct {
	tableName  struct{} `pg:"events"json:"-"`
	id         string
	date       time.Time
	publc      bool
	user_id    string
	society_id string
	created_at time.Time
}

type attendaceLevel string

type EventSociety struct {
	tableName  struct{} `pg:"events_societies"json:"-"`
	society_id string
	event_id   string
	permission attendaceLevel
	created_at time.Time
}

type EventUser struct {
	tableName  struct{} `pg:"events_users"json:"-"`
	user_id    string
	event_id   string
	permission attendaceLevel
	created_at time.Time
}

type EventTrash struct {
	tableName struct{} `pg:"events_trash"json:"-"`
	event_id  string
	trash_id  string
}
