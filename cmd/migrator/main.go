package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vitaliy-ukiru/test-bank/internal/config"
	"github.com/vitaliy-ukiru/test-bank/pkg/client/pg"
)

func main() {
	envPath := flag.String("env-path", "", "Path to .env file")
	flag.Parse()

	err := config.LoadConfig(envPath)
	if err != nil {
		panic(err)
	}

	cfg := config.Get()
	m, err := migrate.New(
		"file://migrations",
		pg.PgxConnString(
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Database,
			cfg.Database.Host,
			cfg.Database.Port,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	version, _, _ := m.Version()
	fmt.Println("migrated to %d", version)
}
