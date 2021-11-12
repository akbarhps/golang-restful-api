package exception

import (
	"go-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func badRequest(res *model.WebResponse, err string) {
	res.Code = http.StatusBadRequest
	res.Status = "Bad Request"
	res.Error = err
}

func notFound(res *model.WebResponse, err string) {
	res.Code = http.StatusNotFound
	res.Status = "Record Not Found"
	res.Error = err
}

func internalServerError(res *model.WebResponse, err string) {
	res.Code = http.StatusInternalServerError
	res.Status = "Internal Server Error"
	res.Error = err
}

func unauthorizedError(res *model.WebResponse, err string) {
	res.Code = http.StatusUnauthorized
	res.Status = "Unauthorized Access"
	res.Error = err
}

func ErrorHandler(c *gin.Context, err interface{}) {
	res := &model.WebResponse{}

	switch err.(type) {
	case RecordDuplicateError:
		badRequest(res, err.(RecordDuplicateError).Error())
	case RecordNotFoundError:
		notFound(res, err.(RecordNotFoundError).Error())
	case validator.ValidationErrors:
		badRequest(res, err.(validator.ValidationErrors).Error())
	case InvalidCredentialError:
		badRequest(res, err.(InvalidCredentialError).Error())
	case InvalidSignatureError:
		unauthorizedError(res, err.(InvalidSignatureError).Error())
	case error:
		internalServerError(res, err.(error).Error())
	}

	c.IndentedJSON(res.Code, res)
}
