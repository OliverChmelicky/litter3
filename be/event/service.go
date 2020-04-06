package event

import "github.com/go-pg/pg/v9"

type EventSerivice struct {
	*eventAccess
}

func CreateService(db *pg.DB) *EventService {
	access := &eventAccess{db: db}
	return &EventService{access}
}
