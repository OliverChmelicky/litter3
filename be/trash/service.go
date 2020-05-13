package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	"net/http"
	"strconv"
)

type trashService struct {
	*TrashAccess
}

func CreateService(db *pg.DB) *trashService {
	access := &TrashAccess{Db: db}
	return &trashService{access}
}

func (s *trashService) CreateTrash(c echo.Context) error {
	creatorId, ok := c.Get("userId").(string)
	if !ok {
		creatorId = ""
	}

	trashRequest := new(models.CreateTrashRequest)
	if err := c.Bind(trashRequest); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	trash := new(models.Trash)
	if !trashRequest.Anonymously {
		trash.FinderId = creatorId
	} else {
		trash.FinderId = ""
	}

	trash.Location = trashRequest.Location
	trash.Size = trashRequest.Size
	trash.TrashType = trashRequest.TrashType
	trash.Description = trashRequest.Description
	trash.Accessibility = trashRequest.Accessibility

	trash, err := s.TrashAccess.CreateTrash(trash)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateTrash, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) GetTrashById(c echo.Context) error {
	id := c.Param("id")

	trash, err := s.TrashAccess.GetTrash(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetTrashById, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) GetTrashInRange(c echo.Context) error {
	var err error
	request := new(models.RangeRequest)

	lat, err := strconv.ParseFloat(c.QueryParam("lat"), 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	lng, err := strconv.ParseFloat(c.QueryParam("lng"), 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	radius, err := strconv.ParseFloat(c.QueryParam("radius"), 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	fmt.Println("Radius v metroch")
	fmt.Println(radius)

	request.Location = models.Point{lat, lng}
	request.Radius = radius

	trash, err := s.TrashAccess.GetTrashInRange(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetTrashInRange, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) UpdateTrash(c echo.Context) error {
	trash := new(models.Trash)
	if err := c.Bind(trash); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	trash, err := s.TrashAccess.UpdateTrash(trash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateTrash, err))
	}

	return c.JSON(http.StatusOK, trash)
}

func (s *trashService) DeleteTrash(c echo.Context) error {
	trashId := c.Param("trashId")

	err := s.TrashAccess.DeleteTrash(trashId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteTrash, err))
	}

	return c.NoContent(http.StatusOK)
}

//
//
//
//	TRASH COMMENT
//
//

func (s *trashService) CreateTrashComment(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.TrashCommentRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	comment, err := s.TrashAccess.CreateTrashComment(&models.TrashComment{TrashId: request.Id, UserId: userId, Message: request.Message})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateComment, err))
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *trashService) GetTrashCommentsByTrashId(c echo.Context) error {
	trashId := c.Param("trashId")

	comments, err := s.TrashAccess.GetTrashComments(trashId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetTrashComments, err))
	}

	return c.JSON(http.StatusOK, comments)
}

func (s *trashService) UpdateTrashComment(c echo.Context) error {
	request := new(models.TrashCommentRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	comment, err := s.TrashAccess.GetTrashCommentById(request.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetComment, err))
	}

	comment.Message = request.Message
	comment, err = s.TrashAccess.UpdateTrashComment(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateComment, err))
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *trashService) DeleteTrashComment(c echo.Context) error {
	userId := c.Get("userId")
	commentId := c.Param("commentId")

	comment, err := s.TrashAccess.GetTrashCommentById(commentId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetComment, err))
	}

	if comment.UserId != userId {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrDeleteComment, fmt.Errorf("You have to be a creator of comment")))
	}

	err = s.TrashAccess.DeleteTrashComment(commentId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteComment, err))
	}

	return c.NoContent(http.StatusOK)
}

//
//
//
//	COLLECTION
//
//

func (s *trashService) CreateCollection(c echo.Context) error {
	creator := c.Get("userId").(string)

	request := new(models.CreateCollectionRandomRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	collection, err := s.TrashAccess.CreateCollectionRandom(request, creator)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionRaw, err))
	}

	return c.JSON(http.StatusOK, collection)
}

func (s *trashService) GetCollection(c echo.Context) error {
	collectionId := c.Param("collectionId")

	collection, err := s.TrashAccess.GetCollection(collectionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionRaw, err))
	}

	return c.JSON(http.StatusOK, collection)
}

func (s *trashService) GetCollectionIdsOfUser(c echo.Context) error {
	userId := c.Param("userId")

	collection, err := s.TrashAccess.GetCollectionIdsOfUser(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionRaw, err))
	}

	return c.JSON(http.StatusOK, collection)
}

func (s *trashService) UpdateCollectionRandom(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.Collection)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	collection, err := s.TrashAccess.UpdateCollectionRandom(request, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateCollection, err))
	}

	return c.JSON(http.StatusOK, collection)
}

func (s *trashService) AddPickerToCollection(c echo.Context) error {
	userId := c.Get("userId").(string)
	request := new(models.UserCollection)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	collection, err := s.TrashAccess.AddPickerToCollection(request, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionRaw, err))
	}

	return c.JSON(http.StatusOK, collection)
}

func (s *trashService) DeleteCollectionFromUser(c echo.Context) error {
	userId := c.Get("userId").(string)
	collectionId := c.Param("collectionId")

	err := s.TrashAccess.DeleteCollectionFromUser(collectionId, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateCollectionRaw, err))
	}

	return c.NoContent(http.StatusOK)
}
