package przepisy

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"przepisyapi/internal/sqlc"
	"przepisyapi/internal/sqlite"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/go-http-helpers/utils"
)

type Options struct {
	Address        string
	DBPath         string
	AllowedOrigins string
}

type Server struct {
	Db *sqlite.Client

	Address        string
	allowedOrigins []string
	Router         chi.Router
	auther         Auther
}

func New(address, allowedOrigins, dbPath, jwtSecret, adminPassword string) (*Server, error) {
	ctx := context.Background()
	// db
	godotenv.Load()
	c, err := sqlite.New(ctx, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite: %w", err)
	}

	// http
	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.ResNotFound(w, "page")
	})

	// init
	server := &Server{
		Address:        address,
		allowedOrigins: strings.Split(allowedOrigins, ","),
		Db:             c,
		Router:         r,
		auther: Auther{
			Db:        c,
			JwtSecret: jwtSecret,
		},
	}

	if users, _ := c.Queries.ListUsers(ctx); len(users) == 0 {
		// Use provided admin password or generate a random one
		password := adminPassword
		if password == "" {
			password = rand.Text()
		}
		params := sqlc.CreateUserParams{
			ID:       utils.CreateUUID(),
			Username: "admin",
			Password: password,
			Email:    "admin@cyber-man.pl",
		}
		slog.Info("no users found; creating admin user", "params", params)
		u, err := c.Queries.CreateUser(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create the default user: %w", err)
		}
		now := time.Now()
		c.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:               u.ID,
			EmailConfirmedAt: &now,
		})
	}

	return server, nil
}

func (s *Server) Serve() {
	chii.SetupMiddlewares(s.Router, s.allowedOrigins)
	// s.SetupRoutes()
	fmt.Printf("listening on address '%s'\n", s.Address)
	chii.PrintRoutes(s.Router)
	err := http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Errorf("listen and serve failed s.Address='%s': %w", s.Address, err))
	}
}
