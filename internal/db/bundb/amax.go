package bundb

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/id"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/uptrace/bun"
	"time"
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
			Where("? = ?", bun.Ident("amax.pub_key"), pubKey)

		if err := q.Scan(ctx); err != nil {
			return nil, a.conn.ProcessError(err)
		}

		return &amax, nil
	}, pubKey)
}

func (a *amaxDB) SubmitInfo(ctx context.Context, req *model.AmaxSubmitInfoRequest) (*gtsmodel.Amax, db.Error) {
	amax := &gtsmodel.Amax{}
	id, err := id.NewRandomULID()
	if err != nil {
		return nil, err
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

	// insert the new amax!
	if err := a.PutAmax(ctx, amax); err != nil {
		return nil, err
	}
	return amax, nil
}

func (a *amaxDB) PutAmax(ctx context.Context, amax *gtsmodel.Amax) db.Error {
	return a.state.Caches.GTS.Amax().Store(amax, func() error {
		_, err := a.conn.
			NewInsert().
			Model(amax).
			Exec(ctx)
		return a.conn.ProcessError(err)
	})
}

func (u *amaxDB) UpdateAmax(ctx context.Context, amax *gtsmodel.Amax, columns ...string) db.Error {
	// Update the amax's last-updated
	amax.UpdatedAt = time.Now()

	if len(columns) > 0 {
		// If we're updating by column, ensure "updated_at" is included
		columns = append(columns, "updated_at")
	}

	// Update the amax in DB
	_, err := u.conn.
		NewUpdate().
		Model(amax).
		Where("? = ?", bun.Ident("amax.pub_key"), amax.PubKey).
		Column(columns...).
		Exec(ctx)
	if err != nil {
		return u.conn.ProcessError(err)
	}

	// Invalidate user from cache
	u.state.Caches.GTS.Amax().Invalidate("PubKey", amax.PubKey)
	return nil
}
