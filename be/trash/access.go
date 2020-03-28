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
	fmt.Println(in)
	if err != nil {
		return &Trash{}, err
	}

	return in, nil
}

func (s *trashAccess) GetTrash(in string) (*Trash, error) {
	trash := new(Trash)
	err := s.db.Select(trash)
	if err != nil {
		return &Trash{}, err
	}
	return trash, nil
}

func (s *trashAccess) UpdateTrash(in *Trash) (*Trash, error) {
	return in, nil
}

func (s *trashAccess) DeleteTrash(in string) error {
	return nil
}

func (s *trashAccess) CreateCollection(in *Collection) (*Collection, error) {
	in.Id = uuid.NewV4().String()
	err := s.db.Insert(in)
	if err != nil {
		return &Collection{}, err
	}

	return in, nil
}
