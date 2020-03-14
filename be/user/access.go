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
	user := new(UserModel)
	err := s.db.Model(user).Where("id = ?", in).Select()
	if err != nil {
		return &UserModel{}, err
	}
	return user, nil
}

func (s *userAccess) UpdateUser(in *UserModel) (*UserModel, error) {
	return in, nil
}

func (s *userAccess) DeleteUser(in string) error {
	return nil
}

func (s *userAccess) CreateSociety(in *SocietyModel) (*SocietyModel, error) {
	in.Id = uuid.NewV4().String()
	in.Created = time.Now().Unix()
	s.db.Model()
	err := s.db.Insert(in)
	if err != nil {
		return &SocietyModel{}, err
	}

	return in, nil
}

func (s *userAccess) GetSociety(in string) (*SocietyModel, error) {
	society := new(SocietyModel)
	err := s.db.Model(society).Where("id = ?", in).Select()
	if err != nil {
		return &SocietyModel{}, err
	}
	return society, nil
}

func (s *userAccess) UpdateSociety(in *SocietyModel) (*SocietyModel, error) {
	return in, nil
}

func (s *userAccess) DeleteSociety(in string) error {
	return nil
}
