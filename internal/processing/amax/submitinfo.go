package amax

import (
	"context"
	"errors"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

func (p *processor) SubmitInfo(ctx context.Context, request *model.AmaxSubmitInfoRequest) (*gtsmodel.Amax, gtserror.WithCode) {
	if request == nil {
		return nil, gtserror.NewErrorGone(errors.New("amax request is nil"))
	}

	amax, err := p.db.SubmitInfo(ctx, request)
	return amax, gtserror.NewErrorGone(err)
}
