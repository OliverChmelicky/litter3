package event

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	"github.com/sirupsen/logrus"
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
		Publc:       request.Publc,
		Date:        request.Date,
		Description: request.Description,
	}

	err = tx.Insert(event)
	if err != nil {
		return &Event{}, fmt.Errorf("Error inserting event: %w", err)
	}

	if request.AsSociety {
		creator := &EventSociety{
			Permission: eventPermission("creator"),
			SocietyId:  request.SocietyId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
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
			return &Event{}, fmt.Errorf("Error inserting user creator: %w", err)
		}
		event.UsersIds = append(event.UsersIds, request.UserId)
	}

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
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

func (s *eventAccess) AttendEvent(request *EventPickerRequest) (*EventPickerRequest, error) {
	if request.AsSociety {
		attendee := &EventSociety{
			Permission: eventPermission("viewer"),
			SocietyId:  request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error inserting society attendee: %w", err)
		}
	} else {
		attendee := &EventUser{
			Permission: eventPermission("viewer"),
			UserId:     request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error inserting user attendee: %w", err)
		}
	}

	return request, nil
}

func (s *eventAccess) CannotAttendEvent(request *EventPickerRequest) (*EventPickerRequest, error) {
	if request.AsSociety {
		attendee := &EventSociety{
			SocietyId: request.PickerId,
			EventId:   request.EventId,
		}
		err := s.db.Select(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error finding society for event: %w ", err)
		}
		if attendee.Permission == eventPermission("creator") {
			return &EventPickerRequest{}, fmt.Errorf("You are an organizer ")
		}
		err = s.db.Delete(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error deleting society attendee: %w", err)
		}
	} else {
		attendee := &EventUser{
			UserId:  request.PickerId,
			EventId: request.EventId,
		}
		err := s.db.Select(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error finding user for event: %w ", err)
		}
		if attendee.Permission == eventPermission("creator") {
			return &EventPickerRequest{}, fmt.Errorf("You are an organizer ")
		}
		err = s.db.Delete(attendee)
		if err != nil {
			return &EventPickerRequest{}, fmt.Errorf("Error deleting user attendee: %w", err)
		}
	}

	return request, nil
}

func (s *eventAccess) UpdateEvent(request *EventRequest, userId string) (*Event, error) {
	if request.AsSociety {
		permission, err := s.HasSocietyEventPermission(request.SocietyId, request.Id, &[]eventPermission{"editor", "creator"})
		if err != nil {
			return &Event{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &Event{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	} else {
		permission, err := s.HasUserEventPermission(userId, request.Id, &[]eventPermission{"editor", "creator"})
		if err != nil {
			return &Event{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &Event{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &Event{}, fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	var event = new(Event)
	event.Id = request.Id
	err = s.db.Select(event)
	if err != nil {
		return &Event{}, fmt.Errorf("Error selecting event for update: %w ", err)
	}

	err = tx.Update(event)
	if err != nil {
		return &Event{}, fmt.Errorf("Error updating event: %w ", err)
	}

	err = tx.Select(event)

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
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

func (s *eventAccess) EditEventRights(request *EventPermissionRequest, userWhoDoesOperation string) (*EventPermissionRequest, error) {
	var isCreator bool
	if request.AsSociety {
		permission, err := s.HasSocietyEventPermission(request.SocietyId, request.EventId, &[]eventPermission{"editor", "creator"})
		if err != nil {
			return &EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &EventPermissionRequest{}, fmt.Errorf("You have no permisssion to edit event ")
		}
		isCreator, err = s.HasSocietyEventPermission(request.SocietyId, request.EventId, &[]eventPermission{"creator"})
		if err != nil {
			return &EventPermissionRequest{}, fmt.Errorf("Error check is creator: %w ", err)
		}
	} else {
		permission, err := s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]eventPermission{"editor", "creator"})
		if err != nil {
			return &EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &EventPermissionRequest{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	}

	if request.ChangingRightsTo == userWhoDoesOperation || request.ChangingRightsTo == request.SocietyId {
		return &EventPermissionRequest{}, fmt.Errorf("You cannot alter your permission")
	}

	if isCreator && request.Permission == eventPermission("creator") {
		return &EventPermissionRequest{}, fmt.Errorf("You there can be only one creator of event ")
	}

	if request.AsSociety {
		updating := new(EventSociety)
		updating.Permission = request.Permission
		updating.SocietyId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(&updating)
		if err != nil {
			return &EventPermissionRequest{}, fmt.Errorf("Couldn`t update society ")
		}
	} else {
		updating := new(EventUser)
		updating.Permission = request.Permission
		updating.UserId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(&updating)
		if err != nil {
			return &EventPermissionRequest{}, fmt.Errorf("Couldn`t update user ")
		}
	}

	return request, nil
}

func (s *eventAccess) DeleteEvent(request *EventPickerRequest, userWhoDoesOperation string) error {
	if request.AsSociety {
		isCreator, err := s.HasSocietyEventPermission(request.PickerId, request.EventId, &[]eventPermission{"creator"})
		if err != nil {
			return fmt.Errorf("Error check is creator: %w ", err)
		}

		if !isCreator {
			return fmt.Errorf("You have no permission to delete event ")
		}

	} else {
		isCreator, err := s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]eventPermission{"creator"})
		if err != nil {
			return fmt.Errorf("Error check is creator: %w ", err)
		}

		if !isCreator {
			return fmt.Errorf("You have no permission to delete event ")
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	userEvent := new(EventUser)
	_, err = tx.Model(userEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete users from event %w ", err)
	}

	societyEvent := new(EventSociety)
	_, err = tx.Model(societyEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete societies from event %w ", err)
	}

	trashEvent := new(EventTrash)
	_, err = tx.Model(trashEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete trash from event %w ", err)
	}

	event := &Event{Id: request.EventId}
	err = tx.Delete(event)
	if err != nil {
		return fmt.Errorf("Error delete event %w ", err)
	}

	return tx.Commit()
}

func (s *eventAccess) GetSocietyEvents(societyId string) ([]Event, error) {
	var allActivities []EventSociety
	err := s.db.Model(&allActivities).Where("society_id = ?", societyId).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society participation: %w ", err)
	}

	var eventsArr []string
	for _, activity := range allActivities {
		eventsArr = append(eventsArr, activity.EventId)
	}

	var events []Event
	err = s.db.Model(&events).Where("id IN (?)", pg.In(eventsArr)).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}

	return events, nil
}

func (s *eventAccess) GetUserEvents(userId string) ([]Event, error) {
	var allActivities []EventUser
	err := s.db.Model(&allActivities).Where("user_id = ?", userId).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}

	var eventsArr []string
	for _, activity := range allActivities {
		eventsArr = append(eventsArr, activity.EventId)
	}

	var events []Event
	err = s.db.Model(&events).Where("id IN (?)", pg.In(eventsArr)).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}
	return events, nil
}

func (s *eventAccess) CreateCollections(collectionRequests *trash.CreateCollectionOrganizedRequest) ([]trash.Collection, []error) {
	var errs []error
	var collections []trash.Collection

	collection := &trash.Collection{}
	for _, request := range collectionRequests.Collections {
		collection.EventId = collectionRequests.EventId
		collection.TrashId = request.TrashId
		collection.CleanedTrash = request.CleanedTrash
		collection.Weight = request.Weight
		err := s.db.Insert(collection)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"trashId": request.TrashId,
				"eventId": collectionRequests.EventId,
				"error":   err.Error(),
			}).Error("Error inserting new collection")
			errs = append(errs, err)
			continue
		}
		collections = append(collections, *collection)
		if request.CleanedTrash {
			updating := new(trash.Trash)
			updating.Id = request.TrashId
			_, err = s.db.Model(updating).Column("cleaned").Where("id = ?", request.TrashId).Update()
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		*collection = trash.Collection{}
	}

	return collections, errs
}

func (s *eventAccess) HasUserEventPermission(userId, eventId string, editPermission *[]eventPermission) (bool, error) {
	relation := new(EventUser)
	relation.UserId = userId
	relation.EventId = eventId
	err := s.db.Select(relation)
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
