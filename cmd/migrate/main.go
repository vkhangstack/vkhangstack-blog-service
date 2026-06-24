package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/migrations"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	ctx := context.Background()
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("init: %v", err)
	}

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("up: %v", err)
		}
		if group.IsZero() {
			fmt.Println("no new migrations to run")
		} else {
			fmt.Printf("migrated to %s\n", group)
		}

	case "down":
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("down: %v", err)
		}
		if group.IsZero() {
			fmt.Println("no migrations to roll back")
		} else {
			fmt.Printf("rolled back %s\n", group)
		}

	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("status: %v", err)
		}
		for _, m := range ms {
			status := "pending"
			if m.IsApplied() {
				status = fmt.Sprintf("applied at %s", m.MigratedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("%-45s %s\n", m.Name, status)
		}

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate create <name>")
		}
		files, err := migrator.CreateSQLMigrations(ctx, os.Args[2])
		if err != nil {
			log.Fatalf("create: %v", err)
		}
		for _, f := range files {
			fmt.Printf("created: %s\n", f.Path)
		}

	default:
		fmt.Fprintln(os.Stderr, "usage: migrate <up|down|status|create <name>>")
		os.Exit(1)
	}
}
