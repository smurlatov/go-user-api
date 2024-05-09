package get

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"user-api-service/internals/models"
)

type Response struct {
	Status uint        `json:"status"`
	Error  string      `json:"error,omitempty"`
	User   models.User `json:"user" validate:"required"`
}

type UserGetter interface {
	GetUser(id string) (*models.User, error)
}

func New(log *slog.Logger, storage UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")

		log.Info("get uuid from reques", slog.String("id", id))

		if err := validator.New(validator.WithRequiredStructEnabled()).Var(id, "uuid"); err != nil {

			log.Error("invalid request", err)

			for _, err := range err.(validator.ValidationErrors) {
				log.Error("Validation error: Field '%s', Tag '%s'", err.Field(), err.Tag())
			}

			render.JSON(w, r, "invalid request")

			return
		}

		user, err := storage.GetUser(id)

		if err != nil {
			log.Error("failed to get user", err)

			render.JSON(w, r, "failed to get user")

			return
		}

		log.Info("user found")

		responseOK(w, r, user)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user *models.User) {
	render.JSON(w, r, Response{
		Status: http.StatusOK,
		User:   *user,
	})
}
