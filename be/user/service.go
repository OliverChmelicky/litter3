package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	log "github.com/sirupsen/logrus"
)

type userService struct {
	UserAccess *UserAccess
	Firebase   *auth.Client
}

func CreateService(db *pg.DB, firebase *auth.Client) *userService {
	access := &UserAccess{Db: db}
	return &userService{access, firebase}
}

func (s *userService) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateUser, err))
	}

	user, err := s.UserAccess.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateUser, err))
	}

	claims := map[string]interface{}{
		"userId": user.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = s.Firebase.SetCustomUserClaims(ctx, user.Uid, claims)
	if err != nil {
		_, errDel := s.UserAccess.DeleteUser(user.Id)
		go func() {
			err = s.Firebase.DeleteUser(context.Background(), user.Uid)
			if err != nil {
				log.Error(err)
			}
		}()
		err = fmt.Errorf(err.Error()+" ERROR user deleted %w", errDel)
		return c.JSON(http.StatusGatewayTimeout, custom_errors.WrapError(custom_errors.ErrCreateUser, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := s.UserAccess.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUserById, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetUserByEmail(c echo.Context) error {
	email := c.Param("email")

	user, err := s.UserAccess.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUserByEmail, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetUsers(c echo.Context) error {
	idsString := c.QueryParam("ids")
	ids := strings.Split(idsString, ",")

	users, err := s.UserAccess.GetUsersByIds(ids)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUsers, err))
	}

	return c.JSON(http.StatusOK, users)
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	id := c.Get("userId").(string)

	user, err := s.UserAccess.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetCurrentUser, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) UpdateUser(c echo.Context) error {
	callerId := c.Get("userId").(string)

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
	}

	user.Id = callerId

	user, err := s.UserAccess.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) ApplyForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(models.IdMessage)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}

	isMember, err := s.UserAccess.IsMember(requesterId, request.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}
	if isMember {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("User is already a member")))
	}

	applicant, err := s.UserAccess.AddApplicant(&models.Applicant{SocietyId: request.Id, UserId: requesterId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}

	return c.JSON(http.StatusOK, applicant)
}

func (s *userService) RemoveApplicationForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	societyId := c.Param("societyId")

	err := s.UserAccess.RemoveApplicationForMembership(&models.Applicant{UserId: requesterId, SocietyId: societyId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveApplicationForMembership, err))
	}

	return c.NoContent(http.StatusOK)
}

func (s *userService) DeleteUser(c echo.Context) error {
	userId := c.Get("userId").(string)

	firebaseUid, err := s.UserAccess.DeleteUser(userId)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteUser, err))
	}

	err = s.Firebase.DeleteUser(context.Background(), firebaseUid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteUser, err))
	}

	return c.NoContent(http.StatusOK)
}

//
//
//
//	FRIENDS PART
//
//

func (s *userService) GetMyFriends(c echo.Context) error {
	userId := c.Get("userId").(string)

	requests, err := s.UserAccess.GetUserFriends(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUserFriends, err))
	}

	return c.JSON(http.StatusOK, requests)
}

func (s *userService) GetMyFriendRequests(c echo.Context) error {
	userId := c.Get("userId").(string)

	requests, err := s.UserAccess.GetUserFriendshipRequests(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetMyReqForFriendship, err))
	}

	return c.JSON(http.StatusOK, requests)
}

func (s *userService) ApplyForFriendshipById(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(models.IdMessage)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	if requesterId == request.Id {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("YOU CANNOT BE FRIEND WITH YOURSELF")))
	}
	areFriendsAlready, err := s.UserAccess.AreFriends(&models.Friends{User1Id: requesterId, User2Id: request.Id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}
	if areFriendsAlready {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("YOU ARE FIENDS ALREADY")))
	}

	isrequestAlreadySend, err := s.UserAccess.IsFriendRequestSendAlready(&models.FriendRequest{User1Id: requesterId, User2Id: request.Id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}
	if isrequestAlreadySend {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("REQUEST IS SEND ALREADY")))
	}

	friendRequest := &models.FriendRequest{User1Id: requesterId, User2Id: request.Id}
	applicant, err := s.UserAccess.AddFriendshipRequest(friendRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}

	return c.JSON(http.StatusCreated, applicant)
}

func (s *userService) ApplyForFriendshipByEmail(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(models.EmailMessage)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	wantsToBeFriendWith, err := s.UserAccess.GetUserByEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUserById, errors.New("YOU ARE FIENDS ALREADY")))
	}

	if requesterId == wantsToBeFriendWith.Id {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("YOU CANNOT BE FRIEND WITH YOURSELF")))
	}
	areFriends, err := s.UserAccess.AreFriends(&models.Friends{User1Id: requesterId, User2Id: wantsToBeFriendWith.Id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}
	if areFriends {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("YOU ARE FIENDS ALREADY")))
	}

	isFriendRequestSendAlready, err := s.UserAccess.IsFriendRequestSendAlready(&models.FriendRequest{User1Id: requesterId, User2Id: wantsToBeFriendWith.Id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}
	if isFriendRequestSendAlready {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("REQUEST IS SEND ALREADY")))
	}

	friendRequest := &models.FriendRequest{User1Id: requesterId, User2Id: wantsToBeFriendWith.Id}
	applicant, err := s.UserAccess.AddFriendshipRequest(friendRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}

	return c.JSON(http.StatusCreated, applicant)
}

func (s *userService) RemoveApplicationForFriendship(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	notWanted := c.Param("notWanted")

	application := &models.FriendRequest{User1Id: notWanted, User2Id: requesterId}

	err := s.UserAccess.RemoveApplicationForFriendship(application)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveApplicationForFriendship, err))
	}

	return c.String(http.StatusOK, "")
}

func (s *userService) AcceptFriendship(c echo.Context) error {
	acceptorId := c.Get("userId").(string)

	requesterId := c.Param("wantedUser")

	friendship := &models.Friends{User1Id: requesterId, User2Id: acceptorId}
	areFriends, err := s.UserAccess.AreFriends(friendship)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}
	if areFriends {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("You are friends already")))
	}

	newMember, err := s.UserAccess.ConfirmFriendship(requesterId, acceptorId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}

	return c.JSON(http.StatusCreated, newMember)
}

func (s *userService) RemoveFriend(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	unfriend := c.Param("notWanted")

	friendship := &models.Friends{User1Id: requesterId, User2Id: unfriend}

	err := s.UserAccess.RemoveFriend(friendship)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveFriend, err))
	}

	return c.String(http.StatusOK, "")
}

//
//
//
//	SOCIETY PART
//
//
//

func (s *userService) CreateSociety(c echo.Context) error {
	creatorId := c.Get("userId").(string)

	society := new(models.Society)
	if err := c.Bind(society); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrCreateSociety, err))
	}

	society, err := s.UserAccess.CreateSocietyWithAdmin(society, creatorId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrCreateSociety, err))
	}

	return c.JSON(http.StatusCreated, society)
}

func (s *userService) GetSocietiesWithPaging(c echo.Context) error {
	//can call also unregistered user
	f := c.QueryParam("from")
	t := c.QueryParam("to")

	from, err := strconv.Atoi(f)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	to, err := strconv.Atoi(t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	if to-from < 0 {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, fmt.Errorf("To is smaller than from")))
	}

	societies, allSocieties, err := s.UserAccess.GetSocietiesWithPaging(from, to)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetSociety, err))
	}

	return c.JSON(http.StatusOK, models.SocietyPagingAnsw{
		Societies: s.mapSocietyToSocietyAnswSimple(societies),
		Paging:    models.Paging{From: from, To: to, TotalCount: allSocieties},
	})
}

func (s *userService) GetSociety(c echo.Context) error {
	//can call also unregistered user
	id := c.Param("id")

	society, err := s.UserAccess.GetSociety(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetSociety, err))
	}

	return c.JSON(http.StatusOK, society)
}

func (s *userService) GetSocietiesByIds(c echo.Context) error {
	idsString := c.QueryParam("ids")
	ids := strings.Split(idsString, ",")

	societies, err := s.UserAccess.GetSocieties(ids)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUsers, err))
	}

	return c.JSON(http.StatusOK, societies)
}

func (s *userService) GetMySocieties(c echo.Context) error {
	userId := c.Get("userId").(string)

	societies, err := s.UserAccess.GetUserSocieties(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUserSocieties, err))
	}

	return c.JSON(http.StatusOK, societies)
}

func (s *userService) GetSocietyAdmins(c echo.Context) error {
	societyId := c.Param("societyId")

	admins, err := s.UserAccess.GetSocietyAdmins(societyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetSocietyMembers, err))
	}

	return c.JSON(http.StatusOK, admins)
}

func (s *userService) GetEditableSocieties(c echo.Context) error {
	userId := c.Get("userId").(string)

	societies, err := s.UserAccess.GetEditableSocieties(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetEditableSocieties, err))
	}

	return c.JSON(http.StatusOK, societies)
}

func (s *userService) GetSocietyMembers(c echo.Context) error {
	societyId := c.Param("societyId")

	admins, err := s.UserAccess.GetSocietyMembers(societyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetSocietyMembers, err))
	}

	return c.JSON(http.StatusOK, admins)
}

func (s *userService) UpdateSociety(c echo.Context) error {
	userId := c.Get("userId").(string)

	society := new(models.Society)
	if err := c.Bind(society); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrUpdateSociety, err))
	}

	admin, _, _ := s.UserAccess.IsUserSocietyAdmin(userId, society.Id)
	if !admin {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrUpdateSociety, fmt.Errorf("You are not an admin of society ")))
	}

	society, err := s.UserAccess.UpdateSociety(society)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrUpdateSociety, err))
	}

	return c.JSON(http.StatusOK, society)
}

func (s *userService) GetSocietyRequests(c echo.Context) error {
	societyId := c.Param("societyId")

	requests, err := s.UserAccess.GetSocietyRequests(societyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetSocietyRequests, err))
	}

	return c.JSON(http.StatusOK, requests)
}

func (s *userService) AcceptApplicant(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(models.UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	isAdmin, _, err := s.UserAccess.HasUserSocietyEditorRights(userId, request.SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, err))
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, errors.New("You are not an admin")))
	}

	isMember, err := s.UserAccess.IsMember(request.UserId, request.SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, err))
	}
	if isMember {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, errors.New("You are already a member of society")))
	}

	newMember, err := s.UserAccess.AcceptApplicant(request.UserId, request.SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, err))
	}

	return c.JSON(http.StatusCreated, newMember)
}

func (s *userService) DismissApplicant(c echo.Context) error {
	requester := c.Get("userId").(string)
	societyId := c.Param("societyId")
	removingUserId := c.Param("userId")

	hasRights, _, err := s.UserAccess.HasUserSocietyEditorRights(requester, societyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	if !hasRights {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrDismissApplicant, fmt.Errorf("You are not an admin ")))
	}

	err = s.UserAccess.RemoveApplicationForMembership(&models.Applicant{UserId: removingUserId, SocietyId: societyId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *userService) ChangeMemberRights(c echo.Context) error {
	id := c.Get("userId").(string)

	requests := []models.Member{}
	err := c.Bind(&requests)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}
	if len(requests) < 0 {
		return c.NoContent(http.StatusOK)
	}

	isAdmin, numOfAdmins, err := s.UserAccess.IsUserSocietyAdmin(id, requests[0].SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrChangeMemberRights, err))
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrDismissApplicant, fmt.Errorf("You are not an admin ")))
	}

	if numOfAdmins == 1 && requests[0].UserId == id && requests[0].Permission == models.Membership("member") {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrChangeMemberRights, fmt.Errorf("You are the only one admin in group ")))
	}

	member, err := s.UserAccess.ChangeUserRights(requests, requests[0].SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrChangeMemberRights, err))
	}

	return c.JSON(http.StatusOK, member)
}

func (s *userService) RemoveMember(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	wantsToRemove := c.Param("removingId")
	societyId := c.Param("societyId")

	admin, adminNumber, err := s.UserAccess.IsUserSocietyAdmin(requesterId, societyId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if requesterId == wantsToRemove { //user asks for removing himself
		if admin {
			if adminNumber == 1 {
				return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrRemoveMember, fmt.Errorf("If you want to delete yourself you have to press delete Society ")))
			}

			err = s.UserAccess.RemoveMember(wantsToRemove, societyId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveMember, err))
			}
		} else {
			err = s.UserAccess.RemoveMember(wantsToRemove, societyId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveMember, err))
			}
		}
		return c.NoContent(http.StatusOK)

	} else { //admin removes someone
		if !admin {
			return c.JSON(http.StatusUnauthorized, custom_errors.WrapError(custom_errors.ErrRemoveMember, fmt.Errorf("You are not an admin ")))
		}
		err = s.UserAccess.RemoveMember(wantsToRemove, societyId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveMember, err))
		}
	}

	return c.NoContent(http.StatusOK)
}

func (s *userService) DeleteSociety(c echo.Context) error {
	userId := c.Get("userId").(string)
	societyId := c.Param("societyId")

	admin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, societyId)
	if err != nil {
		return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrUnauthorized, err))
	}

	if !admin {
		return c.JSON(http.StatusUnauthorized, custom_errors.WrapError(custom_errors.ErrUnauthorized, fmt.Errorf("You are not an admin of society ")))
	}

	err = s.UserAccess.DeleteSociety(societyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteSociety, err))
	}

	return c.NoContent(http.StatusOK)
}
