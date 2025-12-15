package main

import (
	"log"

	"przepisyapi/internal/przepisy"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tomek7667/go-multi-logger-slog/logger"
	"github.com/tomek7667/secrets/secretssdk"
)

var secretsClient *secretssdk.Client

type CliOptions struct {
	FrontendUrl    string `env:"FRONTEND_URL" envDefault:"https://przepisy.cyber-man.pl"`
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:7771"`
	DbPath         string `env:"DB_PATH" envDefault:"./przepisy.sqlite"`
	SecretsAddress string `env:"SECRETS_ADDRESS" envDefault:"https://secrets.cyber-man.pl"`
	SecretsToken   string `env:"SECRETS_TOKEN"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
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
			var err error
			secretsClient, err = secretssdk.New(opts.SecretsAddress, opts.SecretsToken)
			if err != nil {
				return err
			}

			srv, err := przepisy.New(
				opts.FrontendUrl,
				opts.Address,
				opts.AllowedOrigins,
				opts.DbPath,
				secretsClient,
			)
			if err != nil {
				return err
			}
			srv.Serve()
			return nil
		},
	}

	// flags override env/defaults
	rootCmd.Flags().StringVar(&opts.FrontendUrl, "frontend-url", opts.FrontendUrl, "frontend url to send the users to")
	rootCmd.Flags().StringVar(&opts.Address, "address", opts.Address, "listen address")
	rootCmd.Flags().StringVar(&opts.DbPath, "db-path", opts.DbPath, "path to sqlite db")
	rootCmd.Flags().StringVar(&opts.AllowedOrigins, "allowed-origins", opts.AllowedOrigins, "comma-separated list of allowed CORS origins")
	rootCmd.Flags().StringVar(&opts.SecretsAddress, "secrets-address", opts.SecretsAddress, "address of secrets service (default local if not specified)")
	rootCmd.Flags().StringVarP(&opts.SecretsToken, "secrets-token", "t", opts.SecretsToken, "token for retrieving secrets")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
