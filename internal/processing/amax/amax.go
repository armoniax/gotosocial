package amax

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

// Processor wraps a bunch of functions for processing user-level actions.
type Processor interface {
	SubmitInfo(ctx context.Context, amax *gtsmodel.Amax, request *model.AmaxSubmitInfoRequest) gtserror.WithCode
}

type processor struct {
	db db.DB
}

// New returns a new user processor
func New(db db.DB) Processor {
	return &processor{db: db}
}
