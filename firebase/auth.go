package firebase

import (
	"context"
	"errors"
	"strings"

	fb "firebase.google.com/go"
	comusic "github.com/sabigara/comusicAPI"
)

var fbapp *fb.App

func Authenticate(credentials ...interface{}) (*comusic.User, error) {
	if len(credentials) != 1 {
		return nil, errors.New("auth: invalid credential")
	}
	credential, ok := credentials[0].(string)
	if !ok {
		return nil, errors.New("auth: invalid credential")
	}
	tokenSlice := strings.Split(credential, "Bearer ")
	if len(tokenSlice) != 2 {
		return nil, errors.New("auth: malformed credential")
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
	return comusic.NewUser(user.UID, user.DisplayName, user.Email), nil
}

func init() {
	var err error
	fbapp, err = fb.NewApp(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}
}
