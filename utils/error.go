package utils

import (
	"cryptoshare/dto"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%v field is required.", field)
	case "oneof":
		return fmt.Sprintf("%v field must be one of %v", field, fe.Param())
	case "email":
		return "Invalid email."
	case "gte", "lte":
		return "invalid length"
	default:
		return "invalid payload" // default error
	}
}

func GenerateValidationErrorMessage(err error) string {
	if vErr, ok := err.(validator.ValidationErrors); ok {
		errMsg := ""
		for _, fieldErr := range vErr {
			errMsg += msgForTag(fieldErr)
		}
		return errMsg
	}
	return err.Error()
}
func GenerateGormErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrMsg = err.Error()

	if IsErrNotFound(err) || IsDuplicate(err) {
		res.ErrCode = 400
		res.HttpStatusCode = http.StatusBadRequest
		return res
	}
	res.ErrCode = 500
	res.HttpStatusCode = http.StatusInternalServerError
	return res
}

func GenerateValidationErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrMsg = err.Error()
	if IsValidationError(err) {
		res.ErrCode = 422
		res.ErrMsg = GenerateValidationErrorMessage(err)
		res.HttpStatusCode = http.StatusUnprocessableEntity
		return res
	}
	res.ErrCode = 500
	res.HttpStatusCode = http.StatusInternalServerError
	return res
}

func GenerateAuthErrorResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 401
	res.ErrMsg = "pemission denied"
	res.HttpStatusCode = http.StatusUnauthorized
	return res
}

func GenerateBadRequestResponse(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 400
	res.ErrMsg = "invalid request"
	res.HttpStatusCode = http.StatusBadRequest
	return res
}

func GenerateServerError(err error) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 500
	res.HttpStatusCode = http.StatusInternalServerError
	return res
}

func GenerateSuccessResponse(data any) *dto.Response {
	res := &dto.Response{}
	res.ErrCode = 0
	res.ErrMsg = "Success"
	res.Data = data
	res.HttpStatusCode = http.StatusOK
	return res
}
