package trash

import (
	"github.com/go-pg/pg/v9"
	"github.com/satori/go.uuid"
	"time"
)

type trashAccess struct {
	db *pg.DB
}

func (s *trashAccess) CreateTrash(in *TrashModel) (*TrashModel, error) {
	in.Id = uuid.NewV4().String()
	in.Created = time.Now().Unix()
	s.db.Model()
	err := s.db.Insert(in)
	if err != nil {
		return &TrashModel{}, err
	}

	return in, nil
}

func (s *trashAccess) GetTrash(in string) (*TrashModel, error) {
	trash := new(TrashModel)
	err := s.db.Model(trash).Where("id = ?", in).Select()
	if err != nil {
		return &TrashModel{}, err
	}
	return trash, nil
}

func (s *trashAccess) UpdateTrash(in *TrashModel) (*TrashModel, error) {
	return in, nil
}

func (s *trashAccess) DeleteTrash(in string) error {
	return nil
}
