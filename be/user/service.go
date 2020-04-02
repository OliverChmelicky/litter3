package user

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
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
		return c.String(http.StatusNotFound, fmt.Sprintf("User with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	id := c.Get("userId").(string)

	user, err := s.UserAccess.GetUser(id)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User with id %s does not exist", id))
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

	request := new(UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	isMember, err := s.UserAccess.IsMember(request.UserId, request.SocietyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if isMember {
		return c.String(http.StatusConflict, "User is already a member")
	}

	applicant, err := s.UserAccess.AddApplicant(&Applicant{SocietyId: request.SocietyId, UserId: requesterId})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNotImplemented, applicant)
}

func (s *userService) RemoveApplicationForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	societyId := c.Param("societyId")

	err := s.UserAccess.RemoveApplicationForMembership(&Applicant{UserId: requesterId, SocietyId: societyId})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "")
}

//len nazvy funkcii su ok, doimplementuj
//func (s *userService) ApplyForFriendship(c echo.Context) error {
//	requesterId := c.Get("userId").(string)
//
//	request := new(UserGroupRequest)
//	if err := c.Bind(request); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid arguments")
//	}
//
//	isMember, err := s.UserAccess.IsMember(request.UserId, request.SocietyId)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, err.Error())
//	}
//	if isMember {
//		return c.String(http.StatusConflict, "User is already a member")
//	}
//
//	applicant, err := s.UserAccess.AddApplicant(&Applicant{SocietyId: request.SocietyId, UserId: requesterId})
//	if err != nil {
//		return c.String(http.StatusInternalServerError, err.Error())
//	}
//
//	return c.JSON(http.StatusNotImplemented, applicant)
//}
//
//func (s *userService) RemoveApplicationForFriendship(c echo.Context) error {
//	requesterId := c.Get("userId").(string)
//	societyId := c.Param("societyId")
//
//	err := s.UserAccess.RemoveApplicationForMembership(&Applicant{UserId: requesterId, SocietyId: societyId})
//	if err != nil {
//		return c.String(http.StatusInternalServerError, err.Error())
//	}
//
//	return c.String(http.StatusOK, "")
//}

//func (s *userService) AcceptFriendship(c echo.Context) error {
//	userId := c.Get("userId").(string)
//
//	newMemberId := c.Param("userId")
//	societyId := c.Param("societyId")
//
//	isMember, err := s.UserAccess.IsMember(userId, societyId)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, err.Error())
//	}
//	if isMember {
//		return c.String(http.StatusConflict, "You are already a member of society")
//	}
//
//	newMember, err := s.UserAccess.AcceptApplicant(newMemberId, societyId)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, err.Error())
//	}
//
//	return c.JSON(http.StatusCreated, newMember)
//}
//
//func (s *userService) DismissFriendship(c echo.Context) error {
//	admin := c.Get("userId").(string)
//	societyId := c.Param("societyId")
//	removingUserId := c.Param("userId")
//
//	isAdmin, _, err := s.isUserSocietyAdmin(admin, societyId)
//	if err != nil {
//		return c.String(http.StatusNotFound, err.Error())
//	}
//
//	if !isAdmin {
//		return c.String(http.StatusUnauthorized, "You are not an admin")
//	}
//
//	err = s.UserAccess.RemoveApplicationForMembership(&Applicant{UserId: removingUserId, SocietyId: societyId})
//	if err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	}
//
//	return c.String(http.StatusOK, "")
//}

func (s *userService) DeleteUser(c echo.Context) error {
	//remove comments
	//what with events
	//what with collected trash
	//remove notifications
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

//
//
//
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

	admin, _, _ := s.isUserSocietyAdmin(userId, society.Id)
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

	newMemberId := c.Param("userId")
	societyId := c.Param("societyId")

	isMember, err := s.UserAccess.IsMember(userId, societyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if isMember {
		return c.String(http.StatusConflict, "You are already a member of society")
	}

	newMember, err := s.UserAccess.AcceptApplicant(newMemberId, societyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newMember)
}

func (s *userService) DismissApplicant(c echo.Context) error {
	admin := c.Get("userId").(string)
	societyId := c.Param("societyId")
	removingUserId := c.Param("userId")

	isAdmin, _, err := s.isUserSocietyAdmin(admin, societyId)
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

	isAdmin, numOfAdmins, err := s.isUserSocietyAdmin(id, request.SocietyId)
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
	//DELETE {userId, societyId} ----> prerob, tak by v body nebolo nic
	id := c.Get("userId")
	requesterId := fmt.Sprintf("%v", id)

	request := new(UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}
	admin, adminNumber, err := s.isUserSocietyAdmin(requesterId, request.SocietyId)
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
	//remove members
	id := c.Get("userId")
	userId := fmt.Sprintf("%v", id)
	societyId := c.Param("id")

	admin, _, err := s.isUserSocietyAdmin(userId, societyId)
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

func (s *userService) isUserSocietyAdmin(userId, societyId string) (bool, int, error) {
	admins, err := s.UserAccess.GetSocietyAdmins(societyId)
	if err != nil {
		return false, 0, err
	}
	if len(admins) == 0 {
		return false, 0, err
	}

	for _, adminId := range admins {
		if adminId == userId {
			return true, len(admins), nil
		}
	}

	return false, len(admins), nil
}
