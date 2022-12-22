package db

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type Amax interface {
	GetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, Error)

	SubmitInfo(ctx context.Context, request *model.AmaxSubmitInfoRequest) (*gtsmodel.Amax, Error)

	PutAmax(ctx context.Context, user *gtsmodel.Amax) Error

	UpdateAmax(ctx context.Context, amax *gtsmodel.Amax, columns ...string) Error
}
