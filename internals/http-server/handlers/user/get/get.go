package get

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "user-api-service/internals/lib/api/responce"
	"user-api-service/internals/models"
)

type Response struct {
	resp.Response
	User models.User `json:"user" validate:"required"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserGetter
type UserGetter interface {
	GetUser(id string) (models.User, error)
}

func New(log *slog.Logger, storage UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")

		log.Info("get uuid from request", slog.String("id", id))

		if err := validator.New(validator.WithRequiredStructEnabled()).Var(id, "uuid"); err != nil {

			log.Error("invalid uuid", err)

			render.JSON(w, r, resp.Error("invalid uuid"))

			return
		}

		user, err := storage.GetUser(id)

		if err != nil {
			log.Error("failed to get user", err)

			render.JSON(w, r, resp.Error("failed to get user"))
			return
		}

		log.Info("user found")

		responseOK(w, r, user)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user models.User) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		User:     user,
	})
}
