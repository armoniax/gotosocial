package processing

import (
	"context"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
)

func (p *processor) AmaxSubmitInfo(ctx context.Context, authed *oauth.Auth, form *apimodel.AmaxSubmitInfoRequest) gtserror.WithCode {
	return p.amaxProcessor.SubmitInfo(ctx, authed.Amax, form)
}
