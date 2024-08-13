package utils

import (
	"github.com/matoous/go-nanoid/v2"
)

func RandomID() (string, error) {
	return gonanoid.New()
}
