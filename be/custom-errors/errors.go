package custom_errors

type ErrorModel struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"errorMessage"`
}

func WrapError(errorType string, err error) ErrorModel {
	return ErrorModel{
		ErrorType: errorType,
		Message:   err.Error(),
	}
}

var (
	ErrBindingRequest = "ERROR BINDING REQUEST"

	ErrGetUser        = "ERROR GET USER"
	ErrGetCurrentUser = "ERROR GET CURRENT USER"

	ErrApplyForMembership             = "ERROR APPLY FOR MEMBERSHIP"
	ErrRemoveApplicationForMembership = "ERROR REMOVE APPLICATION FOR MEMBERSHIP"
	ErrAcceptApplicant                = "ERROR ACCEPT APPLICANT"

	ErrApplyForFriendship             = "ERROR APPLY FOR FRIENDSHIP"
	ErrRemoveApplicationForFriendship = "ERROR REMOVE APPLICATION FOR FRIENDSHIP"

	ErrRemoveFriend = "ERROR REMOVE FRIEND"

	ErrCreateComment = "ERROR CREATE COMMENT"
	ErrGetComment    = "ERROR GET COMMENT"
	ErrUpdateComment = "ERROR UPDATE COMMENT"
	ErrDeleteComment = "ERROR DELETE COMMENT"

	ErrCreateTrash = "ERROR CREATE TRASH"
	ErrGetTrash    = "ERROR GET TRASH"
	ErrUpdateTrash = "ERROR UPDATE TRASH"
	ErrDeleteTrash = "ERROR DELETE TRASH"

	ErrConflict = "CONFLICT"
)
