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
	request := new(TrashCommentRequest)
	if err := c.Bind(TrashCommentRequest{}); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	comment, err := s.trashAccess.GetTrashComment(request.TrashId)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	comment.message = request.message
	comment, err = s.trashAccess.UpdateTrashComment(comment)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating trash")
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *trashService) DeleteTrash(c echo.Context) error {
	//cez param
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

//
//
//	TRASH COMMENT
//
//

func (s *trashService) CreateComment(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(TrashCommentRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	comment, err := s.trashAccess.CreateTrashComment(&TrashComment{TrashId: request.TrashId, UserId: userId, message: request.message})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateComment, err))
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *trashService) GetTrashComments(c echo.Context) error {
	trashId := c.Param("id")

	comments, err := s.trashAccess.GetTrashComments(trashId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetTrash, err))
	}

	return c.JSON(http.StatusOK, comments)
}

func (s *trashService) UpdateComment(c echo.Context) error {
	id := c.Param("id")

	trash, err := s.trashAccess.GetTrash(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with id does not exist")
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) DeleteComment(c echo.Context) error {
	id := c.Param("id")

	trash, err := s.trashAccess.GetTrash(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with id does not exist")
	}

	return c.JSON(http.StatusOK, trash)
}

//
//
//
//
//
//

func (s *trashService) CreateCollection(c echo.Context) error {
	request := new(CreateCollectionRandomRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	collection, err := s.trashAccess.CreateCollectionRandom(request)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNotImplemented, collection)
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
