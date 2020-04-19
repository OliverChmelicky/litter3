package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/olo/litter3/models"
)

type TrashAccess struct {
	Db *pg.DB
}

func (s *TrashAccess) CreateTrash(in *models.Trash) (*models.Trash, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &models.Trash{}, err
	}

	return in, nil
}

func (s *TrashAccess) GetTrash(id string) (*models.Trash, error) {
	trash := &models.Trash{Id: id}
	err := s.Db.Model(trash).Column("trash.*").
		Relation("Collections").
		Relation("Images").
		First()
	if err != nil {
		return &models.Trash{}, err
	}

	return trash, nil
}

func (s *TrashAccess) GetTrashInRange(request *models.RangeRequest) ([]models.Trash, error) {
	//https://postgis.net/docs/PostGIS_FAQ.html#idm1368
	trash := []models.Trash{}
	err := s.Db.Model(&trash).
		Column("trash.*").
		Where("ST_DWithin(location, 'SRID=4326;POINT(? ?)', ?)", request.Location[0], request.Location[1], request.Radius).
		Relation("Images").
		Select()
	if err != nil {
		return nil, err
	}
	return trash, nil
}

func (s *TrashAccess) UpdateTrash(in *models.Trash) (*models.Trash, error) {
	return in, s.Db.Update(in)
}

func (s *TrashAccess) DeleteTrash(trashId string) error {
	trash := new(models.Trash)
	trash.Id = trashId

	tx, err := s.Db.Begin()
	if err != nil {
		return fmt.Errorf("Coudln`t start transaction: %w ", err)
	}
	defer tx.Rollback()

	var collections []models.Collection
	err = tx.Model(&collections).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return fmt.Errorf("Error collections relelevant to trash: %w ", err)
	}
	if len(collections) != 0 {
		return fmt.Errorf("Error trash has some collections already ")
	}

	var eventTrash []models.EventTrash
	err = tx.Model(&eventTrash).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return fmt.Errorf("Error events relelevant to trash: %w ", err)
	}
	if len(collections) != 0 {
		return fmt.Errorf("Error events are organized to trash ")
	}

	err = s.DeleteTrashComments(trashId, tx)
	if err != nil {
		return err
	}

	//TODO delete images

	err = tx.Delete(trash)
	if err != nil {
		return fmt.Errorf("Error deleting traash: %w ", err)
	}

	return tx.Commit()
}

//
//
//	COLLECTION
//
//

func (s *TrashAccess) CreateCollectionRandom(in *models.CreateCollectionRandomRequest, creatorId string) (*models.Collection, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Collection{}, err
	}
	defer tx.Rollback()

	collection := &models.Collection{TrashId: in.TrashId, CleanedTrash: in.CleanedTrash, Weight: in.Weight}
	err = tx.Insert(collection)
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error creating collection: %w ", err)
	}

	userCollection := &models.UserCollection{}
	userCollection.UserId = creatorId
	userCollection.CollectionId = collection.Id
	err = tx.Insert(userCollection)
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error assigning creator to collection: %w ", err)
	}

	for _, picker := range in.Friends {
		userCollection.UserId = picker
		userCollection.CollectionId = collection.Id
		err = tx.Insert(userCollection)
		if err != nil {
			return &models.Collection{}, fmt.Errorf("Error assigning friends to collection: %w ", err)
		}
	}

	if in.CleanedTrash {
		trash := new(models.Trash)
		trash.Cleaned = true
		_, err = s.Db.Model(trash).Column("cleaned").Where("id = ?", in.TrashId).Update()
		if err != nil {
			return &models.Collection{}, err
		}
	}

	return collection, tx.Commit()
}

func (s *TrashAccess) GetCollection(id string) (*models.Collection, error) {
	collection := new(models.Collection)
	err := s.Db.Model(collection).Where("id = ? ", id).
		Relation("Images", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("collection_id = ?", id), nil
		}).Select()
	if err != nil {
		return &models.Collection{}, err
	}
	return collection, nil
}

func (s *TrashAccess) GetCollectionIdsOfUser(userId string) (*models.UserCollection, error) {
	userCollection := new(models.UserCollection)
	err := s.Db.Model(userCollection).Where("user_id = ?", userId).Select()
	if err != nil {
		return &models.UserCollection{}, err
	}
	return userCollection, nil
}

func (s *TrashAccess) UpdateCollectionRandom(request *models.Collection, userId string) (*models.Collection, error) {
	attended, err := s.isUserInCollection(request.Id, userId)
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error verifying if user is in collection: %w ", err)
	}
	if !attended {
		return &models.Collection{}, fmt.Errorf("You cannot update this collection: %w ", err)
	}

	oldCollection := new(models.Collection)
	oldCollection.Id = request.Id
	if err := s.Db.Select(oldCollection); err != nil {
		return &models.Collection{}, fmt.Errorf("Error getting old collection: %w ", err)
	}

	if oldCollection.TrashId != request.TrashId {
		return &models.Collection{}, fmt.Errorf("You cannot change TrashId: %w ", err)
	}

	if request.EventId != "" {
		return &models.Collection{}, fmt.Errorf("You cannot update EventId ")
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error creating transaction: %w ", err)
	}
	defer tx.Rollback()

	err = tx.Update(request)
	if err != nil {
		return &models.Collection{}, fmt.Errorf("Error update collection: %w ", err)
	}

	if request.CleanedTrash != oldCollection.CleanedTrash {
		trash := new(models.Trash)
		trash.Id = request.TrashId
		trash.Cleaned = request.CleanedTrash
		_, err = tx.Model(trash).Column("cleaned").Where("id = ?", request.TrashId).Update()
		if err != nil {
			return &models.Collection{}, err
		}
	}

	return request, tx.Commit()
}

func (s *TrashAccess) AddPickerToCollection(request *models.UserCollection, givesAnother string) (*models.UserCollection, error) {
	attended, err := s.isUserInCollection(request.CollectionId, givesAnother)
	if err != nil {
		return &models.UserCollection{}, fmt.Errorf("Error verifying if user is in collection: %w ", err)
	}
	if !attended {
		return &models.UserCollection{}, fmt.Errorf("You cannot update this collection: %w ", err)
	}

	err = s.Db.Insert(request)
	if err != nil {
		return &models.UserCollection{}, fmt.Errorf("Error adding picker to collection: %w ", err)
	}

	return request, nil
}

func (s *TrashAccess) DeleteCollectionFromUser(collectionId string, userId string) error {
	attended, err := s.isUserInCollection(collectionId, userId)
	if err != nil {
		return fmt.Errorf("Error verifying if user is in collection: %w ", err)
	}
	if !attended {
		return fmt.Errorf("You are not in this collection: %w ", err)
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return fmt.Errorf("Couldn`t create transaction for deletion: %w ", err)
	}
	defer tx.Rollback()

	userCollection := new(models.UserCollection)
	userCollection.CollectionId = collectionId
	userCollection.UserId = userId
	err = tx.Delete(userCollection)
	if err != nil {
		return fmt.Errorf("Error deleting user from collection: %w ", err)
	}

	//check if user was the last one
	err = tx.Model(userCollection).Where("collection_id = ?", collectionId).Select()
	if err == pg.ErrNoRows {
		collection := &models.Collection{Id: collectionId}
		if err = tx.Select(collection); err != nil {
			return fmt.Errorf("Error checking collection for cleaned property: %w ", err)
		}

		//change trash back to be not cleaned
		if collection.CleanedTrash {
			trash := new(models.Trash)
			trash.Cleaned = false
			_, err = tx.Model(trash).Column("cleaned").Where("id = ?", collection.TrashId).Update()
			if err != nil {
				return fmt.Errorf("Error reverting trash cleaned property: %w ", err)
			}
		}

		err = tx.Delete(collection)
		if err != nil {
			return fmt.Errorf("Error deleting collection: %w ", err)
		}
	} else if err != nil {
		return fmt.Errorf("Error querry deleting collection: %w ", err)
	}

	return tx.Commit()
}

func (s *TrashAccess) isUserInCollection(collectionId string, userId string) (bool, error) {
	userCollection := new(models.UserCollection)
	err := s.Db.Model(userCollection).Where("collection_id = ? and user_id = ?", collectionId, userId).
		Select()
	if err == pg.ErrNoRows {
		return false, fmt.Errorf("You are not a member of collection")
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

//
//
//
//	COMMENTS
//
//

func (s *TrashAccess) CreateTrashComment(in *models.TrashComment) (*models.TrashComment, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &models.TrashComment{}, fmt.Errorf("CREATE TRASH COMMENT: %w", err)
	}
	return in, nil
}

func (s *TrashAccess) GetTrashCommentById(trashId string) (*models.TrashComment, error) {
	comment := new(models.TrashComment)
	err := s.Db.Model(comment).Where("id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENT: %w", err)
	}
	return comment, nil
}

func (s *TrashAccess) GetTrashCommentByTrashId(trashId string) (*models.TrashComment, error) {
	comment := new(models.TrashComment)
	err := s.Db.Model(comment).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENT: %w", err)
	}
	return comment, nil
}

func (s *TrashAccess) GetTrashComments(trashId string) ([]models.TrashComment, error) {
	var comments []models.TrashComment
	err := s.Db.Model(&comments).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENTS: %w", err)
	}
	return comments, nil
}

func (s *TrashAccess) UpdateTrashComment(in *models.TrashComment) (*models.TrashComment, error) {
	return in, s.Db.Update(in)
}

func (s *TrashAccess) DeleteTrashComment(id string) error {
	comment := new(models.TrashComment)
	_, err := s.Db.Model(comment).Where("id = ?", id).Delete()
	return err
}

func (s *TrashAccess) DeleteTrashComments(trashId string, tx *pg.Tx) error {
	comment := new(models.TrashComment)
	_, err := tx.Model(comment).Where("trash_id = ?", trashId).Delete()
	return err
}

func (s *TrashAccess) DeleteUserComments(userId string) error {
	comment := new(models.TrashComment)
	_, err := s.Db.Model(comment).Where("user_id = ?", userId).Delete()
	return err
}
