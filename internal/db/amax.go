package db

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type Amax interface {
	GetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, Error)

	SubmitInfo(ctx context.Context, userID, clientID, redirectUri, responseType, scopes, pubKey, username string) (*gtsmodel.Amax, Error)

	PutAmax(ctx context.Context, user *gtsmodel.Amax) Error
}
