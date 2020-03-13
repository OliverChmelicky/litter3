package user

import (
	"firebase.google.com/go/auth"
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
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DeleteUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) CreateSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) GetSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) UpdateSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DeleteSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}
