package bundb

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
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

	// Invalidate amax from cache
	u.state.Caches.GTS.Amax().Invalidate("PubKey", amax.PubKey)
	return nil
}
