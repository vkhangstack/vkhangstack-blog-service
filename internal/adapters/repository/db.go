package repository

import (
	firebase "firebase.google.com/go/v4"
	"github.com/uptrace/bun"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/cache"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/snowflake"
)

type DB struct {
	db            *bun.DB
	cache         *cache.RedisCache
	adminFirebase *firebase.App
	snowflakeNode *snowflake.Node
}

func NewDB(db *bun.DB, cache *cache.RedisCache, adminFirebase *firebase.App, snowflakeNode *snowflake.Node) *DB {
	return &DB{
		db:            db,
		cache:         cache,
		adminFirebase: adminFirebase,
		snowflakeNode: snowflakeNode,
	}
}
