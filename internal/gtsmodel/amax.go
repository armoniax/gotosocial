package gtsmodel

import (
	"time"
)

type Amax struct {
	ID           string    `validate:"required,ulid" bun:"type:CHAR(26),pk,nullzero,notnull,unique"` // id of this item in the database
	ClientName   string    `validate:"required" bun:",notnull"`
	RedirectUri  string    `validate:"required" bun:",notnull"`
	Scope        string    `validate:"required" bun:",notnull"`
	GrantType    string    `validate:"required" bun:",notnull"`
	ClientId     string    `validate:"required,ulid" bun:"type:CHAR(26),nullzero,notnull"` // ID of the client who owns this token
	ClientSecret string    `validate:"required" bun:",notnull"`
	Reason       string    `validate:"required" bun:",notnull"`
	Email        string    `validate:"required" bun:",notnull"`
	Username     string    `validate:"required" bun:",notnull"`
	Agreement    bool      `validate:"required" bun:",notnull"`
	Locale       string    `validate:"required" bun:",notnull"`
	PubKey       string    `validate:"required" bun:",notnull"`
	CreatedAt    time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item created
	UpdatedAt    time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"` // when was item last updated
}
