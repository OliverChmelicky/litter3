package user

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/olo/litter3/models"
	log "github.com/sirupsen/logrus"
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

func (s *UserAccess) GetUserById(id string) (*models.User, error) {
	user := new(models.User)
	user.Id = id

	err := s.Db.Model(user).Column("user.*").
		Relation("Societies").Relation("Collections").
		Where("id = ?", id).First()

	for i, collection := range user.Collections {
		var images []models.CollectionImage
		err = s.Db.Model(&images).Where("collection_id = ?", collection.Id).Select()
		if err != nil {
			log.Error(err)
			continue
		}
		user.Collections[i].Images = images
	}

	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (s *UserAccess) GetUsersByIds(ids []string) ([]models.User, error) {
	users := []models.User{}
	if len(ids) > 0 {
		err := s.Db.Model(&users).Where("id IN (?)", pg.In(ids)).Select()
		if err != nil {
			return []models.User{}, err
		}
	}
	if len(users) == 0 {
		return []models.User{}, fmt.Errorf("No record GetUsersByIds")
	}

	return users, nil
}

func (s *UserAccess) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := s.Db.Model(user).Where("email = ?", email).Select()
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (s *UserAccess) UpdateUser(in *models.User) (*models.User, error) {
	usr, err := s.GetUserById(in.Id)
	if err != nil {
		return &models.User{}, err
	}
	usr.FirstName = in.FirstName
	usr.LastName = in.LastName
	usr.Email = in.Email
	err = s.Db.Update(usr)

	return usr, err
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

func (s *UserAccess) DeleteUser(userId string) (string, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	user := new(models.User)
	err = tx.Model(user).Where("id = ?", userId).Select()
	if err != nil {
		return "", err
	}
	firebaseUid := user.Uid

	var members []models.Member
	err = tx.Model(&members).Where("user_id = ? and permission = ?", userId, "admin").Select()
	if err != nil {
		return firebaseUid, err
	}

	var testNumOfAdmins []models.Member
	var societies []string
	for _, member := range members {
		err = tx.Model(&testNumOfAdmins).Where("society_id = ? and permission = ?", member.SocietyId, "admin").Select()
		if err != nil {
			return firebaseUid, fmt.Errorf("Error check number of adims in society: %w ", err)
		}
		if len(testNumOfAdmins) == 1 {
			societies = append(societies, member.SocietyId)
		}
	}

	err = s.DeleteSocieties(societies, tx)

	_, err = tx.Model(user).Where("id = ?", userId).Delete()
	if err != nil {
		log.Error(err)
	}

	//return firebaseUid, nil
	return firebaseUid, tx.Commit()
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

func (s *UserAccess) GetSocietiesWithPaging(from, to int) ([]models.Society, int, error) {
	limit := to - from
	societies := []models.Society{}
	err := s.Db.Model(&societies).Order("created_at ASC").Select()
	if err != nil {
		return []models.Society{}, 0, err
	}

	if len(societies) < from {
		return []models.Society{}, 0, fmt.Errorf("No records starting from FROM ")
	}
	if len(societies[from:]) < limit {
		to = from + len(societies[from:])
	}
	length := len(societies)
	societies = societies[from:to]

	var ids []string
	for _, society := range societies {
		ids = append(ids, society.Id)
	}
	if len(ids) > 0 {
		err = s.Db.Model(&societies).Column("society.*").
			Relation("Users").
			Where("id IN (?)", pg.In(ids)).Select()
		if err != nil {
			return []models.Society{}, 0, err
		}
	}

	return societies, length, nil
}

func (s *UserAccess) GetSociety(id string) (*models.Society, error) {
	society := &models.Society{Id: id}
	err := s.Db.Model(society).Column("society.*").
		Relation("Users").Relation("Applicants").
		Relation("MemberRights").Relation("ApplicantsIds").
		Where("id = ?", id).First()
	if err != nil {
		return &models.Society{}, err
	}
	return society, nil
}

func (s *UserAccess) GetSocieties(ids []string) ([]models.Society, error) {
	societies := []models.Society{}
	if len(ids) > 0 {
		err := s.Db.Model(&societies).Column("society.*").
			Relation("Users").Relation("Applicants").
			Relation("MemberRights").Relation("ApplicantsIds").
			Where("id IN (?)", pg.In(ids)).Select()
		if err != nil {
			return []models.Society{}, err
		}
	}
	return societies, nil
}

func (s *UserAccess) GetUserSocieties(id string) ([]models.Society, error) {
	var societies []models.Society
	var tableBetween []models.Member

	err := s.Db.Model(&tableBetween).Where("user_id = ?", id).Select()
	if err != nil {
		return []models.Society{}, err
	}

	var societiesIds []string
	for _, relation := range tableBetween {
		societiesIds = append(societiesIds, relation.SocietyId)
	}

	if len(societiesIds) > 0 {
		err = s.Db.Model(&societies).Where("id IN (?)", pg.In(societies)).
			Where("id = ?", id).Select()
		if err != nil {
			return []models.Society{}, err
		}
	}

	return societies, nil
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

func (s *UserAccess) GetSocietyEditors(societyId string) ([]string, error) {
	members := new(models.Member)
	var admins []string
	err := s.Db.Model(members).Column("user_id").Where("(permission = 'admin' or permission = 'editor') and society_id = ? ", societyId).Select(&admins)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, pg.ErrNoRows
	}

	return admins, nil
}

func (s *UserAccess) GetEditableSocieties(userId string) ([]models.Society, error) {
	var memberships []models.Member
	err := s.Db.Model(&memberships).Where("(permission = 'admin' or permission = 'editor') and user_id = ? ", userId).Select()
	if err != nil {
		log.Error(err)
		return []models.Society{}, err
	}

	var societiesIds []string
	for _, m := range memberships {
		societiesIds = append(societiesIds, m.SocietyId)
	}

	var societies []models.Society
	if len(societiesIds) > 0 {
		err = s.Db.Model(&societies).Where("id IN (?)", pg.In(societiesIds)).Select()
		if err != nil {
			return []models.Society{}, err
		}
	}

	return societies, nil
}

func (s *UserAccess) GetSocietyMembers(societyId string) ([]models.Member, error) {
	var members []models.Member
	err := s.Db.Model(&members).Where("society_id = ? ", societyId).Select(&members)
	if err != nil {
		return nil, err
	}

	return members, nil
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
	society := new(models.Society)
	if err := s.Db.Model(society).Where("id = ?", in.Id).Select(); err != nil {
		return &models.Society{}, fmt.Errorf("Error find society: %w ", err)
	}

	return in, s.Db.Update(in)
}

func (s *UserAccess) GetSocietyRequests(societyId string) ([]models.Applicant, error) {
	var requests []models.Applicant
	if err := s.Db.Model(&requests).Where("society_id = ?", societyId).Select(); err != nil {
		return []models.Applicant{}, fmt.Errorf("Error find society applicants: %w ", err)
	}

	return requests, nil
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

func (s *UserAccess) ChangeUserRights(requests []models.Member, societyId string) ([]models.Member, error) {
	var members []models.Member
	tx, err := s.Db.Begin()
	if err != nil {
		return []models.Member{}, err
	}
	defer tx.Rollback()

	for _, request := range requests {
		member := new(models.Member)
		member.UserId = request.UserId
		member.SocietyId = societyId
		err = tx.Select(member)
		if err != nil {
			continue
		}

		member.Permission = request.Permission
		err = tx.Update(member)
		if err != nil {
			continue
		}

		members = append(members, *member)
	}

	return members, tx.Commit()
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

func (s *UserAccess) DeleteSociety(id string) error {
	tx, err := s.Db.Begin() //there is trigger running but I am not sure if it catches - I guess not
	if err != nil {
		return err
	}
	defer tx.Rollback()

	society := new(models.Society)
	_, err = tx.Model(society).Where("id = ?", id).Delete()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserAccess) DeleteSocieties(ids []string, tx *pg.Tx) error {
	for _, societyId := range ids {
		_, err := tx.Model(&models.Society{}).Where("id = ?", societyId).Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

//
//
//
//	FRIEDNSHIP PART
//
//

func (s *UserAccess) GetUserFriends(userId string) ([]models.Friends, error) {
	var friends []models.Friends
	err := s.Db.Model(&friends).Where("user1_id = ? OR user2_id = ?", userId, userId).Select()
	if err != nil {
		return []models.Friends{}, err
	}

	return friends, nil
}

func (s *UserAccess) GetUserFriendshipRequests(userId string) ([]models.FriendRequest, error) {
	//TODO pg.rec.nofound
	var requests []models.FriendRequest
	err := s.Db.Model(&requests).Where("user1_id = ? OR user2_id = ?", userId, userId).Select()
	if err != nil {
		return []models.FriendRequest{}, err
	}

	return requests, nil
}

func (s *UserAccess) AreFriends(friendship *models.Friends) (bool, error) {
	err := s.Db.Model(friendship).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", friendship.User1Id, friendship.User2Id, friendship.User2Id, friendship.User1Id).Limit(1).Select()
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
	_, err := s.Db.Model(application).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", request.User1Id, request.User2Id, request.User2Id, request.User1Id).Delete()
	return err
}

func (s *UserAccess) ConfirmFriendship(requesterId, acceptorId string) (*models.Friends, error) {
	request := &models.FriendRequest{}

	tx, err := s.Db.Begin()
	if err != nil {
		return &models.Friends{}, fmt.Errorf("Error creating transaction %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Model(request).Where("(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?)", requesterId, acceptorId, acceptorId, requesterId).Delete()
	if err != nil {
		return &models.Friends{}, fmt.Errorf("Error deleting Friendship request %w", err)
	}
	if res.RowsAffected() == 0 {
		return &models.Friends{}, fmt.Errorf("Request did not exist before %w ", err)
	}

	friendship := &models.Friends{User1Id: requesterId, User2Id: acceptorId}
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

// HELPERS
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

func (s *UserAccess) HasUserSocietyEditorRights(userId string, societyId string) (bool, int, error) {
	editors, err := s.GetSocietyEditors(societyId) //admins + edotors
	if err != nil {
		return false, 0, err
	}
	if len(editors) == 0 {
		return false, 0, pg.ErrNoRows
	}

	for _, adminId := range editors {
		if adminId == userId {
			return true, len(editors), nil
		}
	}
	return false, len(editors), nil
}
