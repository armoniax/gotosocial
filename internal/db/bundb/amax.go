package bundb

import (
	"context"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/uptrace/bun"
)

type amaxDB struct {
	conn  *DBConn
	state *state.State
}

func (a *amaxDB) GetAmaxByPubKey(ctx context.Context, pubKey string) (*gtsmodel.Amax, db.Error) {
	return a.state.Caches.GTS.Amax().Load("ID", func() (*gtsmodel.Amax, error) {
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
