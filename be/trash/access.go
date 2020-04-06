package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
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

	return collection, tx.Commit()
}

//from event
func (s *TrashAccess) CreateCollectionOrganized(in *Collection) (*Collection, error) {
	in.Id = uuid.NewV4().String()
	err := s.Db.Insert(in)
	if err != nil {
		return &Collection{}, err
	}

	return in, nil
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
