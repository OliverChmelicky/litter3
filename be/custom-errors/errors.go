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
	ErrNoToken                = "NO AUTHORIZATION TOKEN"
	ErrUnauthorized           = "UNAUTHORIZED"

	ErrCreateUser     = "ERROR CREATE USER"
	ErrGetUserById    = "ERROR GET USER BY ID"
	ErrGetUsers       = "ERROR GET BATCH USERS BY ID"
	ErrGetUserByEmail = "ERROR GET USER BY EMAIL"
	ErrGetCurrentUser = "ERROR GET CURRENT USER"
	ErrUpdateUser     = "ERR UPDATE USER"
	ErrDeleteUser     = "ERROR DELETE USER"

	ErrChangeMemberRights   = "ERROR CHANGE MEMBER RIGHTS"
	ErrGetSocietyMembers    = "ERROR GET SOCIETY MEMBERS"
	ErrGetEditableSocieties = "ERROR GET EDITABLE SOCIETIES"
	ErrGetSocietyRequests   = "ERROR GET SOCIETY REQUESTS"
	ErrRemoveMember         = "ERROR REMOVE MEMBER"

	ErrApplyForMembership             = "ERROR APPLY FOR MEMBERSHIP"
	ErrRemoveApplicationForMembership = "ERROR REMOVE APPLICATION FOR MEMBERSHIP"
	ErrAcceptApplicant                = "ERROR ACCEPT APPLICANT"
	ErrDismissApplicant               = "ERROR DISMISS APPLICANT"
	ErrGetUserFriends                 = "ERROR GET USER FRIENDS"
	ErrGetMyReqForFriendship          = "ERROR GET MY REQUESTS FOR FRIENDSHIP"

	ErrApplyForFriendship             = "ERROR APPLY FOR FRIENDSHIP"
	ErrRemoveApplicationForFriendship = "ERROR REMOVE APPLICATION FOR FRIENDSHIP"
	ErrRemoveFriend                   = "ERROR REMOVE FRIEND"

	ErrCreateSociety    = "ERROR CREATE SOCIETY"
	ErrGetSociety       = "ERROR GET SOCIETY"
	ErrGetUserSocieties = "ERROR GET USER SOCIETIES"
	ErrUpdateSociety    = "ERROR UPDATE SOCIETY"
	ErrDeleteSociety    = "ERROR DELETE SOCIETY"

	ErrCreateComment = "ERROR CREATE COMMENT"
	ErrGetComment    = "ERROR GET COMMENT"
	ErrUpdateComment = "ERROR UPDATE COMMENT"
	ErrDeleteComment = "ERROR DELETE COMMENT"

	ErrCreateTrash      = "ERROR CREATE TRASH"
	ErrGetTrashById     = "ERROR GET TRASH BY ID"
	ErrGetTrashInRange  = "ERROR GET TRASH IN RANGE"
	ErrGetTrashComments = "ERROR GET TRASH COMMENTS"
	ErrUpdateTrash      = "ERROR UPDATE TRASH"
	ErrDeleteTrash      = "ERROR DELETE TRASH"

	ErrCreateCollectionRaw = "CREATE RAW COLLECTION"
	ErrGetCollectionRaw = "GET RAW COLLECTION"
	ErrUpdateCollection    = "UPDATE COLLECTION"
	ErrAddPickerToCollection = "ERROR ADD PICKER TO COLLECTION"

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

	ErrUploadImage = "ERROR UPLOAD IMAGE"
	ErrLoadImage   = "ERROR LOAD IMAGE"
	ErrDeleteImage = "ERROR DELETE IMAGE"
	ErrDeleteCollectionImage = "ERROR DELETE COLLECTION IMAGE"
)
