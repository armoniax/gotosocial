package processing

import (
	"context"
	"github.com/go-errors/errors"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
)

func (p *processor) AmaxSubmitInfo(ctx context.Context, authed *oauth.Auth, form *apimodel.AmaxSubmitInfoRequest) gtserror.WithCode {
	if len(form.UserID) == 0 && len(form.PubKey) > 0 {
		if user, err := p.db.GetUserByUnconfirmedEmail(ctx, form.PubKey+"@amax.com"); err != nil {
			return gtserror.NewErrorGone(errors.Errorf("amax GetUserByUnconfirmedEmail failed: %v", form.PubKey))
		} else {
			form.UserID = user.ID
		}
	}
	return p.amaxProcessor.SubmitInfo(ctx, authed.Amax, form)
}
