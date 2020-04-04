package trash

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/satori/go.uuid"
)

type trashAccess struct {
	db *pg.DB
}

func (s *trashAccess) CreateTrash(in *Trash) (*Trash, error) {
	in.Id = uuid.NewV4().String()
	err := s.db.Insert(in)
	if err != nil {
		return &Trash{}, err
	}

	return in, nil
}

func (s *trashAccess) GetTrash(in string) (*Trash, error) {
	trash := &Trash{Id: in}
	err := s.db.Select(trash)
	if err != nil {
		return &Trash{}, err
	}
	return trash, nil
}

func (s *trashAccess) GetTrashInRange(request *RangeRequest) ([]Trash, error) {
	//https://postgis.net/docs/PostGIS_FAQ.html#idm1368
	var trash []Trash
	err := s.db.Model(trash).Where("ST_DWithin(location, 'SRID=4326;POINT(? ?)', ?)", request.Location[0], request.Location[1], request.Radius).Select()
	if err != nil {
		return nil, err
	}
	return trash, nil
}

func (s *trashAccess) UpdateTrash(in *Trash) (*Trash, error) {
	return in, s.db.Update(in)
}

func (s *trashAccess) DeleteTrash(in string) error {
	return nil
}

//
//
//	COLLECTION
//
//

func (s *trashAccess) CreateCollectionRandom(in *CreateCollectionRandomRequest) (*Collection, error) {
	tx, err := s.db.Begin()
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
func (s *trashAccess) CreateCollectionOrganized(in *Collection) (*Collection, error) {
	in.Id = uuid.NewV4().String()
	err := s.db.Insert(in)
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

func (s *trashAccess) CreateTrashComment(in *TrashComment) (*TrashComment, error) {
	err := s.db.Insert(in)
	if err != nil {
		return &TrashComment{}, fmt.Errorf("CREATE TRASH COMMENT: %w", err)
	}
	return in, nil
}

func (s *trashAccess) GetTrashComment(trashId string) (*TrashComment, error) {
	comment := new(TrashComment)
	err := s.db.Model(comment).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENT: %w", err)
	}
	return comment, nil
}

func (s *trashAccess) GetTrashComments(trashId string) ([]TrashComment, error) {
	var comments []TrashComment
	err := s.db.Model(comments).Where("trash_id = ?", trashId).Select()
	if err != nil {
		return nil, fmt.Errorf("GET TRASH COMMENTS: %w", err)
	}
	return comments, nil
}

func (s *trashAccess) UpdateTrashComment(in *TrashComment) (*TrashComment, error) {
	return in, s.db.Update(in)
}

func (s *trashAccess) DeleteTrashComment(id string) (interface{}, interface{}) {

}

//will I need it?
func (s *trashAccess) DeleteTrashComments(id string) (interface{}, interface{}) {

}
