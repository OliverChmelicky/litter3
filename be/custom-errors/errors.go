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
	ErrBindingRequest         = "ERROR BINDING REQUEST"
	ErrConflict               = "CONFLICT"
	ErrInsufficientPermission = "INSUFFICIENT PERMISSION"

	ErrCreateUser     = "ERROR CREATE USER"
	ErrGetUser        = "ERROR GET USER"
	ErrGetCurrentUser = "ERROR GET CURRENT USER"
	ErrUpdateUser     = "ERR UPDATE USER"
	ErrDeleteUser     = "ERROR DELETE USER"

	ErrGetSocietyMembers = "ERROR GET SOCIETY MEMBERS"

	ErrApplyForMembership             = "ERROR APPLY FOR MEMBERSHIP"
	ErrRemoveApplicationForMembership = "ERROR REMOVE APPLICATION FOR MEMBERSHIP"
	ErrAcceptApplicant                = "ERROR ACCEPT APPLICANT"

	ErrApplyForFriendship             = "ERROR APPLY FOR FRIENDSHIP"
	ErrRemoveApplicationForFriendship = "ERROR REMOVE APPLICATION FOR FRIENDSHIP"
	ErrRemoveFriend                   = "ERROR REMOVE FRIEND"

	ErrCreateComment = "ERROR CREATE COMMENT"
	ErrGetComment    = "ERROR GET COMMENT"
	ErrUpdateComment = "ERROR UPDATE COMMENT"
	ErrDeleteComment = "ERROR DELETE COMMENT"

	ErrCreateTrash = "ERROR CREATE TRASH"
	ErrGetTrash    = "ERROR GET TRASH"
	ErrUpdateTrash = "ERROR UPDATE TRASH"
	ErrDeleteTrash = "ERROR DELETE TRASH"

	ErrCreateCollectionRaw = "CREATE RAW COLLECTION"
	ErrUpdateCollection    = "UPDATE COLLECTION"

	ErrCreateEvent               = "ERROR CREATE EVENT"
	ErrCreateCollectionFromEvent = "CREATE COLLECTION FROM EVENT"
	ErrGetEvent                  = "ERROR GET EVENT"
	ErrGetSocietyEvent           = "ERROR GET SOCIETY EVENTS"
	ErrGetUserEvent              = "ERROR GET USER EVENTS"
	ErrEditEventRights           = "ERROR EDIT EVENT RIGHTS"
	ErrAttendEvent               = "ERROR ATTEND EVENT"
	ErrCannotAttendEvent         = "ERROR CANNOT EVENT"
	ErrUpdateEvent               = "ERROR UPDATE EVENT"
	ErrDeleteEvent               = "ERROR DELETE EVENT"
)
