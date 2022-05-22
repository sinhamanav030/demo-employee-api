package customErrors

import (
	"net/http"
)

const (
	ErrorInvalidRequest    = "invalid request"
	ErrorInternalServer    = "internal server error"
	ErrorUserExists        = "user exists"
	ErrorDataNotFound      = "data not found"
	ErrorAuthFailed        = "authorization failed"
	ErrorUnAuthorized      = "unauthorized"
	ErrorValidationRequest = "failed to perform validation"
	ErrorValidation        = "request validation failed"
	ErrorExpiredToken      = "token has expired"
	ErrorInvalidToken      = "token is invalid"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (e ErrorResponse) DataNotFound(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusNotFound
}

func (e ErrorResponse) InvalidRequest(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusBadRequest
}

func (e ErrorResponse) UnAuthorized(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusUnauthorized
}

func (e ErrorResponse) Conflict(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusConflict
}

func (e ErrorResponse) InvalidValidationRequest(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusBadRequest
}

func (e ErrorResponse) ValidationFailed(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusBadRequest
}

func (e ErrorResponse) ExpiredToken(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusBadRequest
}

func (e ErrorResponse) InavlidToken(msg string) (ErrorResponse, int) {
	e.Message = msg
	return e, http.StatusBadRequest
}

func FindErrorType(err string) (ErrorResponse, int) {
	errResp := ErrorResponse{}
	switch err {
	case ErrorDataNotFound:
		return errResp.DataNotFound(err)

	case ErrorInvalidRequest:
		return errResp.InvalidRequest(err)

	case ErrorUnAuthorized:
		return errResp.UnAuthorized(err)

	case ErrorUserExists:
		return errResp.Conflict(err)

	case ErrorValidationRequest:
		return errResp.InvalidValidationRequest(err)

	case ErrorValidation:
		return errResp.ValidationFailed(err)

	case ErrorInvalidToken:
		return errResp.InavlidToken(err)

	case ErrorExpiredToken:
		return errResp.ExpiredToken(err)

	default:
		errResp.Message = err
		return errResp, http.StatusInternalServerError
	}
}
