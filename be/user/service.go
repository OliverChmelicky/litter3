package user

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"net/http"
)

type userService struct {
	UserAccess *UserAccess
}

func CreateService(db *pg.DB) *userService {
	access := &UserAccess{Db: db}
	return &userService{access}
}

func (s *userService) CreateUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.UserAccess.CreateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNotImplemented, user)
}

func (s *userService) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := s.UserAccess.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetUser, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	id := c.Get("userId").(string)

	user, err := s.UserAccess.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, custom_errors.WrapError(custom_errors.ErrGetCurrentUser, err))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) UpdateUser(c echo.Context) error {
	updatorId := c.Get("userId").(string)

	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}
	user.Id = updatorId

	_, err := s.UserAccess.GetUser(user.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	user, err = s.UserAccess.UpdateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating user")
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) ApplyForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(IdMessage)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	isMember, err := s.UserAccess.IsMember(requesterId, request.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}
	if isMember {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("User is already a member")))
	}

	applicant, err := s.UserAccess.AddApplicant(&Applicant{SocietyId: request.Id, UserId: requesterId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}

	return c.JSON(http.StatusOK, applicant)
}

func (s *userService) RemoveApplicationForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	societyId := c.Param("societyId")

	err := s.UserAccess.RemoveApplicationForMembership(&Applicant{UserId: requesterId, SocietyId: societyId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveApplicationForMembership, err))
	}

	return c.String(http.StatusOK, "")
}

//func (s *userService) DeleteUser(c echo.Context) error {
//	//TODO check ci neorganizuje event a potom vymaz event
//	//TODO check ci nie je jediny admin skupiny a potom vymaz
//}

//
//
//
//	FRIENDS PART
//
//

func (s *userService) ApplyForFriendship(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(IdMessage)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	isRequestSend, err := s.UserAccess.AreFriends(&Friends{User1Id: requesterId, User2Id: request.Id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}
	if isRequestSend {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("YOU ARE FIENDS ALREADY")))
	}

	friendRequest := &FriendRequest{User1Id: requesterId, User2Id: request.Id}
	applicant, err := s.UserAccess.AddFriendshipRequest(friendRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForFriendship, err))
	}

	return c.JSON(http.StatusNotImplemented, applicant)
}

func (s *userService) RemoveApplicationForFriendship(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	notWanted := c.Param("unfriendId")

	application := &FriendRequest{User1Id: notWanted, User2Id: requesterId}

	err := s.UserAccess.RemoveApplicationForFriendship(application)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrRemoveApplicationForFriendship, err))
	}

	return c.String(http.StatusOK, "")
}

func (s *userService) AcceptFriendship(c echo.Context) error {
	requesterId := c.Get("userId").(string)

	request := new(FriendRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if requesterId != request.User1Id && requesterId != request.User2Id {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("You are not in this relation")))
	}

	friendship := &Friends{User1Id: request.User1Id, User2Id: request.User2Id}
	areFriends, err := s.UserAccess.AreFriends(friendship)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}
	if areFriends {
		return c.JSON(http.StatusConflict, custom_errors.WrapError(custom_errors.ErrConflict, errors.New("You are friends already")))
	}

	newMember, err := s.UserAccess.ConfirmFriendship(request.User1Id, request.User2Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrApplyForMembership, err))
	}

	return c.JSON(http.StatusCreated, newMember)
}

func (s *userService) RemoveFriend(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	unfriend := c.Param("unfriendId")

	friendship := &Friends{User1Id: requesterId, User2Id: unfriend}

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

func (s *userService) CreateSociety(c echo.Context) error {
	creatorId := c.Get("userId").(string)

	society := new(Society)
	if err := c.Bind(society); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	society, err := s.UserAccess.CreateSocietyWithAdmin(society, creatorId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, society)
}

func (s *userService) GetSociety(c echo.Context) error {
	id := c.Param("id")

	society, err := s.UserAccess.GetSociety(id)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Society with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, society)
}

//TODO
func (s *userService) GetSocietyMembers(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "IMPLEMENT ME")
}

func (s *userService) UpdateSociety(c echo.Context) error {
	userId := c.Get("userId").(string)

	society := new(Society)
	if err := c.Bind(society); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	_, err := s.UserAccess.GetSociety(society.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Society with provided Id does not exist")
	}

	admin, _, _ := s.UserAccess.IsUserSocietyAdmin(userId, society.Id)
	if !admin {
		return c.String(http.StatusForbidden, "You have no right to update society")
	}

	society, err = s.UserAccess.UpdateSociety(society)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating society")
	}

	return c.JSON(http.StatusOK, society)
}

func (s *userService) AcceptApplicant(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, request.SocietyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, err))
	}
	if !isAdmin {
		return c.JSON(http.StatusUnauthorized, custom_errors.WrapError(custom_errors.ErrAcceptApplicant, errors.New("You are not an admin")))
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
	admin := c.Get("userId").(string)
	societyId := c.Param("societyId")
	removingUserId := c.Param("userId")

	isAdmin, _, err := s.UserAccess.IsUserSocietyAdmin(admin, societyId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if !isAdmin {
		return c.String(http.StatusUnauthorized, "You are not an admin")
	}

	err = s.UserAccess.RemoveApplicationForMembership(&Applicant{UserId: removingUserId, SocietyId: societyId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "")
}

func (s *userService) ChangeMemberRights(c echo.Context) error {
	id := c.Get("userId").(string)

	request := new(Member)
	err := c.Bind(request)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	isAdmin, numOfAdmins, err := s.UserAccess.IsUserSocietyAdmin(id, request.SocietyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !isAdmin {
		return c.String(http.StatusForbidden, "You are not an admin")
	}

	if numOfAdmins == 1 && request.UserId == id && request.Permission == membership("member") {
		return c.String(http.StatusConflict, "You are the only one admin in group")
	}

	member, err := s.UserAccess.ChangeUserRights(request)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNotImplemented, member)
}

func (s *userService) RemoveMember(c echo.Context) error {
	//TODO DELETE {userId, societyId} ----> prerob, tak by v body nebolo nic
	id := c.Get("userId")
	requesterId := fmt.Sprintf("%v", id)

	request := new(UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}
	admin, adminNumber, err := s.UserAccess.IsUserSocietyAdmin(requesterId, request.SocietyId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if request.UserId == c.Get("useId") { //user asks for removing himself
		if admin {
			if adminNumber == 1 {
				err = s.UserAccess.DeleteSociety(request.SocietyId)
			} else {
				err = s.UserAccess.RemoveMember(request.UserId, request.SocietyId)
			}

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		} else {
			err = s.UserAccess.RemoveMember(request.UserId, request.SocietyId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
		return c.String(http.StatusOK, "You were removed")

	} else { //admin removes someone
		if !admin {
			return c.String(http.StatusUnauthorized, "You are not an admin")
		}
		err = s.UserAccess.RemoveMember(request.UserId, request.SocietyId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

	}

	return c.String(http.StatusOK, "Success removing user")
}

func (s *userService) DeleteSociety(c echo.Context) error {
	//TODO check pred tym ci nahodou neorganizuje event a vymaz event potom
	id := c.Get("userId")
	userId := fmt.Sprintf("%v", id)
	societyId := c.Param("id")

	admin, _, err := s.UserAccess.IsUserSocietyAdmin(userId, societyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !admin {
		return c.String(http.StatusUnauthorized, "You are not an admin")
	}

	err = s.UserAccess.DeleteSociety(societyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Implement me")
}
