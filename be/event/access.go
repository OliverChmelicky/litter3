package event

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/olo/litter3/models"
	"github.com/sirupsen/logrus"
)

type eventAccess struct {
	db *pg.DB
}

func (s *eventAccess) CreateEvent(request *models.EventRequest) (*models.Event, error) {
	creatorUser := new(models.User)
	creatorUser.Id = request.UserId
	err := s.db.Select(creatorUser)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error select user: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error creating transaction: %w", err)
	}
	defer tx.Rollback()

	event := &models.Event{
		Publc:       request.Publc,
		Date:        request.Date,
		Description: request.Description,
	}

	err = tx.Insert(event)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error inserting event: %w", err)
	}

	if request.AsSociety {
		creator := &models.EventSociety{
			Permission: models.EventPermission("creator"),
			SocietyId:  request.SocietyId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			return &models.Event{}, fmt.Errorf("Error inserting society creator: %w", err)
		}
		event.SocietiesIds = append(event.SocietiesIds, request.SocietyId)
	} else {
		creator := &models.EventUser{
			Permission: models.EventPermission("creator"),
			UserId:     request.UserId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			return &models.Event{}, fmt.Errorf("Error inserting user creator: %w", err)
		}
		event.UsersIds = append(event.UsersIds, request.UserId)
	}

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error assigning trash: %w", err)
	}

	return event, tx.Commit()
}

func (s *eventAccess) GetEvent(eventId string) (*models.Event, error) {
	event := new(models.Event)
	event.Id = eventId
	err := s.db.Select(event)
	if err != nil {
		return &models.Event{}, err
	}
	var trash []models.EventTrash
	err = s.db.Model(&trash).Where("event_id = ?", eventId).Select(&trash)
	if err != nil {
		return &models.Event{}, err
	}
	var users []models.EventUser
	err = s.db.Model(&users).Where("event_id = ?", eventId).Select(&users)
	if err != nil {
		return &models.Event{}, err
	}
	var societies []models.EventSociety
	err = s.db.Model(&societies).Where("event_id = ?", eventId).Select(&societies)
	if err != nil {
		return &models.Event{}, err
	}

	mapEvent(event, trash, users, societies)

	return event, err
}

func mapEvent(event *models.Event, trash []models.EventTrash, users []models.EventUser, societies []models.EventSociety) {
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

func (s *eventAccess) AttendEvent(request *models.EventPickerRequest) (*models.EventPickerRequest, error) {
	if request.AsSociety {
		attendee := &models.EventSociety{
			Permission: models.EventPermission("viewer"),
			SocietyId:  request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error inserting society attendee: %w", err)
		}
	} else {
		attendee := &models.EventUser{
			Permission: models.EventPermission("viewer"),
			UserId:     request.PickerId,
			EventId:    request.EventId,
		}
		err := s.db.Insert(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error inserting user attendee: %w", err)
		}
	}

	return request, nil
}

func (s *eventAccess) CannotAttendEvent(request *models.EventPickerRequest) (*models.EventPickerRequest, error) {
	if request.AsSociety {
		attendee := &models.EventSociety{
			SocietyId: request.PickerId,
			EventId:   request.EventId,
		}
		err := s.db.Select(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error finding society for event: %w ", err)
		}
		if attendee.Permission == "creator" {
			return &models.EventPickerRequest{}, fmt.Errorf("You are an organizer ")
		}
		err = s.db.Delete(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error deleting society attendee: %w", err)
		}
	} else {
		attendee := &models.EventUser{
			UserId:  request.PickerId,
			EventId: request.EventId,
		}
		err := s.db.Select(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error finding user for event: %w ", err)
		}
		if attendee.Permission == "creator" {
			return &models.EventPickerRequest{}, fmt.Errorf("You are an organizer ")
		}
		err = s.db.Delete(attendee)
		if err != nil {
			return &models.EventPickerRequest{}, fmt.Errorf("Error deleting user attendee: %w", err)
		}
	}

	return request, nil
}

func (s *eventAccess) UpdateEvent(request *models.EventRequest, userId string) (*models.Event, error) {
	if request.AsSociety {
		permission, err := s.HasSocietyEventPermission(request.SocietyId, request.Id, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.Event{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.Event{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	} else {
		permission, err := s.HasUserEventPermission(userId, request.Id, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.Event{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.Event{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	var event = new(models.Event)
	event.Id = request.Id
	err = s.db.Select(event)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error selecting event for update: %w ", err)
	}

	err = tx.Update(event)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error updating event: %w ", err)
	}

	err = tx.Select(event)

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		return &models.Event{}, fmt.Errorf("Error assigning trash: %w ", err)
	}

	var users []models.EventUser
	err = s.db.Model(&users).Where("event_id = ?", event.Id).Select(&users)
	if err != nil {
		return &models.Event{}, err
	}
	var societies []models.EventSociety
	err = s.db.Model(&societies).Where("event_id = ?", event.Id).Select(&societies)
	if err != nil {
		return &models.Event{}, err
	}

	mapEvent(event, nil, users, societies)

	return event, tx.Commit()
}

func (s *eventAccess) EditEventRights(request *models.EventPermissionRequest, userWhoDoesOperation string) (*models.EventPermissionRequest, error) {
	var isCreator bool
	if request.AsSociety {
		permission, err := s.HasSocietyEventPermission(request.SocietyId, request.EventId, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.EventPermissionRequest{}, fmt.Errorf("You have no permisssion to edit event ")
		}
		isCreator, err = s.HasSocietyEventPermission(request.SocietyId, request.EventId, &[]models.EventPermission{"creator"})
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check is creator: %w ", err)
		}
	} else {
		permission, err := s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.EventPermissionRequest{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	}

	if request.ChangingRightsTo == userWhoDoesOperation || request.ChangingRightsTo == request.SocietyId {
		return &models.EventPermissionRequest{}, fmt.Errorf("You cannot alter your permission ")
	}

	if isCreator && request.Permission == models.EventPermission("creator") {
		return &models.EventPermissionRequest{}, fmt.Errorf("You there can be only one creator of event ")
	}

	if request.AsSociety {
		updating := new(models.EventSociety)
		updating.Permission = request.Permission
		updating.SocietyId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(&updating)
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Couldn`t update society ")
		}
	} else {
		updating := new(models.EventUser)
		updating.Permission = request.Permission
		updating.UserId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(&updating)
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Couldn`t update user ")
		}
	}

	return request, nil
}

func (s *eventAccess) DeleteEvent(request *models.EventPickerRequest, userWhoDoesOperation string) error {
	if request.AsSociety {
		isCreator, err := s.HasSocietyEventPermission(request.PickerId, request.EventId, &[]models.EventPermission{"creator"})
		if err != nil {
			return fmt.Errorf("Error check is creator: %w ", err)
		}

		if !isCreator {
			return fmt.Errorf("You have no permission to delete event ")
		}

	} else {
		isCreator, err := s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]models.EventPermission{"creator"})
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

	userEvent := new(models.EventUser)
	_, err = tx.Model(userEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete users from event %w ", err)
	}

	societyEvent := new(models.EventSociety)
	_, err = tx.Model(societyEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete societies from event %w ", err)
	}

	trashEvent := new(models.EventTrash)
	_, err = tx.Model(trashEvent).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete trash from event %w ", err)
	}

	//maybe navráť stav ak nie je later collection
	collection := new(models.Collection)
	_, err = tx.Model(collection).Where("event_id = ?", request.EventId).Delete()
	if err != nil {
		return fmt.Errorf("Error delete collection from event %w ", err)
	}

	event := &models.Event{Id: request.EventId}
	err = tx.Delete(event)
	if err != nil {
		return fmt.Errorf("Error delete event %w ", err)
	}

	return tx.Commit()
}

func (s *eventAccess) GetSocietyEvents(societyId string) ([]models.Event, error) {
	var allActivities []models.EventSociety
	err := s.db.Model(&allActivities).Where("society_id = ?", societyId).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society participation: %w ", err)
	}

	var eventsArr []string
	for _, activity := range allActivities {
		eventsArr = append(eventsArr, activity.EventId)
	}

	var events []models.Event
	err = s.db.Model(&events).Where("id IN (?)", eventsArr).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}

	return events, nil
}

func (s *eventAccess) GetUserEvents(userId string) ([]models.Event, error) {
	var allActivities []models.EventUser
	err := s.db.Model(&allActivities).Where("user_id = ?", userId).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}

	var eventsArr []string
	for _, activity := range allActivities {
		eventsArr = append(eventsArr, activity.EventId)
	}

	var events []models.Event
	err = s.db.Model(&events).Where("id IN (?)", pg.In(eventsArr)).Select()
	if err != nil {
		return nil, fmt.Errorf("Error get society events: %w ", err)
	}
	return events, nil
}

func (s *eventAccess) CreateCollectionsOrganized(collectionRequests *models.CreateCollectionOrganizedRequest) ([]models.Collection, []error) {
	var errs []error
	var collections []models.Collection

	rights, err := s.CheckPickerRights(collectionRequests.OrganizerId, collectionRequests.EventId, collectionRequests.AsSociety)
	if err != nil {
		errs = append(errs, fmt.Errorf("Error verifying rights for creating organized collection: %w ", err))
		return nil, errs
	}
	if !rights {
		errs = append(errs, fmt.Errorf("No rights for creating organized collection "))
		return nil, errs
	}

	collection := &models.Collection{}
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
			updating := new(models.Trash)
			updating.Id = request.TrashId
			_, err = s.db.Model(updating).Column("cleaned").Where("id = ?", request.TrashId).Update()
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		*collection = models.Collection{}
	}

	return collections, errs
}

func (s *eventAccess) UpdateCollectionOrganized(request *models.UpdateCollectionOrganizedRequest) (*models.Collection, error) {
	rights, err := s.CheckPickerRights(request.OrganizerId, request.EventId, request.AsSociety)
	if err != nil {
		return nil, err
	}
	if !rights {
		return nil, err
	}

	oldCollection := new(models.Collection)
	oldCollection.Id = request.Collection.Id
	if err := s.db.Select(oldCollection); err != nil {
		return &models.Collection{}, fmt.Errorf("Couldn`t get old collection: %w ", err)
	}
	if oldCollection.TrashId != request.Collection.TrashId {
		return &models.Collection{}, fmt.Errorf("Cannot change trashId to collection ")
	}
	if oldCollection.EventId != request.Collection.EventId {
		return &models.Collection{}, fmt.Errorf("Cannot change trashId to collection ")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Couldn`t create transaction: %w ", err)
	}
	defer tx.Rollback()

	updatedCollection := &models.Collection{
		Id:           request.Collection.Id,
		Weight:       request.Collection.Weight,
		CleanedTrash: request.Collection.CleanedTrash,
		TrashId:      oldCollection.TrashId,
		EventId:      oldCollection.EventId,
		CreatedAt:    oldCollection.CreatedAt,
	}
	err = tx.Update(updatedCollection)
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error update collection: %w ", err)
	}

	if oldCollection.CleanedTrash != request.Collection.CleanedTrash {
		updateTrash := new(models.Trash)
		updateTrash.Id = request.Collection.TrashId
		updateTrash.Cleaned = request.Collection.CleanedTrash
		_, err = tx.Model(updateTrash).Column("cleaned").Where("id = ?", request.Collection.TrashId).Update()
		if err != nil {
			return &models.Collection{}, err
		}
	}

	return updatedCollection, tx.Commit()
}

func (s *eventAccess) DeleteCollectionOrganized(organizerId, collectionId, eventId string, asSociety bool) error {
	rights, err := s.CheckPickerRights(organizerId, eventId, asSociety)
	if err != nil {
		return err
	}
	if !rights {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("Couldn`t create transaction: %w ", err)
	}
	defer tx.Rollback()

	oldCollection := new(models.Collection)
	oldCollection.Id = collectionId
	if err := tx.Select(oldCollection); err != nil {
		return fmt.Errorf("Couldn`t get old collection: %w ", err)
	}
	if oldCollection.EventId != eventId {
		return fmt.Errorf("Collection does not belong to event: %w ", err)
	}

	if oldCollection.CleanedTrash {
		updateTrash := new(models.Trash)
		updateTrash.Id = oldCollection.TrashId
		updateTrash.Cleaned = false
		_, err = tx.Model(updateTrash).Column("cleaned").Where("id = ?", oldCollection.TrashId).Update()
		if err != nil {
			return fmt.Errorf("Error updating trash before deleting collection: %w ", err)
		}
	}

	err = s.db.Delete(oldCollection)
	if err != nil {
		return fmt.Errorf("Error deleting collection: %w ", err)
	}

	return tx.Commit()
}

func (s *eventAccess) HasUserEventPermission(userId, eventId string, editPermission *[]models.EventPermission) (bool, error) {
	relation := new(models.EventUser)
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

func (s *eventAccess) HasSocietyEventPermission(societyId, eventId string, editPermission *[]models.EventPermission) (bool, error) {
	relation := new(models.EventSociety)
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

func (s *eventAccess) AssignTrashToEvent(tx *pg.Tx, event *models.Event) error {
	relation := new(models.EventTrash)
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

func (s *eventAccess) CheckPickerRights(organizerId string, eventId string, asSociety bool) (bool, error) {
	if asSociety {
		return s.HasSocietyEventPermission(organizerId, eventId, &[]models.EventPermission{"creator", "editor"})

	} else {
		return s.HasUserEventPermission(organizerId, eventId, &[]models.EventPermission{"creator", "editor"})
	}
}
