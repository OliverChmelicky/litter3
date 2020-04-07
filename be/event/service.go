package event

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	"net/http"
)

type EventService struct {
	*eventAccess
	*user.UserAccess
	*trash.TrashAccess
}

func CreateService(db *pg.DB, userAccess *user.UserAccess, trashAccess *trash.TrashAccess) *EventService {
	access := &eventAccess{db: db}
	return &EventService{access, userAccess, trashAccess}
}

func (s *EventService) CreateEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(EventRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, request.SocietyId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateEvent, err))
		}
		if !isAdmin {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.UserId = userId
	newTrash, err := s.eventAccess.CreateEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateEvent, err))
	}

	return c.JSON(http.StatusOK, newTrash)
}
func (s *EventService) GetEvent(c echo.Context) error {
	eventId := c.Param("eventId")

	event, err := s.eventAccess.GetEvent(eventId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetEvent, err))
	}

	return c.JSON(http.StatusOK, event)
}

func (s *EventService) AttendEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(EventAttendanceRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, request.PickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrAttendEvent, err))
		}
		if !isAdmin {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.PickerId = userId
	trash, err := s.eventAccess.AttendEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrAttendEvent, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *EventService) CannotAttendEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(EventAttendanceRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, request.PickerId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCannotAttendEvent, err))
		}
		if !isAdmin {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.PickerId = userId
	trash, err := s.eventAccess.CannotAttendEvent(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCannotAttendEvent, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *EventService) EditEvent(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(EventRequest)
	err := c.Bind(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if request.AsSociety {
		isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, request.SocietyId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateEvent, err))
		}
		if !isAdmin {
			return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrInsufficientPermission, err))
		}
	}

	request.UserId = userId
	newTrash, err := s.eventAccess.EditEvent(request, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrUpdateEvent, err))
	}

	return c.JSON(http.StatusOK, newTrash)
}

func (s *EventService) EditEventRights(c echo.Context) error {

}

//
//func(s *EventService) DeleteEvent(c echo.Context) error{
//
//}
//
//func(s *EventService) GetSocietyEvents(c echo.Context) error{
//
//}
//
//func(s *EventService) GetUserEvents(c echo.Context) error{
//
//}
//
//func(s *EventService) CreateCollection(c echo.Context) error{
//
//}
