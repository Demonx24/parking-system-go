package request

import "parking-system-go/model/database"

type User struct{}
type Login struct {
	User                 database.User `json:"user"`
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt int64         `json:"access_token_expires_at"`
}
