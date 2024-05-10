package update

import (
	"errors"
	"github.com/go-chi/chi/v5"
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
	Error string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserUpdater
type UserUpdater interface {
	UpdateUser(user models.User, id string) error
}

func New(log *slog.Logger, storage UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, "empty request")

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Any("err", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {

			log.Error("invalid request")

			render.JSON(w, r, resp.ValidationError(err.(validator.ValidationErrors)))

			return
		}

		id := chi.URLParam(r, "id")

		if err := validator.New(validator.WithRequiredStructEnabled()).Var(id, "uuid"); err != nil {

			log.Error("id validation error ", slog.Any("err", err.Error()))

			render.JSON(w, r, resp.Error("invalid uuid"))

			return
		}

		//TODO make it cleaner
		user := models.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Age:       req.Age,
		}
		err = storage.UpdateUser(user, id)

		if err != nil {
			log.Error("failed to update user", slog.Any("err", err.Error()))

			render.JSON(w, r, resp.Error("failed to update user"))

			return
		}

		log.Info("user updated", slog.String("id", id))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
