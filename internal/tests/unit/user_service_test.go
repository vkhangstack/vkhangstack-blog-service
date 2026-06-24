package unit

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/cache"
	"github.com/vkhangstack/hexagonal-architecture/internal/adapters/repository"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func setUpDB() *repository.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://test:test@localhost:5433/template1?sslmode=disable"),
	))
	db := bun.NewDB(sqldb, pgdialect.New())

	ctx := context.Background()
	for _, model := range []interface{}{
		(*domain.Message)(nil),
		(*domain.Customer)(nil),
		(*domain.Payment)(nil),
	} {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			panic(err)
		}
	}

	redisCache, err := cache.NewRedisCache("localhost:6379", "")
	if err != nil {
		panic(err)
	}

	return repository.NewDB(db, redisCache, nil, nil)
}

func TestCreateUser(t *testing.T) {
	db := setUpDB()

	email := "alanmoore@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
}

/*
func TestReadUser(t *testing.T) {
	db := setUpDB()

	email := "test@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	cachedUser, err := db.ReadUser(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, cachedUser)
	assert.Equal(t, user.ID, cachedUser.ID)
	assert.Equal(t, user.Email, cachedUser.Email)
	assert.Equal(t, user.Password, cachedUser.Password)

	time.Sleep(time.Second * 3)

	cachedUser, err = db.ReadUser(user.ID)
	assert.Error(t, err)
	assert.Nil(t, cachedUser)
}

func TestReadUsers(t *testing.T) {
	db := setUpDB()

	email := "test@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	users, err := db.ReadUsers()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.NotEmpty(t, users)
}

func TestUpdateUser(t *testing.T) {
	db := setUpDB()

	email := "test@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	newEmail := "new@example.com"
	newPassword := "newpassword"

	err = db.UpdateUser(user.ID, newEmail, newPassword)
	assert.NoError(t, err)

	cachedUser, err := db.ReadUser(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, cachedUser)
	assert.Equal(t, newEmail, cachedUser.Email)
	assert.NotEqual(t, password, cachedUser.Password)
}

func TestDeleteUser(t *testing.T) {
	db := setUpDB()

	email := "test@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	err = db.DeleteUser(user.ID)
	assert.NoError(t, err)

	cachedUser, err := db.ReadUser(user.ID)
	assert.Error(t, err)
	assert.Nil(t, cachedUser)

	users, err := db.ReadUsers()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Empty(t, users)
}

func TestCreateUserAlreadyExists(t *testing.T) {
	db := setUpDB()

	email := "test@example.com"
	password := "password"

	user, err := db.CreateUser(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = db.CreateUser(email, password)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestReadUserNotFound(t *testing.T) {
	db := setUpDB()

	user, err := db.ReadUser(0)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUpdateUserNotFound(t *testing.T) {
	db := setUpDB()

	err := db.UpdateUser("nonexistent", "new@example.com", "newpassword")
	assert.Error(t, err)
}

func TestDeleteUserNotFound(t *testing.T) {
	db := setUpDB()

	err := db.DeleteUser(0)
	assert.Error(t, err)
}
*/
