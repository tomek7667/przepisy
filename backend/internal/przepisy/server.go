package przepisy

import (
	"context"
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
	"github.com/tomek7667/secrets/secretssdk"
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
	secreter       *secretssdk.Client
}

func New(address, allowedOrigins, dbPath string, secretsClient *secretssdk.Client) (*Server, error) {
	ctx := context.Background()
	jwtSecret, err := secretsClient.GetSecret("przepisy/jwt-token")
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets jwt token: %w", err)
	}
	adminPassword, err := secretsClient.GetSecret("przepisy/admin-password")
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets admin password: %w", err)
	}
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
		secreter:       secretsClient,
		auther: Auther{
			Db:        c,
			JwtSecret: jwtSecret.Value,
		},
	}

	if users, _ := c.Queries.ListUsers(ctx); len(users) == 0 {
		// Use provided admin password or generate a random one
		params := sqlc.CreateUserParams{
			ID:       utils.CreateUUID(),
			Username: "admin",
			Password: adminPassword.Value,
			Email:    "admin@cyber-man.pl",
		}
		slog.Info("no users found; creating admin user")
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
	s.SetupRoutes()
	fmt.Printf("listening on address '%s'\n", s.Address)
	chii.PrintRoutes(s.Router)
	err := http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Errorf("listen and serve failed s.Address='%s': %w", s.Address, err))
	}
}
