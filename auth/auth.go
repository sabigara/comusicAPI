package auth

import (
	"context"
	"errors"

	firebase "firebase.google.com/go"
	comusic "github.com/sabigara/comusicAPI"
)

var fbapp *firebase.App

func Authenticate(credentials ...interface{}) (*comusic.User, error) {
	if len(credentials) != 1 {
		return nil, errors.New("auth: invalid credential")
	}
	idToken, ok := credentials[0].(string)
	if !ok {
		return nil, errors.New("auth: invalid credential")
	}

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
	fbapp, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}
}
