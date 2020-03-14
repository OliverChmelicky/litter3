package user

import (
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"net/http"
)

type userService struct {
	userAccess   *userAccess
	firebaseAuth *auth.Client
}

func CreateService(db *pg.DB, authConn *auth.Client) *userService {
	access := &userAccess{db: db}
	return &userService{access, authConn}
}

func (s *userService) CreateUser(c echo.Context) error {
	user := new(UserModel)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.userAccess.CreateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNotImplemented, user)
}

func (s *userService) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := s.userAccess.GetUser(id)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	id := c.Get("userId")
	strId := fmt.Sprintf("%v", id)

	user, err := s.userAccess.GetUser(strId)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) UpdateUser(c echo.Context) error {
	user := new(UserModel)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	_, err := s.userAccess.GetUser(user.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	user, err = s.userAccess.UpdateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating user")
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) ApplyForMembership(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DeleteUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) CreateSociety(c echo.Context) error {
	society := new(SocietyModel)
	if err := c.Bind(society); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	society, err := s.userAccess.CreateSociety(society)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, society)
}

func (s *userService) GetSociety(c echo.Context) error {
	id := c.Param("id")

	user, err := s.userAccess.GetSociety(id)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Society with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) UpdateSociety(c echo.Context) error {
	user := new(UserModel)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	_, err := s.userAccess.GetUser(user.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	user, err = s.userAccess.UpdateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating user")
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) AcceptApplicant(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DismissApplicant(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) RemoveMember(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DeleteSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}
