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
	"user-api-service/internals/models"
)

type Request struct {
	FirstName string `json:"first-name" validate:"required,alpha"`
	LastName  string `json:"last-name" validate:"required,alpha"`
	Email     string `json:"e-mail" validate:"required,email"`
	Age       uint   `json:"age" validate:"required,gte=0,lte=150"`
}

type Response struct {
	Status uint   `json:"status"`
	Error  string `json:"error,omitempty"`
}

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
			log.Error("failed to decode request body", err)

			render.JSON(w, r, "failed to decode request")

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {

			log.Error("invalid request", err)

			for _, err := range err.(validator.ValidationErrors) {
				log.Error("Validation error: Field '%s', Tag '%s'", err.Field(), err.Tag())
			}

			render.JSON(w, r, "invalid request")

			return
		}

		id := chi.URLParam(r, "id")

		//TODO make it clenear
		user := models.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Age:       req.Age,
		}
		err = storage.UpdateUser(user, id)

		if err != nil {
			log.Error("failed to update user", err)

			render.JSON(w, r, "failed to update user")

			return
		}

		log.Info("user updated", slog.String("id", id))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: http.StatusOK,
	})
}