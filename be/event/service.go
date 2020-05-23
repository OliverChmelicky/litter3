package event

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type EventService struct {
	*eventAccess
	*user.UserAccess
	*trash.TrashAccess
}

func CreateService(db *pg.DB) *EventService {
	access := &eventAccess{db: db}
	userAccess := &user.UserAccess{Db: db}
	trashAccess := &trash.TrashAccess{Db: db}
	return &EventService{access, userAccess, trashAccess}
}

func (s *EventService) CreateEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.SocietyId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.UserId = userId
	newTrash, err := s.eventAccess.CreateEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateEvent, err))
	}

	return c.JSON(http.StatusCreated, newTrash)
}
func (s *EventService) GetEvent(c echo.Context) error {
	eventId := c.Param("eventId")

	event, err := s.eventAccess.GetEventWithCollection(eventId)
	if err != nil {
		logrus.Error("Co je tu zza err: ", err)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetEvent, err))
	}

	fmt.Println("Vraciam event: ",event)
	fmt.Println("Collections are: ", event.Collections)
	//fmt.Println("images in col: ", event.Collections[0].Images)

	return c.JSON(http.StatusOK, event)
}

func (s *EventService) GetEventsWithPaging(c echo.Context) error {
	//can call also unregistered user
	f := c.QueryParam("from")
	t := c.QueryParam("to")

	from, err := strconv.Atoi(f)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	to, err := strconv.Atoi(t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	if to-from < 0 {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, fmt.Errorf("To is smaller than from")))
	}

	events, allSocieties, err := s.eventAccess.GetEventsWithPaging(from, to)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetSociety, err))
	}

	return c.JSON(http.StatusOK, models.EventPagingAnsw{
		Events: events,
		Paging: models.Paging{From: from, To: to, TotalCount: allSocieties},
	})
}

func (s *EventService) AttendEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventPickerRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.PickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrAttendEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	} else {
		request.PickerId = userId
	}

	attendeeRelation, err := s.eventAccess.AttendEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrAttendEvent, err))
	}

	return c.JSON(http.StatusCreated, attendeeRelation)
}

func (s *EventService) CannotAttendEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventPickerRequest)
	eventId := c.QueryParam("event")
	pickerId := c.QueryParam("picker")
	asSociety, err := strconv.ParseBool(c.QueryParam("asSociety"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	request.EventId = eventId
	request.PickerId = pickerId
	request.AsSociety = asSociety

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.PickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCannotAttendEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	} else {
		request.PickerId = userId
	}

	_, err = s.eventAccess.CannotAttendEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCannotAttendEvent, err))
	}

	return c.NoContent(http.StatusOK)
}

func (s *EventService) UpdateEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.SocietyId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.UserId = userId
	updatedEvent, err := s.eventAccess.UpdateEvent(request, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateEvent, err))
	}

	return c.JSON(http.StatusOK, updatedEvent)
}

func (s *EventService) EditEventRights(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventPermissionRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.SocietyId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrEditEventRights, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	newTrash, err := s.eventAccess.EditEventRights(request, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrEditEventRights, err))
	}

	return c.JSON(http.StatusOK, newTrash)
}

func (s *EventService) DeleteEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.EventPickerRequest)
	eventId := c.QueryParam("event")
	pickerId := c.QueryParam("picker")
	asSociety, err := strconv.ParseBool(c.QueryParam("asSociety"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	request.EventId = eventId
	request.PickerId = pickerId
	request.AsSociety = asSociety

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.PickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrDeleteEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	err = s.eventAccess.DeleteEvent(request, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrUpdateEvent, err))
	}

	return c.JSON(http.StatusOK, "")
}

func (s *EventService) GetSocietyEvents(c echo.Context) error {
	societyId := c.QueryParam("societyId")

	events, err := s.eventAccess.GetSocietyEvents(societyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrGetSocietyEvent, err))
	}

	return c.JSON(http.StatusOK, events)
}

func (s *EventService) GetUserEvents(c echo.Context) error {
	searchedUserId := c.QueryParam("picker")

	events, err := s.eventAccess.GetUserEvents(searchedUserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrGetSocietyEvent, err))
	}

	return c.JSON(http.StatusOK, events)
}

func (s *EventService) CreateCollectionsOrganized(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.CreateCollectionOrganizedRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.OrganizerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateCollectionFromEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	} else {
		request.OrganizerId = userId
	}

	collections, errs := s.eventAccess.CreateCollectionsOrganized(request)
	if len(errs) != 0 {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionFromEvent, err))
	}

	return c.JSON(http.StatusCreated, collections)
}

func (s *EventService) UpdateCollectionOrganized(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.UpdateCollectionOrganizedRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.OrganizerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateCollectionFromEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	} else {
		request.OrganizerId = userId
	}

	collections, err := s.eventAccess.UpdateCollectionOrganized(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionFromEvent, err))
	}

	return c.JSON(http.StatusCreated, collections)
}

func (s *EventService) DeleteCollection(c echo.Context) error {
	userId := c.Get("userId").(string)

	eventId := c.QueryParam("event")
	pickerId := c.QueryParam("picker")
	collectionId := c.QueryParam("collectionId")
	asSociety, err := strconv.ParseBool(c.QueryParam("asSociety"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if asSociety {
		hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, pickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCannotAttendEvent, err))
		}
		if !hasRights {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	} else {
		pickerId = userId
	}

	err = s.eventAccess.DeleteCollectionOrganized(pickerId, collectionId, eventId, asSociety)

	return err
}
