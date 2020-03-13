package user

import (
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"time"
)

type userAccess struct {
	db *pg.DB
}

func (s *userAccess) CreateUser(in *UserModel) (*UserModel, error) {
	in.Id = uuid.NewV4().String()
	in.Created = time.Now().Unix()
	s.db.Model()
	err := s.db.Insert(in)
	if err != nil {
		return &UserModel{}, err
	}

	return in, nil
}

func (s *userAccess) GetUser(in string) (*UserModel, error) {
	return &UserModel{}, nil
}

func (s *userAccess) UpdateUser(in *UserModel) (*UserModel, error) {
	return in, nil
}

func (s *userAccess) DeleteUser(in string) error {
	return nil
}
