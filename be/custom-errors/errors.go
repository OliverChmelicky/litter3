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
	ErrBindingRequest                 = "ERROR BINDING REQUEST"
	ErrApplyForMembership             = "ERROR APPLY FOR MEMBERSHIP"
	ErrRemoveApplicationForMembership = "ERROR REMOVE APPLICATION FOR MEMBERSHIP"
	ErrApplyForFriendship             = "ERROR APPLY FOR FRIENDSHIP"
	ErrRemoveApplicationForFriendship = "ERROR REMOVE APPLICATION FOR FRIENDSHIP"
	ErrRemoveFriend                   = "ERROR REMOVE FRIEND"
	ErrConflict                       = "CONFLICT"
)
