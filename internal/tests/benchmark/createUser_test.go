package benchmark

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/cache"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/repository"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func BenchmarkCreateUser(b *testing.B) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://postgres:password@localhost:5432/postgres?sslmode=disable"),
	))
	db := bun.NewDB(sqldb, pgdialect.New())

	redisCache, err := cache.NewRedisCache("localhost:6379", "")
	if err != nil {
		panic(err)
	}

	store := repository.NewDB(db, redisCache, nil, nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("test_user_%d@example.com", i)
		password := "password"
		// Delete user if it exists
		var customer domain.Customer
		if err := db.NewSelect().Model(&customer).Where("email = ?", email).Limit(1).Scan(ctx); err == nil {
			if _, err := db.NewDelete().Model(&customer).Where("email = ?", email).Exec(ctx); err != nil {
				b.Fatalf("failed to delete customer: %v", err)
			}
		}
		_, err := store.CreateUser(email, password)
		if err != nil {
			b.Fatalf("failed to create test customer: %v", err)
		}
	}
}
