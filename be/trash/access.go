package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	middlewareService "github.com/olo/litter3/middleware"
	"github.com/satori/go.uuid"
)

type TrashAccess struct {
	Db *pg.DB
}

func (s *TrashAccess) CreateTrash(in *Trash) (*Trash, error) {
	in.Id = uuid.NewV4().String()
	err := s.Db.Insert(in)
	if err != nil {
		return &Trash{}, err
	}

	return in, nil
}

func (s *TrashAccess) GetTrash(in string) (*Trash, error) {
	trash := &Trash{Id: in}
	err := s.Db.Select(trash)
	if err != nil {
		return &Trash{}, err
	}
	return trash, nil
}

func (s *TrashAccess) GetTrashInRange(request *RangeRequest) ([]Trash, error) {
	//https://postgis.net/docs/PostGIS_FAQ.html#idm1368
	trash := []Trash{}
	err := s.Db.Model(&trash).Where("ST_DWithin(location, 'SRID=4326;POINT(? ?)', ?)", request.Location[0], request.Location[1], request.Radius).Select()
	if err != nil {
		return nil, err
	}
	return trash, nil
}

func (s *TrashAccess) UpdateTrash(in *Trash) (*Trash, error) {
	return in, s.Db.Update(in)
}

func (s *TrashAccess) DeleteTrash(in string) error {
	return nil
}

//
//
//	COLLECTION
//
//

func (s *TrashAccess) CreateCollectionRandom(in *CreateCollectionRandomRequest) (*Collection, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return &Collection{}, err
	}
	defer tx.Rollback()

	collection := &Collection{TrashId: in.TrashId, CleanedTrash: in.CleanedTrash}
	err = tx.Insert(collection)
	if err != nil {
		return &Collection{}, err
	}

	userCollection := &UserCollection{}
	for _, picker := range in.UsersIds {
		userCollection.UserId = picker
		userCollection.CollectionId = collection.Id
		err = tx.Insert()
		if err != nil {
			return &Collection{}, err
		}
	}

	if in.CleanedTrash {
		trash := new(Trash)
		trash.Id = in.TrashId
		_, err = s.Db.Model(trash).Column("cleaned").Where("id = ?", in.TrashId).Update()
		if err != nil {
			return &Collection{}, err
		}
	}

	return collection, tx.Commit()
}

func (s *TrashAccess) GetCollection(id string) (*Collection, error) {
	s.Db.AddQueryHook(middlewareService.DbMiddleware{})
	collection := new(Collection)
	collection.Id = id
	err := s.Db.Select(collection)
	if err != nil {
		return &Collection{}, err
	}
	return collection, nil
}

func (s *TrashAccess) GetCollectionsOfUser(userId string) (*UserCollection, error) {
	userCollection := new(UserCollection)
	userCollection.UserId = userId
	err := s.Db.Model(userCollection).Where("user_id = ?", userId).Select()
	if err != nil {
		return &UserCollection{}, err
	}
	return userCollection, nil
}

func (s *TrashAccess) UpdateCollectionRandom(request *Collection, userId string) (*Collection, error) {
	//neriesi collection z eventu
	attended, err := s.isUserInCollection(request.Id, userId)
	if err != nil {
		return &Collection{}, fmt.Errorf("Error verifying if user is in collection: %w ", err)
	}
	if !attended {
		return &Collection{}, fmt.Errorf("You cannot update this collection: %w ", err)
	}

	err = s.Db.Update(request)
	if err != nil {
		return &Collection{}, fmt.Errorf("Error update collection: %w ", err)
	}

	if request.CleanedTrash {
		trash := new(Trash)
		trash.Id = request.TrashId
		_, err = s.Db.Model(trash).Column("cleaned").Where("id = ?", request.TrashId).Update()
		if err != nil {
			return &Collection{}, err
		}
	}

	return request, nil
}

func (s *TrashAccess) AddPickerToCollection(request *UserCollection, givesAnother string) (*UserCollection, error) {
	attended, err := s.isUserInCollection(request.CollectionId, givesAnother)
	if err != nil {
		return &UserCollection{}, fmt.Errorf("Error verifying if user is in collection: %w ", err)
	}
	if !attended {
		return &UserCollection{}, fmt.Errorf("You cannot update this collection: %w ", err)
	}

	err = s.Db.Insert(request)
	if err != nil {
		return &UserCollection{}, fmt.Errorf("Error adding picker to collection: %w ", err)
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

	userCollection := new(UserCollection)
	userCollection.CollectionId = collectionId
	userCollection.UserId = userId
	err = s.Db.Delete(userCollection)
	if err != nil {
		return fmt.Errorf("Error deleting from collection: %w ", err)
	}

	return nil
}

func (s *TrashAccess) isUserInCollection(collectionId string, userId string) (bool, error) {
	userCollection := new(UserCollection)
	userCollection.CollectionId = collectionId
	userCollection.UserId = userId
	err := s.Db.Select(userCollection)
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

func (s *TrashAccess) CreateTrashComment(in *TrashComment) (*TrashComment, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &TrashComment{}, fmt.Errorf("CREATE TRASH COMMENT: %w", err)
	}
	return in, nil
}

func (s *TrashAccess) GetTrashCommentById(trashId string) (*TrashComment, error) {
	comment := new(TrashComment)
	err := s.Db.Model(comment).Where("id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENT: %w", err)
	}
	return comment, nil
}

func (s *TrashAccess) GetTrashCommentByTrashId(trashId string) (*TrashComment, error) {
	comment := new(TrashComment)
	err := s.Db.Model(comment).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENT: %w", err)
	}
	return comment, nil
}

func (s *TrashAccess) GetTrashComments(trashId string) ([]TrashComment, error) {
	var comments []TrashComment
	err := s.Db.Model(&comments).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENTS: %w", err)
	}
	return comments, nil
}

func (s *TrashAccess) UpdateTrashComment(in *TrashComment) (*TrashComment, error) {
	return in, s.Db.Update(in)
}

func (s *TrashAccess) DeleteTrashComment(id string) error {
	comment := new(TrashComment)
	_, err := s.Db.Model(comment).Where("id = ?", id).Delete()
	return err
}

//will I need it?
func (s *TrashAccess) DeleteTrashComments(trashId string) error {
	comment := new(TrashComment)
	_, err := s.Db.Model(comment).Where("trash_id = ?", trashId).Delete()
	return err
}

func (s *TrashAccess) DeleteUserComments(userId string) error {
	comment := new(TrashComment)
	_, err := s.Db.Model(comment).Where("user_id = ?", userId).Delete()
	return err
}
