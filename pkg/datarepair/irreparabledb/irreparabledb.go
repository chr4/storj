// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package irreparabledb

import (
	"context"

	"github.com/zeebo/errs"
)

// Error is the default irreparabledb errs class
var Error = errs.Class("irreparabledb error")

// DB interface for database operations
type DB interface {
	// IncrementRepairAttempts increments the repair attempt
	IncrementRepairAttempts(context.Context, *RemoteSegmentInfo) error
	// Get a irreparable's segment info from the db
	Get(context.Context, []byte) (*RemoteSegmentInfo, error)
	// Delete a irreparable's segment info from the db
	Delete(ctx context.Context, segmentPath []byte) error
}

// Database implements the irreparable services
type Database struct {
	db DB
}

// RemoteSegmentInfo is info about a single entry stored in the irreparable db
type RemoteSegmentInfo struct {
	EncryptedSegmentPath   []byte
	EncryptedSegmentDetail []byte //contains marshaled info of pb.Pointer
	LostPiecesCount        int64
	RepairUnixSec          int64
	RepairAttemptCount     int64
}

// NewServer creates instance of Server
func New(db DB) *Database {
	return &Database{
		db: db,
	}
}

// IncrementRepairAttempts a db entry for to increment the repair attempts field
func (db *Database) IncrementRepairAttempts(ctx context.Context, segmentInfo *RemoteSegmentInfo) (err error) {
	return db.IncrementRepairAttempts(ctx, segmentInfo)
}

// Get a irreparable's segment info from the db
func (db *Database) Get(ctx context.Context, segmentPath []byte) (resp *RemoteSegmentInfo, err error) {
	return db.Get(ctx, segmentPath)
}

// Delete a irreparable's segment info from the db
func (db *Database) Delete(ctx context.Context, segmentPath []byte) (err error) {
	return db.Delete(ctx, segmentPath)
}
