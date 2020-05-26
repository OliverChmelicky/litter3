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

func (s *eventAccess) CreateEvent(request *models.EventRequest) (*models.CreateEvent, error) {
	creatorUser := new(models.User)
	creatorUser.Id = request.UserId
	err := s.db.Select(creatorUser)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error select user: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error creating transaction: %w", err)
	}
	defer tx.Rollback()

	event := &models.CreateEvent{
		Date:        request.Date,
		Description: request.Description,
	}

	err = tx.Insert(event)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error inserting event: %w", err)
	}

	if request.AsSociety {
		creator := &models.EventSociety{
			Permission: models.EventPermission("creator"),
			SocietyId:  request.SocietyId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			return &models.CreateEvent{}, fmt.Errorf("Error inserting society creator: %w", err)
		}
	} else {
		creator := &models.EventUser{
			Permission: models.EventPermission("creator"),
			UserId:     request.UserId,
			EventId:    event.Id,
		}
		err = tx.Insert(creator)
		if err != nil {
			return &models.CreateEvent{}, fmt.Errorf("Error inserting user creator: %w", err)
		}
	}

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error assigning trash: %w", err)
	}

	return event, tx.Commit()
}

func (s *eventAccess) GetEvent(eventId string) (*models.Event, error) {
	event := new(models.Event)
	event.Id = eventId

	err := s.db.Model(event).Where("id = ?", eventId).Column("event.*").
		Relation("Trash").Relation("SocietiesIds").Relation("UsersIds").First()
	if err != nil {
		return &models.Event{}, err
	}

	return event, nil
}

func (s *eventAccess) GetEventWithCollection(eventId string) (*models.EventWithCollections, error) {
	event := new(models.Event)

	err := s.db.Model(event).Where("id = ?", eventId).Column("event.*").
		Relation("Trash").Relation("SocietiesIds").Relation("UsersIds").First()
	if err != nil {
		return &models.EventWithCollections{}, err
	}

	var collections []models.Collection
	err = s.db.Model(&collections).Where("event_id = ?", event.Id).Select()
	if err != nil {
		return &models.EventWithCollections{}, err
	}

	eventWithCollections := new(models.EventWithCollections)
	eventWithCollections.Id = event.Id
	eventWithCollections.CreatedAt = event.CreatedAt
	eventWithCollections.Description = event.Description
	eventWithCollections.Date = event.Date
	eventWithCollections.Trash = event.Trash
	eventWithCollections.SocietiesIds = event.SocietiesIds
	eventWithCollections.UsersIds = event.UsersIds
	eventWithCollections.Collections = collections

	for i, c := range eventWithCollections.Collections {
		var images []models.CollectionImage

		err = s.db.Model(&images).Where("collection_id = ?", c.Id).Select()
		if err != nil {
			eventWithCollections.Collections[i].Images = []models.CollectionImage{}
		}
		eventWithCollections.Collections[i].Images = images
	}

	return eventWithCollections, nil
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

func (s *eventAccess) UpdateEvent(request *models.EventRequest, userId string) (*models.CreateEvent, error) {
	if request.AsSociety {
		permission, err := s.HasSocietyEventPermission(request.SocietyId, request.Id, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.CreateEvent{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.CreateEvent{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	} else {
		permission, err := s.HasUserEventPermission(userId, request.Id, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.CreateEvent{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.CreateEvent{}, fmt.Errorf("You have no permisssion to edit event ")
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	var event = new(models.CreateEvent)
	event.Id = request.Id
	err = s.db.Select(event)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error selecting event for update: %w ", err)
	}

	event.Date = request.Date
	event.Description = request.Description
	err = tx.Update(event)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error updating event: %w ", err)
	}

	event.TrashIds = request.Trash
	err = s.AssignTrashToEvent(tx, event)
	if err != nil {
		return &models.CreateEvent{}, fmt.Errorf("Error assigning trash: %w ", err)
	}

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
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if isCreator && request.ChangingRightsTo == request.SocietyId {
			return &models.EventPermissionRequest{}, fmt.Errorf("You are creator and cannot change rights to yourself ")
		}
	} else {
		permission, err := s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]models.EventPermission{"editor", "creator"})
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}
		if !permission {
			return &models.EventPermissionRequest{}, fmt.Errorf("You have no permisssion to edit event ")
		}
		isCreator, err = s.HasUserEventPermission(userWhoDoesOperation, request.EventId, &[]models.EventPermission{"creator"})
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Error check permission: %w ", err)
		}

		if isCreator && request.ChangingRightsTo == userWhoDoesOperation {
			return &models.EventPermissionRequest{}, fmt.Errorf("You are creator and cannot change rights to yourself ")
		}
	}

	if request.ChangingToSociety {
		updating := new(models.EventSociety)
		updating.Permission = request.Permission
		updating.SocietyId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(updating)
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Couldn`t update society permission: %w ", err)
		}
	} else {
		updating := new(models.EventUser)
		updating.Permission = request.Permission
		updating.UserId = request.ChangingRightsTo
		updating.EventId = request.EventId
		err := s.db.Update(updating)
		if err != nil {
			return &models.EventPermissionRequest{}, fmt.Errorf("Couldn`t update user permision %w ", err)
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

	event := &models.Event{Id: request.EventId}
	err := s.db.Delete(event)
	if err != nil {
		return fmt.Errorf("Error delete event %w ", err)
	}

	//tx, err := s.db.Begin()
	//if err != nil {
	//	return fmt.Errorf("Error creating transaction: %w ", err)
	//}
	//defer tx.Rollback()
	//
	//userEvent := new(models.EventUser)
	//_, err = tx.Model(userEvent).Where("event_id = ?", request.EventId).Delete()
	//if err != nil {
	//	return fmt.Errorf("Error delete users from event %w ", err)
	//}
	//
	//societyEvent := new(models.EventSociety)
	//_, err = tx.Model(societyEvent).Where("event_id = ?", request.EventId).Delete()
	//if err != nil {
	//	return fmt.Errorf("Error delete societies from event %w ", err)
	//}
	//
	//trashEvent := new(models.EventTrash)
	//_, err = tx.Model(trashEvent).Where("event_id = ?", request.EventId).Delete()
	//if err != nil {
	//	return fmt.Errorf("Error delete trash from event %w ", err)
	//}
	//
	////maybe navráť stav ak nie je later collection
	//collection := new(models.Collection)
	//_, err = tx.Model(collection).Where("event_id = ?", request.EventId).Delete()
	//if err != nil {
	//	return fmt.Errorf("Error delete collection from event %w ", err)
	//}
	//
	//event := &models.Event{Id: request.EventId}
	//err = tx.Delete(event)
	//if err != nil {
	//	return fmt.Errorf("Error delete event %w ", err)
	//}

	return nil
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
	if len(eventsArr) > 0 {
		err = s.db.Model(&events).Where("id IN (?)", pg.In(eventsArr)).Select()
		if err != nil {
			return nil, fmt.Errorf("Error get society events: %w ", err)
		}
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

	//newCollections, updateCollections, err = s.filterUpdateColectionAndCreateNew(collectionRequests)
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
			_, err := s.db.Model(updating).Set("cleaned = ?", true).Where("id = ?", request.TrashId).Update()
			if err != nil {
				logrus.Error("Error update trash to cleaned: ", err)
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
	if err := s.db.Model(oldCollection).Where("id = ?", oldCollection.Id).First(); err != nil {
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
	if err := tx.Model(oldCollection).Where("id = ?", collectionId).First(); err != nil {
		return fmt.Errorf("Couldn`t get old collection: %w ", err)
	}
	if oldCollection.EventId != eventId {
		return fmt.Errorf("Collection does not belong to event: %w ", err)
	}

	var laterCollections []models.Collection
	err = tx.Model(&laterCollections).Where("trash_id = ? and created_at > ?", oldCollection.TrashId, oldCollection.CreatedAt).Select()
	if err != nil {
		return fmt.Errorf("Error checking if no later collections: %w ", err)
	}
	if len(laterCollections) == 0 {
		if oldCollection.CleanedTrash {
			trash := new(models.Trash)
			trash.Cleaned = false
			_, err = tx.Model(trash).Column("cleaned").Where("id = ?", oldCollection.TrashId).Update()
			if err != nil {
				return fmt.Errorf("Error updating trash before deleting collection: %w ", err)
			}
		}
	}

	err = tx.Delete(oldCollection)
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
	err := s.db.Select(relation)
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

func (s *eventAccess) AssignTrashToEvent(tx *pg.Tx, event *models.CreateEvent) error {
	relation := new(models.EventTrash)
	_, err := tx.Model(relation).Where("event_id = ?", event.Id).Delete()
	if err != nil {
		return fmt.Errorf("Error delete previous trash %w", err)
	}

	relation.EventId = event.Id
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

func (s *eventAccess) GetEventsWithPaging(from int, to int) ([]models.Event, int, error) {
	limit := to - from
	events := []models.Event{}
	err := s.db.Model(&events).Order("created_at DESC").Select()
	if err != nil {
		return []models.Event{}, 0, err
	}

	if len(events) < from {
		return []models.Event{}, 0, fmt.Errorf("No records starting from FROM ")
	}
	if len(events[from:]) < limit {
		to = from + len(events[from:])
	}

	var ids []string
	for _, event := range events {
		ids = append(ids, event.Id)
	}

	for i, event := range events {
		wholeEvent, err := s.GetEvent(event.Id)
		if err != nil {
			return []models.Event{}, 0, fmt.Errorf("Error filling event details: %w ", err)
		}
		events[i] = *wholeEvent
	}

	return events[from:to], len(events), nil
}

//func (s *eventAccess) filterUpdateColectionAndCreateNew(requests *models.CreateCollectionOrganizedRequest) ([]models.CreateCollectionRandomRequest, []models.CreateCollectionRandomRequest, error) {
//	newRequests := []models.CreateCollectionOrganizedRequest{}
//	updateRequests := []models.CreateCollectionOrganizedRequest{}
//	for _, request := range requests.Collections {
//		err := s.db.Model(&models.Collection{}).Where("event_id = ? and trash_id = ?", requests.EventId, request.TrashId).First()
//		if err == nil {
//
//		}
//	}
//}
