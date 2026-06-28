package repository

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
	"google.golang.org/api/iterator"
)

func (f *DB) GetUser(ctx context.Context, uid string) *auth.UserRecord {
	client, err := f.adminFirebase.Auth(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("error getting firebase auth client")
	}

	u, err := client.GetUser(ctx, uid)
	if err != nil {
		logger.Log.WithError(err).WithField("uid", uid).Error("error getting firebase user")
	}
	return u
}

func (f *DB) ListUsers(ctx context.Context) {
	client, err := f.adminFirebase.Auth(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("error getting firebase auth client")
	}

	iter := client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Log.WithError(err).Fatal("error listing firebase users")
		}
		logger.Log.WithField("user", user).Debug("firebase user")
	}

	pager := iterator.NewPager(client.Users(ctx, ""), 100, "")
	for {
		var users []*auth.ExportedUserRecord
		nextPageToken, err := pager.NextPage(&users)
		if err != nil {
			logger.Log.WithError(err).Fatal("firebase paging error")
		}
		for _, u := range users {
			logger.Log.WithField("user", u).Debug("firebase user")
		}
		if nextPageToken == "" {
			break
		}
	}
}
