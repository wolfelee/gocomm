package juuid

import (
	"github.com/renstrom/shortuuid"
)

func ShortUUID() string {
	return shortuuid.New()
}
