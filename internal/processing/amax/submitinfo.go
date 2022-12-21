package amax

import (
	"context"
	"errors"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

func (p *processor) SubmitInfo(ctx context.Context, amax *gtsmodel.Amax, request *model.AmaxSubmitInfoRequest) gtserror.WithCode {
	if request == nil {
		return gtserror.NewErrorGone(errors.New("amax request is nil"))
	}

	if amax == nil {
		amax = new(gtsmodel.Amax)
	}

	if _, err := p.db.SubmitInfo(ctx, request.UserID, request.ClientID, request.RedirectUri, request.ResponseType, request.Scopes, request.PubKey, request.Username); err != nil {
		return gtserror.NewErrorGone(err, "processor submitInfo create failed")
	}

	return nil
}
