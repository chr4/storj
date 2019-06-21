// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// Package attribution implements value attribution from docs/design/value-attribution.md
package attribution

import (
	"context"
	"time"

	"github.com/skyrings/skyring-common/tools/uuid"
)

// Info describing value attribution from partner to bucket
type Info struct {
	ProjectID  uuid.UUID
	BucketName []byte
	PartnerID  uuid.UUID
	CreatedAt  time.Time
}

// DB implements the database for value attribution table
type DB interface {
	// Get retrieves attribution info using project id and bucket name.
	Get(ctx context.Context, projectID uuid.UUID, bucketName []byte) (*Info, error)
	// Insert creates and stores new Info
	Insert(ctx context.Context, info *Info) (*Info, error)
}
