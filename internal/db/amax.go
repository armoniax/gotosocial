package db

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type Amax interface {
	GetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, Error)
}
