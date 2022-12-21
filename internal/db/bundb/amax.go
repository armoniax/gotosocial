package bundb

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/log"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/uptrace/bun"
)

type amaxDB struct {
	conn  *DBConn
	state *state.State
}

func (a *amaxDB) GetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, db.Error) {
	return a.state.Caches.GTS.Amax().Load("PubKey", func() (*gtsmodel.Amax, error) {
		var amax gtsmodel.Amax

		q := a.conn.
			NewSelect().
			Model(&amax).
			Relation("Account").
			Where("? = ?", bun.Ident("user.pub_key"), pubKey)

		if err := q.Scan(ctx); err != nil {
			return nil, a.conn.ProcessError(err)
		}

		return &amax, nil
	}, pubKey)
}

func (a *amaxDB) SubmitInfo(ctx context.Context, userID, clientID, redirectUri, responseType, scopes, pubKey, username string) (*gtsmodel.Amax, db.Error) {
	// if something went wrong while creating a user, we might already have an account, so check here first...
	log.Infof("userId: %v, clientID: %v, redirectUrl: %v,responseType: %v, scposes: %v, pubkey: %v, username: %v", userID, clientID, redirectUri, responseType, scopes, pubKey, username)

	amax := &gtsmodel.Amax{}
	//if err := a.conn.
	//	NewSelect().
	//	Model(amax).
	//	Where("? = ?", bun.Ident("amax.pub_key"), pubKey).
	//	Scan(ctx); err != nil {
	//	err = a.conn.ProcessError(err)
	//	if err != db.ErrNoEntries {
	//		log.Errorf("error checking for existing account: %s", err)
	//		return nil, err
	//	}
	//}

	amax.UserID = userID
	amax.ClientID = clientID
	amax.RedirectURI = redirectUri
	amax.ResponseType = responseType
	amax.Scopes = scopes
	amax.PubKey = pubKey
	amax.Username = username

	// insert the new account!
	if err := a.PutAmax(ctx, amax); err != nil {
		return nil, err
	}
	return amax, nil
}

func (a *amaxDB) PutAmax(ctx context.Context, amax *gtsmodel.Amax) db.Error {
	err := a.state.Caches.GTS.Amax().Store(amax, func() error {
		_, err := a.conn.
			NewInsert().
			Model(amax).
			Exec(ctx)
		return a.conn.ProcessError(err)
	})
	log.Errorf("mark: PutAmax error: %v", err)
	return err
}
