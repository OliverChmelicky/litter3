package event

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/olo/litter3/user"
)

type eventAccess struct {
	db *pg.DB
}

func (s *eventAccess) CreateEvent(request *EventRequest) (*Event, error) {
	creatorUser := new(user.User)
	creatorUser.Id = request.UserId
	err := s.db.Select(creatorUser)
	if err != nil {
		return &Event{}, fmt.Errorf("Error select user: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &Event{}, fmt.Errorf("Error creating transaction: %w", err)
	}
	defer tx.Rollback()

	event := &Event{
		Publc: request.Publc,
		Date:  request.Date,
	}

	err = tx.Insert(event)
	if err != nil {
		tx.Rollback()
		return &Event{}, fmt.Errorf("Error creating event: %w", err)
	}

	err = tx.Select(event)
	fmt.Println(err)
	fmt.Println("pici je")
	fmt.Println(event)

	if request.AsSociety {
		creator := &EventSociety{
			Permission: attendaceLevel("creator"),
			SocietyId:  request.SocietyId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			tx.Rollback()
			return &Event{}, fmt.Errorf("Error inserting society creator: %w", err)
		}
		event.SocietiesIds = append(event.SocietiesIds, request.SocietyId)
	} else {
		creator := &EventUser{
			Permission: attendaceLevel("creator"),
			UserId:     request.UserId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			tx.Rollback()
			return &Event{}, fmt.Errorf("Error inserting user creator: %w", err)
		}
		event.UsersIds = append(event.UsersIds, request.UserId)
		fmt.Println("A tu preslo....")
	}

	event.TrashIds = request.Trash
	//event.TrashIds = append(event.TrashIds, trashId)
	fmt.Println("Som tu?")
	fmt.Println(event.TrashIds)
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		tx.Rollback()
		return &Event{}, fmt.Errorf("Error assigning trash: %w", err)
	}

	return event, tx.Commit()
}

func (s *eventAccess) GetEvent(eventId string) (*Event, error) {
	event := new(Event)
	err := s.db.Model(event).
		Column("event.*").
		Relation("events_societies", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("event_id = ?", eventId), nil
		}).
		Relation("users_societies", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("event_id = ?", eventId), nil
		}).
		First()
	if err != nil {
		return &Event{}, err
	}

	return event, err
}

func (s *eventAccess) AttendEvent(request *EventAttendanceRequest) (*EventAttendanceRequest, error) {
	if request.AsSociety {
		attendee := &EventSociety{
			Permission: attendaceLevel("viewer"),
			SocietyId:  request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &EventAttendanceRequest{}, fmt.Errorf("Error inserting society attendee: %w", err)
		}
	} else {
		attendee := &EventUser{
			Permission: attendaceLevel("viewer"),
			UserId:     request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &EventAttendanceRequest{}, fmt.Errorf("Error inserting user attendee: %w", err)
		}
	}

	return request, nil
}

func (s *eventAccess) CannotAttendEvent(request *EventAttendanceRequest) (*EventAttendanceRequest, error) {
	if request.AsSociety {
		attendee := &EventSociety{
			SocietyId: request.PickerId,
			EventId:   request.EventId,
		}
		err := s.db.Delete(attendee)
		if err != nil {
			return &EventAttendanceRequest{}, fmt.Errorf("Error deleting society attendee: %w", err)
		}
	} else {
		attendee := &EventUser{
			UserId:  request.PickerId,
			EventId: request.EventId,
		}
		err := s.db.Delete(attendee)
		if err != nil {
			return &EventAttendanceRequest{}, fmt.Errorf("Error deleting user attendee: %w", err)
		}
	}

	return request, nil
}

//func (s *eventAccess) EditEvent() (Event, error) {
//
//}
//
//func (s *eventAccess) EditEventRights() (Event, error) {
//
//}
//
//func (s *eventAccess) DeleteEvent() error {
//
//}
//
//func (s *eventAccess) GetSocietyEvents() {
//
//}
//
//func (s *eventAccess) GetUserEvents() {
//
//}
//
//func (s *eventAccess) CreateCollection() {
//
//}

func (s *eventAccess) AssignTrashToEvent(tx *pg.Tx, event *Event) error {
	relation := new(EventTrash)
	relation.EventId = event.Id
	fmt.Println("EventId")
	fmt.Println(event.Id)
	_, err := tx.Model(relation).Where("event_id = ?", event.Id).Delete()
	if err != nil {
		return fmt.Errorf("Error delete previous trash %w", err)
	}

	for _, trashId := range event.TrashIds {
		relation.TrashId = trashId
		err = tx.Insert(relation)
		if err != nil {
			return fmt.Errorf("Error insert trasId %s: %w", trashId, err)
		}
	}

	return nil
}
