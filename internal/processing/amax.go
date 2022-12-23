package processing

import (
	"context"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

func (p *processor) AmaxSubmitInfo(ctx context.Context, form *apimodel.AmaxSubmitInfoRequest) (*gtsmodel.Amax, gtserror.WithCode) {
	return p.amaxProcessor.SubmitInfo(ctx, form)
}

func (p *processor) AmaxGetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, gtserror.WithCode) {
	amax, err := p.db.GetAmaxByPubKey(ctx, pubKey)
	if err != nil {
		return nil, gtserror.NewErrorGone(err)
	}
	return amax, nil
}
