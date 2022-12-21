package gtsmodel

import (
	"time"
)

type Amax struct {
	ID           string    `validate:"required,ulid" bun:"type:CHAR(26),pk,nullzero,notnull,unique"`        // id of this item in the database
	UserID       string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull,unique"`           // The id of the local gtsmodel.Account entry for this user.
	ClientID     string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull"`                  // ID of the client who owns this token
	CreatedAt    time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item created
	UpdatedAt    time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item last updated
	Scopes       string    `validate:"required" bun:",notnull"`                                             // scopes requested when this app was created
	RedirectURI  string    `validate:"required" bun:",notnull"`
	ResponseType string    `validate:"required" bun:",notnull"`
	PubKey       string    `validate:"required" bun:",notnull"`
	UserNam      string    `validate:"required" bun:",notnull"`
}
