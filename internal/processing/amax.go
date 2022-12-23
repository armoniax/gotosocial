package processing

import (
	"context"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/db"
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

func (p *processor) AmaxSignatureLogin(ctx context.Context, form *apimodel.AmaxSignatureLoginRequest) (*gtsmodel.User, gtserror.WithCode) {
	amax, err := p.db.GetAmaxByPubKey(ctx, form.PubKey)
	switch err {
	case nil:
		return p.login(ctx, amax)
	case db.ErrNoEntries:
		return p.register(ctx, form)
	default:
		return nil, gtserror.NewErrorGone(err)
	}
}

func (p *processor) register(ctx context.Context, form *apimodel.AmaxSignatureLoginRequest) (*gtsmodel.User, gtserror.WithCode) {
	return nil, nil
}

func (p *processor) login(ctx context.Context, amax *gtsmodel.Amax) (*gtsmodel.User, gtserror.WithCode) {
	return nil, nil
}
