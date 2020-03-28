package user

import (
	"errors"
	"github.com/go-pg/pg/v9"
)

type UserAccess struct {
	Db *pg.DB
}

func (s *UserAccess) CreateUser(in *User) (*User, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &User{}, err
	}

	return in, nil
}

func (s *UserAccess) GetUser(in string) (*User, error) {
	user := new(User)
	err := s.Db.Model(user).Where("id = ?", in).Select()
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (s *UserAccess) UpdateUser(in *User) (*User, error) {
	return in, s.Db.Update(in)
}

func (s *UserAccess) AddApplicant(in *Applicant) (*Applicant, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &Applicant{}, err
	}

	return in, err
}

func (s *UserAccess) RemoveApplicationForMembership(in *Applicant) error {
	applicant := new(Applicant)
	_, err := s.Db.Model(applicant).Where("user_id = ? and society_id = ?", in.UserId, in.SocietyId).Delete()
	return err
}

func (s *UserAccess) DeleteUser(in string) error {
	return nil
}

//
//
//
//
//

func (s *UserAccess) CreateSocietyWithAdmin(in *Society, adminId string) (*Society, error) {
	tx, err := s.Db.Begin()
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

func (s *UserAccess) GetSociety(in string) (*Society, error) {
	society := new(Society)
	err := s.Db.Model(society).Where("id = ?", in).Select()
	if err != nil {
		return &Society{}, err
	}
	return society, nil
}

func (s *UserAccess) GetSocietyAdmins(societyId string) ([]string, error) {
	members := new(Member)
	var admins []string
	err := s.Db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", societyId).Select(&admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (s *UserAccess) CountSocietyAdmins(in string) (int, error) {
	members := new(Member)
	num, err := s.Db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", in).Count()
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (s *UserAccess) UpdateSociety(in *Society) (*Society, error) {
	return in, s.Db.Update(in)
}

func (s *UserAccess) AcceptApplicant(userId, societyId string) (*Member, error) {
	applicant := new(Applicant)
	err := s.Db.Model(applicant).Where("user_id = ? and society_id = ?", userId, societyId).Select()
	if err != nil {
		return &Member{}, err
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return &Member{}, err
	}
	defer tx.Rollback()

	err = tx.Delete(applicant)
	if err != nil {
		return &Member{}, err
	}

	newMember := &Member{UserId: userId, SocietyId: societyId, Permission: membership("member")}
	err = tx.Insert(newMember)
	if err != nil {
		return &Member{}, err
	}

	return newMember, tx.Commit()
}

func (s *UserAccess) ChangeUserRights(request *Member) (*Member, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return &Member{}, err
	}
	defer tx.Rollback()

	member := new(Member)
	member.UserId = request.UserId
	member.SocietyId = request.SocietyId
	err = tx.Select(member)
	if err != nil {
		return &Member{}, err
	}

	member.Permission = request.Permission
	err = tx.Update(member)
	if err != nil {
		return &Member{}, err
	}

	return member, tx.Commit()
}

func (s *UserAccess) RemoveMember(userId, societyId string) error {
	member := new(Member)
	_, err := s.Db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Delete()
	return err
}

func (s *UserAccess) IsMember(userId, societyId string) (bool, error) {
	member := new(Member)
	err := s.Db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Select()
	if err == pg.ErrNoRows { //record was not found
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserAccess) DeleteSociety(in string) error {
	//transaction

	//remove comments
	//what with events
	//what with collected trash --> removedSociety

	return errors.New("Uninplemented")
}
