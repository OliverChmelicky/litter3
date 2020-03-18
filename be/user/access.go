package user

import (
	"errors"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type userAccess struct {
	db *pg.DB
}

func (s *userAccess) CreateUser(in *UserModel) (*UserModel, error) {
	in.Id = uuid.NewV4().String()
	in.Created = time.Now().Unix()
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

func (s *userAccess) AddApplicant(in *ApplicantModel) (*ApplicantModel, error) {
	in.Created = time.Now().Unix()
	s.db.Model()
	err := s.db.Insert(in)
	if err != nil {
		return &ApplicantModel{}, err
	}

	return in, err
}

func (s *userAccess) RemoveApplicationForMembership(in *ApplicantModel) error {
	applicant := new(ApplicantModel)
	_, err := s.db.Model(applicant).Where("user_id = ? and society_id = ?", in.UserId, in.SocietyId).Delete()
	return err
}

func (s *userAccess) DeleteUser(in string) error {
	return nil
}

func (s *userAccess) CreateSocietyWithAdmin(in *SocietyModel, adminId string) (*SocietyModel, error) {
	//Make it transactional
	//put Id creation, created into db middleware
	in.Id = uuid.NewV4().String()
	in.Created = time.Now().Unix()

	err := s.db.Insert(in)
	if err != nil {
		return &SocietyModel{}, err
	}

	admin := &MemberModel{

		SocietyId:  in.Id,
		UserId:     adminId,
		Permission: "admin",
	}

	err = s.db.Insert(admin)
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

func (s *userAccess) GetSocietyAdmins(in string) ([]string, error) {
	members := new(MemberModel)
	var admins []string
	//error: pg: Model(non-pointer []string)
	err := s.db.Model(members).Column("user_id").Where("membership = admin and society_id = ? ", in).Select(admins) //does it work like this, is slice pointer itself?
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.Info(admins)
	return admins, nil
}

func (s *userAccess) UpdateSociety(in *SocietyModel) (*SocietyModel, error) {
	return in, nil
}

func (s *userAccess) RemoveMember(userId, societyId string) error {
	member := new(MemberModel)
	_, err := s.db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Delete()
	return err
}

func (s *userAccess) IsMember(userId, societyId string) (bool, error) {
	member := new(MemberModel)
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
