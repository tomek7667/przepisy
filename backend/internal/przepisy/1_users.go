package przepisy

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"przepisyapi/internal/crypto"
	"przepisyapi/internal/mails"
	"przepisyapi/internal/sqlc"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/go-http-helpers/utils"
)

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UpdateUserDto struct {
	Password string `json:"password"`
}

func (s *Server) AddUsersRoutes() {
	s.Router.With(h.WithRateLimit(s.loginRateLimiter)).Post("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
		dto, err := h.GetDto[CreateUserDto](r)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}
		randomCode := crypto.RandCode(6)
		newuser, err := s.Db.Queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			ID:               utils.CreateUUID(),
			Username:         dto.Username,
			Password:         dto.Password,
			Email:            dto.Email,
			EmailConfirmCode: &randomCode,
		})
		if err != nil {
			h.ResErr(w, err)
			return
		}
		go s.mailer.SendMail(s.agentID, mails.Options{
			From:    "Przepisy",
			To:      dto.Email,
			Subject: "[Przepisy] Potwierdź e-mail",
			HTML: fmt.Sprintf(
				`Oto kod do potwierdzenia maila: %s`+
					`<br />`+
					`<br />`+
					`Możesz też kliknąć <a href="%s/api/users/%s/confirm?code=%s">tutaj</a> aby automatycznie potwierdzić mail.`,
				randomCode,
				s.FrontendUrl,
				newuser.ID,
				randomCode),
		})
		h.ResSuccess(w, newuser)
	})

	s.Router.With(h.WithRateLimit(s.confirmCodeRateLimiter)).Get("/api/users/{id}/confirm", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		userId := chi.URLParam(r, "id")
		slog.Info("confirm code action", "code", code, "user id", userId)
		user, err := s.Db.Queries.GetUserByID(r.Context(), userId)
		if err != nil {
			h.ResNotFound(w, "user")
			return
		}
		if user.EmailConfirmedAt != nil || user.EmailConfirmCode == nil {
			h.ResErr(w, fmt.Errorf("email already confirmed"))
			return
		}
		if *user.EmailConfirmCode != code {
			h.ResErr(w, fmt.Errorf("invalid confirm code"))
			return
		}
		now := time.Now()
		err = s.Db.Queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
			ID:               user.ID,
			EmailConfirmedAt: &now,
			EmailConfirmCode: nil,
		})
		if err != nil {
			h.ResErr(w, err)
			return
		}
		updatedUser, err := s.Db.Queries.GetUserByID(r.Context(), user.ID)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		token, err := s.auther.GetToken(&updatedUser)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, map[string]string{
			"token": token,
		})
	})

	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/users", func(r chi.Router) {
		// TODO: confirmed account / admins / permissions
		// r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// 	user := chii.GetUser[sqlc.User](r)
		// 	id := chi.URLParam(r, "id")
		// 	dto, err := h.GetDto[UpdateUserDto](r)
		// 	if err != nil {
		// 		h.ResBadRequest(w, err)
		// 		return
		// 	}
		// 	toBeUpdated, err := s.Db.Queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
		// 		ID: id,
		// 	})
		// 	if err != nil {
		// 		h.ResErr(w, err)
		// 		return
		// 	}
		// 	h.ResSuccess(w, toBeUpdated)
		// })

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			_, err := s.Db.Queries.GetUserByID(r.Context(), id)
			if err != nil || user.ID != id || user.Email != "admin@cyber-man.pl" {
				h.ResNotFound(w, "user")
				return
			}

			err = s.Db.Queries.DeleteUser(r.Context(), id)
			if err != nil {
				h.ResErr(w, err)
				return
			}
			h.ResSuccess(w, nil)
		})
	})
}
