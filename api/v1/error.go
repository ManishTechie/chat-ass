package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// ErrorConstant is a constant type of known API Errors.
// They are mapped to HTTP Response statuses like 404, 500 etc.
type ErrorConstant string

const (
	// when something went wrong with internal server
	InternalServerError ErrorConstant = "INTERNAL_SERVER_ERROR"

	// when error comes while fetching data from db
	FetchDBDataError ErrorConstant = "FETCH_DB_DATA_ERROR"

	// when error comes while no data found in db
	EmptyDBDataError ErrorConstant = "EMPTY_DB_DATA_ERROR"

	// when error comes while inserting data from db
	InsertDBDataError ErrorConstant = "INSERT_DB_DATA_ERROR"

	// when error comes while deleting data from db
	DeleteDBDataError ErrorConstant = "DELETE_DB_DATA_ERROR"

	// when error comes while updating data from db
	UpdateDBDataError ErrorConstant = "UPDATE_DB_DATA_ERROR"

	// when any common validation error occurs
	ValidationError ErrorConstant = "VALIDATION_ERROR"

	// when error occurs during parsing request-body
	RequestParseError ErrorConstant = "REQUEST_PARSE_ERROR"

	// when an entity is already exist, and you don't want to add one more like that
	ConflictError ErrorConstant = "CONFLICT_ERROR"

	// when an entity is unauthorized to access any resource
	UnAuthorizedError ErrorConstant = "UNAUTHORIZED_ERROR"

	// token is sent but in wrong syntax or format
	TokenFormatERROR ErrorConstant = "TOKEN_FORMAT_ERROR"
	// when range out of bound
	MaxLimitError ErrorConstant = "MAX_LIMIT_REACHED_ERROR"

	// when service is unavailable
	ServiceUnavailableError ErrorConstant = "SERVICE_UNAVAILABLE_ERROR"

	// Duplicate
	Duplicate ErrorConstant = "DUPLICATE"

	// when any common validation error occurs
	BadRequestError ErrorConstant = "BAD_REQUEST_ERROR"

	// Forbidden Error
	ForbiddenError ErrorConstant = "INVALID_KEY_OR_DATA"
)

// customErrors is an internal map of known API Errors, so we can write HTTP Statuse accordingly.
var customErrors = map[ErrorConstant]int{
	InternalServerError:     http.StatusInternalServerError,
	FetchDBDataError:        http.StatusInternalServerError,
	InsertDBDataError:       http.StatusInternalServerError,
	UpdateDBDataError:       http.StatusInternalServerError,
	DeleteDBDataError:       http.StatusInternalServerError,
	EmptyDBDataError:        http.StatusNotFound,
	ValidationError:         http.StatusBadRequest,
	RequestParseError:       http.StatusBadRequest,
	ConflictError:           http.StatusConflict,
	UnAuthorizedError:       http.StatusUnauthorized,
	TokenFormatERROR:        http.StatusUnauthorized,
	MaxLimitError:           http.StatusTooManyRequests,
	ServiceUnavailableError: http.StatusServiceUnavailable,
	Duplicate:               http.StatusBadRequest,
	BadRequestError:         http.StatusBadRequest,
	ForbiddenError:          http.StatusForbidden,
}

// APIError is an internal type of API errors.
type APIError struct {
	Code          int           `json:"code"`
	ErrorConstant ErrorConstant `json:"error_constant"`
	Message       []string      `json:"message"`
}

// APIErrorResponse is the JSON objet in the error response's body.
type APIErrorResponse struct {
	Error APIError `json:"error"`
}

// NewAPIError creates a new api error with a known ErrorConstant and zero or more error messages.
func NewAPIError(errorConst ErrorConstant, message ...string) *APIError {
	return &APIError{
		Code:          customErrors[errorConst],
		ErrorConstant: errorConst,
		Message:       message,
	}
}

// Error returns the API Error as a string, so APIError satisfies the Go error interface.
func (apiErr *APIError) Error() string {
	msg := fmt.Sprintf("API Error Code: %d Constant: %s", apiErr.Code, apiErr.ErrorConstant)
	if apiErr.Message != nil && len(apiErr.Message) > 0 {
		msg += fmt.Sprintf(" Message: %s", strings.Join(apiErr.Message, ", "))
	}
	return msg
}

// Abort returns the APIError's code and the APIError itself, so it can be used in
// gin.Context.JSON or AbortWithJSON methods.
func (apiErr *APIError) Abort() (int, *APIError) {
	return apiErr.Code, apiErr
}

// IsAPIError returns true and the APIError if an error is an APIError,
// and returns false, nil if the error is not an APIError
func IsAPIError(err error) (ok bool, apiErr *APIError) {
	apiErr, ok = err.(*APIError)
	return
}

// TError returns the error of error typed interface
func (apiErr *APIError) TError() error {
	return errors.New(apiErr.Error())
}
