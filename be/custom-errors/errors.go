package custom_errors

import "errors"

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
	ErrPageConfigNotFound  = errors.New("pageConfig not found")
	ErrPageDataNotFound    = errors.New("no pageData found")
	ErrInterfaceConversion = errors.New("interface conversion failed")
	ErrProductNotActive    = errors.New("product is not valid")
	ErrIndexOutOfRange     = errors.New("index out of range")
	ErrSchemaNotFound      = errors.New("schema not found")
	ErrDetailQueryNotValid = errors.New("query is not valid")
	ErrAmendNotSet         = errors.New("amend is not set")
	ErrProductNotFound     = errors.New("no product found")
	ErrMinAggNotFound      = errors.New("min aggregation not found")
	ErrMaxAggNotFound      = errors.New("max aggregation not found")
	ErrMissingToken        = errors.New("token is missing")
	ErrInvalidFilter       = errors.New("invalid filter query value")
)

var (
	ErrBindingRequest = "ERROR BINDING REQUEST"
)
