//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
//go:generate go run github.com/pressly/goose/v3/cmd/goose -v -dir schema sqlite3 ./internal/sqlite/przepisy.sqlite up
package generate
