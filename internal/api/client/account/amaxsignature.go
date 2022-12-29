package account

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/armoniax/eos-go"
	"github.com/armoniax/eos-go/ecc"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/superseriousbusiness/gotosocial/internal/api"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"net/http"
	"strings"
)

const (
	EosURL = "https://test-chain.ambt.art"
)

func (m *Module) AccountSignaturePOSTHandler(c *gin.Context) {
	if _, err := api.NegotiateAccept(c, api.JSONAcceptHeaders...); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGet)
		return
	}

	form := &model.AmaxSignatureRequest{}
	if err := c.ShouldBind(form); err != nil {
		api.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGet)
		return
	}

	t, errWithCode := m.signature(c.Request.Context(), form)
	if errWithCode != nil {
		api.ErrorHandler(c, errWithCode, m.processor.InstanceGet)
		return
	}

	c.JSON(http.StatusOK, t)
}

func (m *Module) signature(ctx context.Context, form *model.AmaxSignatureRequest) (*token, gtserror.WithCode) {
	isSucc, username, pubkey, err := newTxSignDigest(ctx, form.Signature, form.Address, form.Scope, form.Authority, form.Expiration, form.ChainId)
	if err != nil {
		return nil, gtserror.NewError(err)
	}

	if !isSucc {
		return nil, gtserror.NewError(errors.Errorf("signature verify no pass"))
	}

	t, errWithCode := signatureLogin(getAddr(), username, pubkey)
	if errWithCode != nil {
		return nil, errWithCode
	}

	return t, nil
}

func signatureLogin(addr string, username, pubKey string) (*token, gtserror.WithCode) {
	data := make(map[string]any)
	data["username"] = username
	data["pub_key"] = pubKey

	return clientHttp[token]("POST", addr+SignatureLogin, data, nil, true)
}

// --------------------------signature logic----------------------------

type Identity struct {
	Scope      eos.Name            `json:"scope"`
	Permission eos.PermissionLevel `json:"permission"`
}

func newTxSignDigest(ctx context.Context, signature, address, scope, authority, expiration, chainId string) (bool, string, string, error) {
	digest, err := getTxSignDigest(address, scope, authority, expiration, chainId)
	if err != nil {
		return false, "", "", err
	}

	sig := ecc.MustNewAMASignature(signature)
	pk, err := sig.PublicKey(digest)
	if err != nil {
		return false, "", "", err
	}

	account, err := eos.New(EosURL).GetAccount(ctx, eos.AN(address))
	if err != nil {
		return false, "", "", err
	}

	isSucc, err := getIsVerified(account, authority, pk, digest, sig)
	return isSucc, account.AccountName.String(), pk.String(), err
}

func getIsVerified(account *eos.AccountResp, authority string, pk ecc.PublicKey, digest []byte, sig ecc.Signature) (isSuccess bool, err error) {
	isVerify := false
ACCOUNT:
	for _, perm := range account.Permissions {
		if perm.PermName != authority || len(perm.RequiredAuth.Keys) == 0 {
			continue
		}
		for _, key := range perm.RequiredAuth.Keys {
			if key.PublicKey.String() == pk.String() {
				isVerify = true
				break ACCOUNT
			}
		}
	}

	if !isVerify {
		return false, errors.New("verify signature error")
	}

	return sig.Verify(digest, pk), nil
}

func getTxSignDigest(address, scope, authority, expiration, chainId string) ([]byte, error) {
	signAction, err := getSignAction(address, scope, authority)
	if err != nil {
		return nil, err
	}

	tx := eos.NewTransaction([]*eos.Action{signAction}, nil)
	if tx.Expiration, err = eos.ParseJSONTime(expiration); err != nil {
		return nil, err
	}

	txData, _, err := eos.NewSignedTransaction(tx).PackedTransactionAndCFD()
	if err != nil {
		fmt.Printf("packed tx error: %v\n", err)
		return nil, err
	}

	chainID, err := hex.DecodeString(chainId)
	if err != nil {
		fmt.Printf("hex decode error: %v\n", err)
		return nil, err
	}

	return eos.SigDigest(chainID, txData, nil), nil
}

func getSignAction(address, scope, authority string) (*eos.Action, error) {
	data, err := getAbiData(address, scope, authority)
	if err != nil {
		return nil, err
	}
	actionData := eos.NewActionData(hex.EncodeToString(data))
	return newSignAction(eos.AN(address), eos.AccountName(""), "anchor.link.d", actionData), nil
}

func newSignAction(payer, account eos.AccountName, scope string, actionData eos.ActionData) *eos.Action {
	return &eos.Action{
		Account: account,
		Name:    eos.ActN("identity"),
		Authorization: []eos.PermissionLevel{
			{Actor: payer, Permission: eos.PN("active")},
		},
		ActionData: actionData,
	}
}

func getAbiData(address, scope, authority string) ([]byte, error) {
	abi, err := eos.NewABI(strings.NewReader(abiString))
	if err != nil {
		fmt.Printf("new abi error: %s\n", err.Error())
		return nil, err
	}

	identity := Identity{
		Scope: eos.Name(scope),
		Permission: eos.PermissionLevel{
			Actor:      eos.AN(address),
			Permission: eos.PN(authority),
		},
	}

	abiData, err := json.Marshal(identity)
	if err != nil {
		fmt.Printf("json.Marshal failed: %v", err)
		return nil, err
	}

	return abi.EncodeStruct("identity", abiData)
}

// 客户端扫码签名后返回的结果
var signRes = `
{
  "chainId":"208dacab3cd2e181c86841613cf05d9c60786c677e4ce86b266d0a58884968f7",
  "scope":"anchor.link.d",
  "expiration":"2022-10-09T03:33:13",
  "signer":{
      "actor":"merchantx",
      "permission":"active"
  },
  "signature":"SIG_K1_K9MXXWBPBgEV36XRoUwnJvpfRhBDvxA7GTiu69cRdQm3HF6KRgbjTkyyjqR3Rc1FUipM4j5ZjSuopwUSMXfYpmrCVJUCaC"
}
`

var abiString = `
{
  "version": "app::abi/1.1",
  "types": [
      {
          "new_type_name": "chain_alias",
          "type": "uint8"
      },
      {
          "new_type_name": "chain_id",
          "type": "checksum256"
      },
      {
          "new_type_name": "request_flags",
          "type": "uint8"
      }
  ],
  "structs": [
      {
          "base": "",
          "name": "permission_level",
          "fields": [
              {
                  "name": "actor",
                  "type": "name"
              },
              {
                  "name": "permission",
                  "type": "name"
              }
          ]
      },
      {
          "base": "",
          "name": "action",
          "fields": [
              {
                  "name": "account",
                  "type": "name"
              },
              {
                  "name": "name",
                  "type": "name"
              },
              {
                  "name": "authorization",
                  "type": "permission_level[]"
              },
              {
                  "name": "data",
                  "type": "identity"
              }
          ]
      },
      {
          "base": "",
          "name": "identity",
          "fields": [
              {
                  "name": "scope",
                  "type": "name"
              },
              {
                  "name": "permission",
                  "type": "permission_level?"
              }
          ]
      }
  ],
  "actions": [ ],
  "tables": [ ],
  "ricardian_clauses": [ ],
  "error_messages": [ ],
  "abi_extensions": [ ],
  "variants": [
      {
          "name": "variant_id",
          "types": [
              "chain_alias",
              "chain_id"
          ]
      },
      {
          "name": "variant_req",
          "types": [
              "action",
              "action[]",
              "transaction",
              "identity"
          ]
      }
  ]
}
`
