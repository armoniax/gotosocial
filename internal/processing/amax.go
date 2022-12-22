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
