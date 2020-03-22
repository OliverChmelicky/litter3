package user

import (
	"errors"
	"github.com/go-pg/pg/v9"
)

type userAccess struct {
	db *pg.DB
}

func (s *userAccess) CreateUser(in *User) (*User, error) {
	err := s.db.Insert(in)
	if err != nil {
		return &User{}, err
	}

	return in, nil
}

func (s *userAccess) GetUser(in string) (*User, error) {
	user := new(User)
	err := s.db.Model(user).Where("id = ?", in).Select()
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (s *userAccess) UpdateUser(in *User) (*User, error) {
	return in, s.db.Update(in)
}

func (s *userAccess) AddApplicant(in *Applicant) (*Applicant, error) {
	err := s.db.Insert(in)
	if err != nil {
		return &Applicant{}, err
	}

	return in, err
}

func (s *userAccess) RemoveApplicationForMembership(in *Applicant) error {
	applicant := new(Applicant)
	_, err := s.db.Model(applicant).Where("user_id = ? and society_id = ?", in.UserId, in.SocietyId).Delete()
	return err
}

func (s *userAccess) DeleteUser(in string) error {
	return nil
}

//
//
//
//
//

func (s *userAccess) CreateSocietyWithAdmin(in *Society, adminId string) (*Society, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return &Society{}, err
	}
	defer tx.Rollback()

	err = tx.Insert(in)
	if err != nil {
		return &Society{}, err
	}

	admin := &Member{
		SocietyId:  in.Id,
		UserId:     adminId,
		Permission: "admin",
	}

	err = tx.Insert(admin)
	if err != nil {
		return &Society{}, err
	}

	return in, tx.Commit()
}

func (s *userAccess) GetSociety(in string) (*Society, error) {
	society := new(Society)
	err := s.db.Model(society).Where("id = ?", in).Select()
	if err != nil {
		return &Society{}, err
	}
	return society, nil
}

func (s *userAccess) GetSocietyAdmins(in string) ([]string, error) {
	members := new(Member)
	var admins []string
	err := s.db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", in).Select(&admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (s *userAccess) CountSocietyAdmins(in string) (int, error) {
	members := new(Member)
	num, err := s.db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", in).Count()
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (s *userAccess) UpdateSociety(in *Society) (*Society, error) {
	return in, s.db.Update(in)
}

func (s *userAccess) RemoveMember(userId, societyId string) error {
	member := new(Member)
	_, err := s.db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Delete()
	return err
}

func (s *userAccess) IsMember(userId, societyId string) (bool, error) {
	member := new(Member)
	err := s.db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Select()
	if err != nil {
		return false, err
	}

	//co ak nenajde? Hodi err?
	return true, nil
}

func (s *userAccess) DeleteSociety(in string) error {
	//transaction

	//remove comments
	//what with events
	//what with collected trash --> removedSociety

	return errors.New("Uninplemented")
}
