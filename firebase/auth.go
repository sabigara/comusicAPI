package firebase

import (
	"context"
	"strings"
	"time"

	fb "firebase.google.com/go"
	comusic "github.com/sabigara/comusicAPI"
)

var fbapp *fb.App

func Authenticate(credentials ...interface{}) (*comusic.User, error) {
	if len(credentials) != 1 {
		return nil, comusic.ErrUnauthenticated
	}
	credential, ok := credentials[0].(string)
	if !ok {
		return nil, comusic.ErrUnauthenticated
	}
	tokenSlice := strings.Split(credential, "Bearer ")
	if len(tokenSlice) != 2 {
		return nil, comusic.ErrUnauthenticated
	}
	idToken := tokenSlice[1]
	ctx := context.Background()
	client, err := fbapp.Auth(ctx)
	if err != nil {
		return nil, comusic.ErrAuthProcessFailed
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, comusic.ErrUnauthenticated
	}
	user, err := client.GetUser(ctx, token.UID)
	if err != nil {
		return nil, comusic.ErrAuthProcessFailed
	}
	// TODO: UpdateTimestamp is not provided from firebase.
	return &comusic.User{
		Meta: &comusic.Meta{
			ID: user.UID,
			// CreationTimestamp is in milliseconds so divide by 1000.
			CreatedAt: time.Unix(user.UserMetadata.CreationTimestamp/1000, 0),
		},
		Email: user.Email,
	}, nil
}

func init() {
	var err error
	fbapp, err = fb.NewApp(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}
}
