package trash

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
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
	trash := new(TrashModel)
	if err := c.Bind(trash); err != nil {
		log.Error(trash)
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Info(trash)
	log.Warn(trash.Gps)
	log.Warn(trash.Gps.X)
	log.Warn(trash.Gps.Y)

	trash, err := s.trashAccess.CreateTrash(trash)
	if err != nil {
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

func (s *trashService) GetTrashByInRange(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *trashService) UpdateTrash(c echo.Context) error {
	trash := new(TrashModel)
	if err := c.Bind(trash); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	_, err := s.trashAccess.GetTrash(trash.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Invalid arguments")
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
