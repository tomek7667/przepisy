package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"przepisyapi/internal/przepisy"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/go-multi-logger-slog/logger"
)

type CliOptions struct {
	Address        string `env:"SECRETS_ADDRESS" envDefault:"127.0.0.1:7771"`
	DbPath         string `env:"SECRETS_DB_PATH" envDefault:"./przepisy.sqlite"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
	JwtSecret      string `env:"SECRETS_JWT_SECRET"`
	AdminPassword  string `env:"SECRETS_ADMIN_PASSWORD"`
}

func getJwtSecret() string {
	var secret string
	_, err := os.Stat(".jwtsecret")
	if err == nil || os.IsExist(err) {
		secret, _ = utils.MustReadFile(".jwtsecret")
	} else {
		fmt.Printf("empty jwt secret, creating .jwtsecret\n")
		secret = rand.Text()
		_ = os.WriteFile(".jwtsecret", []byte(secret), os.ModeAppend)
	}
	return secret
}

func main() {
	godotenv.Load()
	logger.SetLogLevel()
	var opts CliOptions

	// load defaults from env
	if err := env.Parse(&opts); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "secretsserver",
		Short: "Run secrets HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.JwtSecret == "" {
				opts.JwtSecret = getJwtSecret()
			}
			srv, err := przepisy.New(
				opts.Address,
				opts.AllowedOrigins,
				opts.DbPath,
				opts.JwtSecret,
				opts.AdminPassword,
			)
			if err != nil {
				return err
			}
			srv.Serve()
			return nil
		},
	}

	// flags override env/defaults
	rootCmd.Flags().StringVar(&opts.Address, "address", opts.Address, "listen address")
	rootCmd.Flags().StringVar(&opts.DbPath, "db-path", opts.DbPath, "path to sqlite db")
	rootCmd.Flags().StringVar(&opts.AllowedOrigins, "allowed-origins", opts.AllowedOrigins, "comma-separated list of allowed CORS origins")
	rootCmd.Flags().StringVar(&opts.JwtSecret, "jwt-secret", opts.JwtSecret, "jwt secret used for users session")
	rootCmd.Flags().StringVar(&opts.AdminPassword, "admin-password", opts.AdminPassword, "admin user password (generated randomly if not provided)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
