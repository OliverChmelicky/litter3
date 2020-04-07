package event

import (
	"fmt"
	"github.com/go-pg/pg/v9"
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

	if request.AsSociety {
		creator := &EventSociety{
			Permission: eventPermission("creator"),
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
			Permission: eventPermission("creator"),
			UserId:     request.UserId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			tx.Rollback()
			return &Event{}, fmt.Errorf("Error inserting user creator: %w", err)
		}
		event.UsersIds = append(event.UsersIds, request.UserId)
	}

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		tx.Rollback()
		return &Event{}, fmt.Errorf("Error assigning trash: %w", err)
	}

	return event, tx.Commit()
}

func (s *eventAccess) GetEvent(eventId string) (*Event, error) {
	event := new(Event)
	event.Id = eventId
	err := s.db.Select(event)
	if err != nil {
		return &Event{}, err
	}
	var trash []EventTrash
	err = s.db.Model(&trash).Where("event_id = ?", eventId).Select(&trash)
	if err != nil {
		return &Event{}, err
	}
	var users []EventUser
	err = s.db.Model(&users).Where("event_id = ?", eventId).Select(&users)
	if err != nil {
		return &Event{}, err
	}
	var societies []EventSociety
	err = s.db.Model(&societies).Where("event_id = ?", eventId).Select(&societies)
	if err != nil {
		return &Event{}, err
	}

	mapEvent(event, trash, users, societies)

	return event, err
}

func mapEvent(event *Event, trash []EventTrash, users []EventUser, societies []EventSociety) {
	if trash != nil {
		for _, v := range trash {
			event.TrashIds = append(event.TrashIds, v.TrashId)
		}
	}
	if users != nil {
		for _, v := range users {
			event.UsersIds = append(event.UsersIds, v.UserId)
		}
	}
	if societies != nil {
		for _, v := range societies {
			event.SocietiesIds = append(event.SocietiesIds, v.SocietyId)
		}
	}
}

func (s *eventAccess) AttendEvent(request *EventAttendanceRequest) (*EventAttendanceRequest, error) {
	if request.AsSociety {
		attendee := &EventSociety{
			Permission: eventPermission("viewer"),
			SocietyId:  request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &EventAttendanceRequest{}, fmt.Errorf("Error inserting society attendee: %w", err)
		}
	} else {
		attendee := &EventUser{
			Permission: eventPermission("viewer"),
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

func (s *eventAccess) EditEvent(request *EventRequest, userId string) (*Event, error) {
	editorUser := new(user.User)
	editorUser.Id = userId
	err := s.db.Select(editorUser)
	if err != nil {
		return &Event{}, fmt.Errorf("Error select user: %w ", err)
	}

	permission, err := s.HasUserEventPermission(userId, request.Id, &[]eventPermission{"editor", "creator"})
	if err != nil {
		return &Event{}, fmt.Errorf("Error check permission: %w ", err)
	}

	if !permission {
		return &Event{}, fmt.Errorf("You have no permisssion to edit event ")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &Event{}, fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	event := &Event{
		Id:    request.Id,
		Publc: request.Publc,
		Date:  request.Date,
	}

	err = tx.Update(event)
	if err != nil {
		tx.Rollback()
		return &Event{}, fmt.Errorf("Error updating event: %w ", err)
	}

	err = tx.Select(event)

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		tx.Rollback()
		return &Event{}, fmt.Errorf("Error assigning trash: %w ", err)
	}

	var users []EventUser
	err = s.db.Model(&users).Where("event_id = ?", event.Id).Select(&users)
	if err != nil {
		return &Event{}, err
	}
	var societies []EventSociety
	err = s.db.Model(&societies).Where("event_id = ?", event.Id).Select(&societies)
	if err != nil {
		return &Event{}, err
	}

	mapEvent(event, nil, users, societies)

	return event, tx.Commit()
}

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

func (s *eventAccess) HasUserEventPermission(userId, eventId string, editPermission *[]eventPermission) (bool, error) {
	relation := new(EventUser)
	relation.UserId = userId
	relation.EventId = eventId
	err := s.db.Model(relation).Select(relation)
	if err != nil {
		return false, fmt.Errorf("Error get user-event relation: %w", err)
	}
	for _, p := range *editPermission {
		if p == relation.Permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *eventAccess) HasSocietyEventPermission(societyId, eventId string, editPermission *[]eventPermission) (bool, error) {
	relation := new(EventSociety)
	relation.SocietyId = societyId
	relation.EventId = eventId
	err := s.db.Model(relation).Select(relation)
	if err != nil {
		return false, fmt.Errorf("Error get society-event relation: %w", err)
	}
	for _, p := range *editPermission {
		if p == relation.Permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *eventAccess) AssignTrashToEvent(tx *pg.Tx, event *Event) error {
	relation := new(EventTrash)
	relation.EventId = event.Id

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
