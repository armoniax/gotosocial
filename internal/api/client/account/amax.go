package account

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/superseriousbusiness/gotosocial/internal/api"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"net/http"
)

func (m *Module) AccountCreateAmaxInfoPOSTHandler(c *gin.Context) {
	_, err := oauth.Authed(c, true, true, true, true)
	if err != nil {
		api.ErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGet)
		return
	}

	if _, err := api.NegotiateAccept(c, api.JSONAcceptHeaders...); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGet)
		return
	}

	form := &model.AmaxSubmitInfoRequest{}
	if err := c.ShouldBind(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	if err := validateCreateAmax(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	amax, errWithCode := m.processor.AmaxSubmitInfo(c.Request.Context(), form)
	if errWithCode != nil {
		api.ErrorHandler(c, errWithCode, m.processor.InstanceGet)
		return
	}
	c.JSON(http.StatusOK, amax)
}

func validateCreateAmax(form *model.AmaxSubmitInfoRequest) error {
	if form == nil {
		return errors.New("form is nil")
	}

	if len(form.ClientId) != 26 {
		return errors.Errorf("client_id length is not correct: %v", form.ClientId)
	}

	return nil
}
