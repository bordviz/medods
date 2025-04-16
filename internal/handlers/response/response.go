package response

import (
	"medods/internal/lib/customerror"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrBadRequest = customerror.NewCustomError("bad request", 400)
)

type Response struct {
	Detail string `json:"detail" example:"response detail"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err customerror.CustomError) {
	render.Status(r, err.Code())
	render.JSON(w, r, Response{Detail: err.Error()})
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
