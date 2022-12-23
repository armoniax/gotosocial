package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/superseriousbusiness/gotosocial/internal/api"
	"github.com/superseriousbusiness/gotosocial/internal/api/client/app"
	"github.com/superseriousbusiness/gotosocial/internal/api/client/auth"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/log"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"io"
	"net"
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

func (m *Module) AccountCreateUserTokenPOSTHandler(c *gin.Context) {
	authed, err := oauth.Authed(c, true, true, false, false)
	if err != nil {
		api.ErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGet)
		return
	}

	if _, err := api.NegotiateAccept(c, api.JSONAcceptHeaders...); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGet)
		return
	}

	form := &model.AccountCreateRequest{}
	if err := c.ShouldBind(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	if err := validateCreateAccount(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	clientIP := c.ClientIP()
	signUpIP := net.ParseIP(clientIP)
	if signUpIP == nil {
		err := errors.New("ip address could not be parsed from request")
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}
	form.IP = signUpIP

	ti, errWithCode := m.processor.AccountCreateUserToken(c.Request.Context(), authed, form)
	if errWithCode != nil {
		api.ErrorHandler(c, errWithCode, m.processor.InstanceGet)
		return
	}

	c.JSON(http.StatusOK, ti)
}

func (m *Module) AccountSignatureLoginPOSTHandler(c *gin.Context) {
	if _, err := api.NegotiateAccept(c, api.JSONAcceptHeaders...); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGet)
		return
	}

	form := &model.AmaxSignatureLoginRequest{}
	if err := c.ShouldBind(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	if err := validateSignatureLoginReq(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	user, errWithCode := m.amaxSignatureLogin(c.Request.Context(), form)
	//user, errWithCode := m.processor.AmaxSignatureLogin(c.Request.Context(), form)
	if errWithCode != nil {
		api.ErrorHandler(c, errWithCode, m.processor.InstanceGet)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (m *Module) amaxSignatureLogin(ctx context.Context, form *model.AmaxSignatureLoginRequest) (*model.Account, gtserror.WithCode) {
	if len(form.PubKey) == 0 {
		return nil, gtserror.NewError(errors.New("form PubKey is empty"))
	}

	notFound := "no entries"
	amax, err := m.processor.AmaxGetAmaxByPubKey(ctx, form.PubKey)

	switch {
	case err != nil && err.Error() == notFound:
		return m.register(form)
	case err == nil && amax != nil:
		return m.login(amax)
	default:
		return nil, err
	}
}

func (m *Module) register(form *model.AmaxSignatureLoginRequest) (account *model.Account, errs gtserror.WithCode) {
	defer func() {
		if account == nil {
			deleteRegisterAllInfo()
		}
	}()

	bindAddress := "http://localhost"
	port := config.GetPort()
	addr := fmt.Sprintf("%s:%d", bindAddress, port)

	//# Step 1: create the app to register the new account
	app, err := createApplication(addr)
	if err != nil {
		errs = err
		return
	}

	//# Step 2: obtain a code for that app
	appt1, err := createAppToken(addr, app.ClientID, app.ClientSecret)
	if err != nil {
		errs = err
		return
	}

	//# Step 3: use the code to register a new account
	appt2, err := createUser(addr, appt1.AccessToken, form.Username, form.PubKey)
	if err != nil {
		errs = err
		return
	}

	//# Step 4: verify the returned access token
	account, err = verifyCredentials(addr, appt2.AccessToken)
	if err != nil {
		return
	}

	//# Step 5: store amax core info
	amax := model.AmaxSubmitInfoRequest{}
	amax.ClientName = app.Name
	amax.RedirectUris = app.RedirectURI
	amax.Scope = appt1.Scope
	amax.GrantType = "client_credentials"
	amax.ClientId = app.ClientID
	amax.ClientSecret = app.ClientSecret
	amax.Reason = "Testing whether or not this dang diggity thing works!"
	amax.Email = form.PubKey + "@amax.com"
	amax.Username = form.Username
	amax.Password = form.PubKey
	amax.Agreement = true
	amax.Locale = "en"
	if err = createAmaxInfo(addr, appt2.AccessToken, &amax); err != nil {
		return nil, gtserror.NewError(err)
	}
	return account, nil
}

func createApplication(addr string) (*model.Application, gtserror.WithCode) {
	data := make(map[string]any)
	data["client_name"] = "amax"
	data["redirect_uris"] = addr
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	resp, err := http.Post(addr+app.BasePath, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	var cnt bytes.Buffer
	if _, err = io.Copy(&cnt, resp.Body); err != nil {
		log.Errorf("application copy failed: %v", err)
		return nil, gtserror.NewError(err)
	}

	app := model.Application{}
	if err = json.Unmarshal(cnt.Bytes(), &app); err != nil {
		log.Errorf("application Unmarshal failed: %v", err)
		return nil, gtserror.NewError(err)
	}
	return &app, nil
}

type appToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func createAppToken(addr, clientId, clientSecret string) (*appToken, gtserror.WithCode) {
	data := make(map[string]any)
	data["scope"] = "read"
	data["grant_type"] = "client_credentials"
	data["client_id"] = clientId
	data["client_secret"] = clientSecret
	data["redirect_uri"] = addr
	bytesData, err := json.Marshal(data)

	resp, err := http.Post(addr+auth.OauthTokenPath, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	var cnt bytes.Buffer
	if _, err = io.Copy(&cnt, resp.Body); err != nil {
		log.Errorf("application auth failed: %v", err)
		return nil, gtserror.NewError(err)
	}

	appt := appToken{}
	if err = json.Unmarshal(cnt.Bytes(), &appt); err != nil {
		log.Errorf("application auth Unmarshal failed: %v", err)
		return nil, gtserror.NewError(err)
	}
	return &appt, nil
}

func createUser(addr, authStr, username, pubKey string) (*appToken, gtserror.WithCode) {
	data := make(map[string]any)
	data["reason"] = "Testing whether or not this dang diggity thing works!"
	data["username"] = username
	data["email"] = pubKey + "@amax.com"
	data["password"] = pubKey
	data["agreement"] = true
	data["locale"] = "en"
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	log.Infof("the value: %+v", data)
	req, err := http.NewRequest("POST", addr+BasePath, bytes.NewReader(bytesData))
	if err != nil {
		return nil, gtserror.NewError(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authStr)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	var cnt bytes.Buffer
	if _, err = io.Copy(&cnt, resp.Body); err != nil {
		log.Errorf("create user copy failed: %v", err)
		return nil, gtserror.NewError(err)
	}

	appt := appToken{}
	if err = json.Unmarshal(cnt.Bytes(), &appt); err != nil {
		log.Infof("create user Unmarshal: %v", cnt.String())
		log.Errorf("create user Unmarshal failed: %v", err)
		return nil, gtserror.NewError(err)
	}
	return &appt, nil
}

func verifyCredentials(addr, authStr string) (*model.Account, gtserror.WithCode) {
	req, err := http.NewRequest("GET", addr+VerifyPath, http.NoBody)
	if err != nil {
		return nil, gtserror.NewError(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authStr)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	var cnt bytes.Buffer
	if _, err = io.Copy(&cnt, resp.Body); err != nil {
		log.Errorf("create user copy failed: %v", err)
		return nil, gtserror.NewError(err)
	}

	account := model.Account{}
	if err = json.Unmarshal(cnt.Bytes(), &account); err != nil {
		log.Errorf("create user Unmarshal failed: %v", err)
		return nil, gtserror.NewError(err)
	}
	return &account, nil
}

func createAmaxInfo(addr, authStr string, amax *model.AmaxSubmitInfoRequest) gtserror.WithCode {
	bytesData, err := json.Marshal(amax)
	req, err := http.NewRequest("POST", addr+SubmitAmaxInfo, bytes.NewReader(bytesData))
	if err != nil {
		return gtserror.NewError(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authStr)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return gtserror.NewError(err)
	}

	return nil
}

func deleteRegisterAllInfo() {
	//del kind of table info
}

func (m *Module) login(amax *gtsmodel.Amax) (*model.Account, gtserror.WithCode) {
	return &model.Account{
		ID: "login ....",
	}, nil
}

func validateSignatureLoginReq(form *model.AmaxSignatureLoginRequest) error {
	if form == nil {
		return errors.New("form is nil")
	}

	if len(form.Username) == 0 {
		return errors.New("Username is empty")
	}

	if len(form.PubKey) == 0 {
		return errors.New("PubKey is empty")
	}

	return nil
}
