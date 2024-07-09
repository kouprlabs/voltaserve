package helper

import (
	"time"
	"voltaserve/infra"
)

func NewExpiry(token *infra.Token) time.Time {
	return time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
}
