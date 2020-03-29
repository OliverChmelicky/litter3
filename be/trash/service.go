package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"net/http"
)

type trashService struct {
	*trashAccess
}

func CreateService(db *pg.DB) *trashService {
	access := &trashAccess{db: db}
	return &trashService{access}
}

func (s *trashService) CreateTrash(c echo.Context) error {
	trash := new(Trash)
	if err := c.Bind(trash); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	trash, err := s.trashAccess.CreateTrash(trash)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) GetTrashById(c echo.Context) error {
	id := c.Param("id")

	trash, err := s.trashAccess.GetTrash(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with id does not exist")
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) GetTrashInRange(c echo.Context) error {
	request := new(RangeRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	fmt.Printf("%+v \n", request)

	trash, err := s.trashAccess.GetTrashInRange(request)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError("GetTrashInRangeAccess", err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) UpdateTrash(c echo.Context) error {
	trash := new(Trash)
	if err := c.Bind(trash); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	_, err := s.trashAccess.GetTrash(trash.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	trash, err = s.trashAccess.UpdateTrash(trash)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating trash")
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) DeleteTrash(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

//
//
//
//
//

func (s *trashService) CreateCollection(c echo.Context) error {
	user := new(Collection)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.trashAccess.CreateCollection(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNotImplemented, user)
}

func (s *trashService) GetCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *trashService) GetCollectionOfUser(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *trashService) GetCollectionOfSociety(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *trashService) UpdateCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *trashService) DeleteCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}
