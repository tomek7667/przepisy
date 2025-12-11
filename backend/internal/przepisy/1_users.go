package przepisy

import (
	"net/http"

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
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/users", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			users, err := s.Db.Queries.ListUsers(r.Context())
			if err != nil {
				h.ResErr(w, err)
				return
			}
			h.ResSuccess(w, users)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			fetchedUser, err := s.Db.Queries.GetUserByID(r.Context(), id)
			if err != nil {
				h.ResNotFound(w, "user")
				return
			}
			h.ResSuccess(w, fetchedUser)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			dto, err := h.GetDto[CreateUserDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			newuser, err := s.Db.Queries.CreateUser(r.Context(), sqlc.CreateUserParams{
				ID:       utils.CreateUUID(),
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
			})
			if err != nil {
				h.ResErr(w, err)
				return
			}
			h.ResSuccess(w, newuser)
		})

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
