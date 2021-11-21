package middleware

import (
	"github.com/gin-gonic/gin"
	"go-api/exception"
	"go-api/model"
	"net/http"

	"github.com/go-playground/validator"
)

func parseError(err []error) []map[string]string {
	var errors []map[string]string
	for _, e := range err {
		switch e.(type) {
		case exception.FieldError:
			fe := e.(exception.FieldError)
			errors = append(errors, map[string]string{
				"field": fe.Field,
				"error": fe.Message,
			})
		default:
			errors = append(errors, map[string]string{
				"error": e.Error(),
			})
		}
	}
	return errors
}

func badRequest(res *model.WebResponse, err []error) {
	res.Code = http.StatusBadRequest
	res.Status = "Bad Request"
	res.Errors = parseError(err)
}

func notFound(res *model.WebResponse, err exception.NotFoundError) {
	res.Code = http.StatusNotFound
	res.Status = "Record Not Found"
	res.Errors = parseError([]error{err})
}

func internalServerError(res *model.WebResponse, err []error) {
	res.Code = http.StatusInternalServerError
	res.Status = "Internal Server Error"
	res.Errors = parseError(err)
}

func unauthorizedError(res *model.WebResponse, err []error) {
	res.Code = http.StatusUnauthorized
	res.Status = "Unauthorized Access"
	res.Errors = parseError(err)
}

func validationError(res *model.WebResponse, err validator.ValidationErrors) {
	res.Code = http.StatusBadRequest
	res.Status = "Bad Request"
	for _, e := range err {
		res.Errors = append(res.Errors, map[string]string{
			"field": e.Field(),
			"error": e.Tag(),
		})
	}
}

func PanicHandler(c *gin.Context, err interface{}) {
	res := &model.WebResponse{}

	switch err.(type) {
	case exception.Errors:
		badRequest(res, err.(exception.Errors).Errors)
	case validator.ValidationErrors:
		validationError(res, err.(validator.ValidationErrors))
	case exception.TokenError:
		unauthorizedError(res, []error{err.(exception.TokenError)})
	case exception.WrongPasswordError:
		badRequest(res, []error{err.(exception.WrongPasswordError)})
	case exception.DuplicateError:
		badRequest(res, []error{err.(exception.DuplicateError)})
	case exception.NotFoundError:
		notFound(res, err.(exception.NotFoundError))
	case exception.DatabaseError:
		internalServerError(res, []error{err.(exception.DatabaseError)})
	case error:
		internalServerError(res, []error{err.(error)})
	}

	c.AbortWithStatusJSON(res.Code, res)
}
