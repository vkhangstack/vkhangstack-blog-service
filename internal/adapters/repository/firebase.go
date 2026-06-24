package repository

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
)

func (f *DB) GetUser(ctx context.Context, uid string) *auth.UserRecord {
	// [START get_user_golang]
	// Get an auth client from the firebase.App
	client, err := f.adminFirebase.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
	}

	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Printf("error getting user %s: %v\n", uid, err)
	}
	// [END get_user_golang]
	return u
}

func (f *DB) ListUsers(ctx context.Context) {
	// [START list_all_users_golang]
	// Note, behind the scenes, the Users() iterator will retrive 1000 Users at a time through the API
	client, err := f.adminFirebase.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
	}

	iter := client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error listing users: %s\n", err)
		}
		log.Printf("read user user: %v\n", user)
	}

	// Iterating by pages 100 users at a time.
	// Note that using both the Next() function on an iterator and the NextPage()
	// on a Pager wrapping that same iterator will result in an error.
	pager := iterator.NewPager(client.Users(ctx, ""), 100, "")
	for {
		var users []*auth.ExportedUserRecord
		nextPageToken, err := pager.NextPage(&users)
		if err != nil {
			log.Fatalf("paging error %v\n", err)
		}
		for _, u := range users {
			log.Printf("read user user: %v\n", u)
		}
		if nextPageToken == "" {
			break
		}
	}
	// [END list_all_users_golang]
}
