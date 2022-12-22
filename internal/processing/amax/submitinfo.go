package amax

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/id"
	"time"
)

func (p *processor) SubmitInfo(ctx context.Context, request *model.AmaxSubmitInfoRequest) (*gtsmodel.Amax, gtserror.WithCode) {
	if request == nil || len(request.Password) == 0 {
		return nil, gtserror.NewErrorGone(errors.New("amax request is nil"))
	}

	amax, err := p.db.GetAmaxByPubKey(ctx, request.Password)
	if err != nil {
		return nil, gtserror.NewErrorGone(errors.Errorf("db GetAmaxByPubKey failed: %v", request.Password))
	}

	if amax != nil {
		return nil, gtserror.NewErrorGone(errors.Errorf("db amax has existed: %v", request.Password))
	}

	err = p.db.PutAmax(ctx, composeAmax(request))
	if err != nil {
		return nil, gtserror.NewErrorGone(errors.Errorf("db PutAmax failed: %v", err))
	}

	amax, err = p.db.GetAmaxByPubKey(ctx, request.Password)
	return amax, gtserror.NewErrorGone(err)
}

func composeAmax(req *model.AmaxSubmitInfoRequest) *gtsmodel.Amax {
	amax := &gtsmodel.Amax{}
	id, err := id.NewRandomULID()
	if err != nil {
		return nil
	}

	amax.ID = id
	amax.CreatedAt = time.Now()
	amax.UpdatedAt = time.Now()
	amax.ClientName = req.ClientName
	amax.RedirectUri = req.RedirectUris
	amax.Scope = req.Scope
	amax.GrantType = req.GrantType
	amax.ClientId = req.ClientId
	amax.ClientSecret = req.ClientSecret
	amax.Reason = req.Reason
	amax.Email = req.Email
	amax.Username = req.Username
	amax.Agreement = req.Agreement
	amax.Locale = req.Locale
	amax.PubKey = req.Password
	return amax
}
