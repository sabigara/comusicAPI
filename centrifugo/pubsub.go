package centrifugo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/centrifugal/gocent"
	"github.com/dgrijalva/jwt-go"
	comusic "github.com/sabigara/comusicAPI"
)

type PubSub struct {
	Centr     *gocent.Client
	JWTSecret string
}

func NewPubSub(c *gocent.Client, secret string) *PubSub {
	return &PubSub{Centr: c, JWTSecret: secret}
}

func (ps *PubSub) GenerateToken(user *comusic.User) (token string, err error) {
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
	}).SignedString([]byte(ps.JWTSecret))
	return
}

func (ps *PubSub) Publish(channel string, message interface{}) error {
	dataBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("centrifugo.Publish: failed to marshal message struct%w", err)
	}
	ctx := context.Background()
	err = ps.Centr.Publish(ctx, channel, dataBytes)
	if err != nil {
		return fmt.Errorf("centrifugo.Publish: failed to publish: %w", err)
	}
	return err
}
