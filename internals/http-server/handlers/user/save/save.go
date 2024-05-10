package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "user-api-service/internals/lib/api/responce"
	"user-api-service/internals/models"
)

type Request struct {
	FirstName string `json:"first-name" validate:"required,alpha"`
	LastName  string `json:"last-name" validate:"required,alpha"`
	Email     string `json:"e-mail" validate:"required,email"`
	Age       uint   `json:"age" validate:"required,gte=0,lte=150"`
}

type Response struct {
	resp.Response
	Id string `json:"id" validate:"required, uuid"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type UserSaver interface {
	SaveUser(user models.User) (string, error)
}

func New(log *slog.Logger, storage UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("request body is empty"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {

			log.Error("invalid request", err)

			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))

			return
		}

		user := models.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Age:       req.Age,
		}
		id, err := storage.SaveUser(user)

		if err != nil {
			log.Error("failed to save user", err)

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		log.Info("user saved", slog.String("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
