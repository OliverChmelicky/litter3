package user

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"net/http"
)

type userService struct {
	userAccess *userAccess
}

func CreateService(db *pg.DB) *userService {
	access := &userAccess{db: db}
	return &userService{access}
}

func (s *userService) CreateUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.userAccess.CreateUser(user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNotImplemented, user)
}

func (s *userService) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := s.userAccess.GetUser(id)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User with id %s does not exist", id))
	}

	return c.JSON(http.StatusOK, user)
}

func (s *userService) GetCurrentUser(c echo.Context) error {
	id := c.Get("userId").(string)

	user, err := s.userAccess.GetUser(id)
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

	_, err := s.userAccess.GetUser(user.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Trash with provided Id does not exist")
	}

	user, err = s.userAccess.UpdateUser(user)
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

	//chceck ci uz je member
	if s.IsMember(requesterId, request.UserId) {
		return c.String(http.StatusConflict, "User is already a member")
	}

	applicant, err := s.userAccess.AddApplicant(&Applicant{SocietyId: request.SocietyId, UserId: requesterId})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNotImplemented, applicant)
}

func (s *userService) RemoveApplicationForMembership(c echo.Context) error {
	requesterId := c.Get("userId").(string)
	societyId := c.Param("societyId")

	err := s.userAccess.RemoveApplicationForMembership(&Applicant{UserId: requesterId, SocietyId: societyId})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "")
}

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

	society, err := s.userAccess.CreateSocietyWithAdmin(society, creatorId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, society)
}

func (s *userService) GetSociety(c echo.Context) error {
	id := c.Param("id")

	society, err := s.userAccess.GetSociety(id)
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

	_, err := s.userAccess.GetSociety(society.Id)
	if err != nil {
		return c.String(http.StatusNotFound, "Society with provided Id does not exist")
	}

	admin, _, _ := s.isUserSocietyAdmin(userId, society.Id)
	if !admin {
		return c.String(http.StatusForbidden, "You have no right to update society")
	}

	society, err = s.userAccess.UpdateSociety(society)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating society")
	}

	return c.JSON(http.StatusOK, society)
}

func (s *userService) AcceptApplicant(c echo.Context) error {
	//transactional...
	//from applicant table to member
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) ChangeMemberRights(c echo.Context) error {
	//caller is admin
	// access.update user membership
	return c.JSON(http.StatusNotImplemented, "Implement me")
}

func (s *userService) DismissApplicant(c echo.Context) error {
	userId := c.Get("userId").(string)

	request := new(UserGroupRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid arguments")
	}

	admin, _, err := s.isUserSocietyAdmin(userId, request.SocietyId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if !admin {
		return c.String(http.StatusUnauthorized, "You are not an admin")
	}

	err = s.userAccess.RemoveApplicationForMembership(&Applicant{UserId: request.UserId, SocietyId: request.SocietyId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "")
}

func (s *userService) RemoveMember(c echo.Context) error {
	//DELETE {userId, societyId}
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
				err = s.userAccess.DeleteSociety(request.SocietyId)
			} else {
				err = s.userAccess.RemoveMember(request.UserId, request.SocietyId)
			}

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		} else {
			err = s.userAccess.RemoveMember(request.UserId, request.SocietyId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
		return c.String(http.StatusOK, "You were removed")

	} else { //admin removes someone
		if !admin {
			return c.String(http.StatusUnauthorized, "You are not an admin")
		}
		err = s.userAccess.RemoveMember(request.UserId, request.SocietyId)
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

	err = s.userAccess.DeleteSociety(societyId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Implement me")
}

func (s *userService) isUserSocietyAdmin(userId, societyId string) (bool, int, error) {
	admins, err := s.userAccess.GetSocietyAdmins(societyId)
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

func (s *userService) IsMember(userId, societyId string) bool {
	member, err := s.userAccess.IsMember(userId, societyId)
	if err != nil {
		return false
	}

	return member
}
