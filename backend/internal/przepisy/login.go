package przepisy

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tomek7667/go-http-helpers/h"
)

type GetLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) PostLogin() {
	s.Router.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
		dto, err := h.GetDto[GetLoginDto](r)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}
		user, err := s.Db.Queries.GetUserByUsername(r.Context(), dto.Username)
		if err != nil || user.Password != dto.Password {
			slog.Warn(
				"s.Dber.GetUserByUsername(r.Context(), dto.Username)",
				"err", err,
				"given pass", dto.Password,
			)
			h.ResErr(w, fmt.Errorf("invalid username or password"))
			return
		}
		token, err := s.auther.GetToken(&user)
		if err != nil {
			h.ResErr(w, err)
			return
		}
		h.ResSuccess(w, map[string]string{
			"token": token,
		})
	})
}
