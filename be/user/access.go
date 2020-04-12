package user

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/olo/litter3/models"
)

type UserAccess struct {
	Db *pg.DB
}

//
//
//
//	USER PART
//
//

func (s *UserAccess) CreateUser(in *models.User) (*models.User, error) {
	err := s.Db.Insert(in)
	if err != nil {
		return &models.User{}, err
	}

	return in, nil
}

func (s *UserAccess) GetUser(id string) (*models.User, error) {
	user := new(models.User)
	user.Id = id
	err := s.Db.Model(user).Select()
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (s *UserAccess) UpdateUser(in *models.User) (*models.User, error) {
	if _, err := s.GetUser(in.Id); err != nil {
		return &models.User{}, err
	}
	return in, s.Db.Update(in)
}

func (s *UserAccess) AddApplicant(in *models.Applicant) (*models.Applicant, error) {
	society := new(models.Society)
	err := s.Db.Model(society).Where("id = ?", in.SocietyId).Select()
	if err != nil {
		return &models.Applicant{}, fmt.Errorf("ERROR FIND SOCIETY %w", err)
	}

	err = s.Db.Insert(in)
	if err != nil {
		return &models.Applicant{}, fmt.Errorf("ERROR INSERT APPLICATION %w", err)
	}

	return in, nil
}

func (s *UserAccess) RemoveApplicationForMembership(in *models.Applicant) error {
	applicant := new(models.Applicant)
	_, err := s.Db.Model(applicant).Where("user_id = ? and society_id = ?", in.UserId, in.SocietyId).Delete()
	return err
}

func (s *UserAccess) DeleteUser(in string) error {
	return nil
}

//
//
//
//	SOCIETY PART
//
//

func (s *UserAccess) CreateSocietyWithAdmin(in *models.Society, adminId string) (*models.Society, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Society{}, err
	}
	defer tx.Rollback()

	err = tx.Insert(in)
	if err != nil {
		return &models.Society{}, err
	}

	admin := &models.Member{
		SocietyId:  in.Id,
		UserId:     adminId,
		Permission: "admin",
	}

	err = tx.Insert(admin)
	if err != nil {
		return &models.Society{}, err
	}

	return in, tx.Commit()
}

func (s *UserAccess) GetSociety(id string) (*models.Society, error) {
	society := &models.Society{Id: id}
	err := s.Db.Model(society).Column("society.*").
		Relation("Users").
		Where("id = ?", id).First()
	if err != nil {
		return &models.Society{}, err
	}
	return society, nil
}

//worht thinking more how to explicitly involve my friends who are in the wanted society
//func (s *UserAccess) GetMyFriendsInSociety(societyId, userId string) ([]models.User, error) {
//	friends := []models.Friends{}
//	err := s.Db.Model(&friends).Where("user1_id = ? or user2_id = ?", userId, userId).
//		Select()
//	if err != nil {
//		return nil, fmt.Errorf("Error querying friends of user: %w ", err)
//	}
//
//	socMemb := []models.User{}
//	err := s.Db.Model(&socMemb).Column("societies_members.*").Where("society_id = ?", societyId).
//		Relation("UserDetails").
//		Select()
//	if err != nil {
//		return nil, fmt.Errorf("Error querying normal members view: %w ", err)
//	}
//
//
//}

func (s *UserAccess) GetSocietyAdminsAll(societyId string) ([]models.User, error) {
	var admins []models.User
	permission := models.Membership("admin")
	err := s.Db.Model(&admins).Column("user.*").
		Relation("Admins", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("permission = ? and society_id = ?", permission, societyId), nil
		}).
		Select()
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, pg.ErrNoRows
	}

	return admins, nil
}

func (s *UserAccess) GetSocietyAdmins(societyId string) ([]string, error) {
	members := new(models.Member)
	var admins []string
	err := s.Db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", societyId).Select(&admins)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, pg.ErrNoRows
	}

	return admins, nil
}

func (s *UserAccess) CountSocietyAdmins(in string) (int, error) {
	members := new(models.Member)
	num, err := s.Db.Model(members).Column("user_id").Where("permission = 'admin' and society_id = ? ", in).Count()
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (s *UserAccess) UpdateSociety(in *models.Society) (*models.Society, error) {
	society := &models.Society{Id: in.Id}
	if err := s.Db.Model(society).Select(); err != nil {
		return &models.Society{}, fmt.Errorf("Error find society: %w ", err)
	}

	return in, s.Db.Update(in)
}

func (s *UserAccess) AcceptApplicant(userId, societyId string) (*models.Member, error) {
	applicant := new(models.Applicant)
	err := s.Db.Model(applicant).Where("user_id = ? and society_id = ?", userId, societyId).Select()
	if err != nil {
		return &models.Member{}, err
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Member{}, err
	}
	defer tx.Rollback()

	err = tx.Delete(applicant)
	if err != nil {
		return &models.Member{}, err
	}

	newMember := &models.Member{UserId: userId, SocietyId: societyId, Permission: models.Membership("member")}
	err = tx.Insert(newMember)
	if err != nil {
		return &models.Member{}, err
	}

	return newMember, tx.Commit()
}

func (s *UserAccess) ChangeUserRights(request *models.Member) (*models.Member, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Member{}, err
	}
	defer tx.Rollback()

	member := new(models.Member)
	member.UserId = request.UserId
	member.SocietyId = request.SocietyId
	err = tx.Select(member)
	if err != nil {
		return &models.Member{}, err
	}

	member.Permission = request.Permission
	err = tx.Update(member)
	if err != nil {
		return &models.Member{}, err
	}

	return member, tx.Commit()
}

func (s *UserAccess) RemoveMember(userId, societyId string) error {
	member := new(models.Member)
	_, err := s.Db.Model(member).Where("user_id = ? and society_id = ?", userId, societyId).Delete()
	return err
}

func (s *UserAccess) IsMember(userId, societyId string) (bool, error) {
	member := new(models.Member)
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
	//TODO
	//transaction

	//what with events
	//what with collected trash --> removedSociety

	return errors.New("Uninplemented")
}

//
//
//
//	FRIEDNSHIP PART
//
//

func (s *UserAccess) AreFriends(friendship *models.Friends) (bool, error) {
	err := s.Db.Model(friendship).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", friendship.User1Id, friendship.User2Id, friendship.User2Id, friendship.User1Id).Select()
	if err == pg.ErrNoRows { //record was not found
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserAccess) IsFriendRequestSendAlready(request *models.FriendRequest) (bool, error) {
	err := s.Db.Model(request).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", request.User1Id, request.User2Id, request.User2Id, request.User1Id).Select()
	if err == pg.ErrNoRows { //record was not found
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserAccess) AddFriendshipRequest(request *models.FriendRequest) (*models.FriendRequest, error) {
	err := s.Db.Insert(request)
	if err != nil {
		return &models.FriendRequest{}, err
	}
	return request, nil
}

func (s *UserAccess) RemoveApplicationForFriendship(request *models.FriendRequest) error {
	application := new(models.FriendRequest)
	res, err := s.Db.Model(application).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", request.User1Id, request.User2Id, request.User2Id, request.User1Id).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("REMOVE APPLICATION FOR FRIENDSHIP: NO ROWS WERE AFFECTED")
	}
	return err
}

func (s *UserAccess) ConfirmFriendship(user1Id, user2Id string) (*models.Friends, error) {
	request := &models.FriendRequest{User1Id: user1Id, User2Id: user2Id}

	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Friends{}, fmt.Errorf("Error creating transaction %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Model(request).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", request.User1Id, request.User2Id, request.User2Id, request.User1Id).Delete()
	if err != nil {
		return &models.Friends{}, fmt.Errorf("Error deleting Friendship request %w", err)
	}

	friendship := &models.Friends{User1Id: request.User1Id, User2Id: request.User2Id}
	err = s.Db.Insert(friendship)
	if err != nil {
		return &models.Friends{}, fmt.Errorf("Error creating Friendship %w", err)
	}
	return friendship, tx.Commit()
}

func (s *UserAccess) RemoveFriend(friendship *models.Friends) error {
	application := new(models.Friends)
	res, err := s.Db.Model(application).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", friendship.User1Id, friendship.User2Id, friendship.User2Id, friendship.User1Id).Delete()
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("REMOVE FRIEND: NO ROWS WERE AFFECTED")
	}
	return err
}

func (s *UserAccess) IsUserSocietyAdmin(userId, societyId string) (bool, int, error) {
	admins, err := s.GetSocietyAdmins(societyId)
	if err != nil {
		return false, 0, err
	}
	if len(admins) == 0 {
		return false, 0, pg.ErrNoRows
	}

	for _, adminId := range admins {
		if adminId == userId {
			return true, len(admins), nil
		}
	}
	return false, len(admins), nil
}
