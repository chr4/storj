// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedb

import (
	"context"
	"database/sql"
	"time"

	"github.com/skyrings/skyring-common/tools/uuid"
	"github.com/zeebo/errs"

	"storj.io/storj/internal/memory"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/satellite/console"
	dbx "storj.io/storj/satellite/satellitedb/dbx"
)

// usagerollups implements console.UsageRollups
type usagerollups struct {
	db *dbx.DB
}

// GetProjectTotal retrieves project usage for a given period
func (db *usagerollups) GetProjectTotal(ctx context.Context, projectID uuid.UUID, since, before time.Time) (usage *console.ProjectUsage, err error) {
	defer mon.Task()(&ctx)(&err)
	since = timeTruncateDown(since)

	storageQuery := db.db.All_BucketStorageTally_By_ProjectId_And_BucketName_And_IntervalStart_GreaterOrEqual_And_IntervalStart_LessOrEqual_OrderBy_Desc_IntervalStart

	roullupsQuery := db.db.Rebind(`SELECT SUM(settled), SUM(inline), action
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND interval_start >= ? AND interval_start <= ?
		GROUP BY action`)

	rollupsRows, err := db.db.QueryContext(ctx, roullupsQuery, projectID[:], since, before)
	if err != nil {
		return nil, err
	}
	defer func() { err = errs.Combine(err, rollupsRows.Close()) }()

	var totalEgress int64
	for rollupsRows.Next() {
		var action pb.PieceAction
		var settled, inline int64

		err = rollupsRows.Scan(&settled, &inline, &action)
		if err != nil {
			return nil, err
		}

		// add values for egress
		if action == pb.PieceAction_GET || action == pb.PieceAction_GET_AUDIT || action == pb.PieceAction_GET_REPAIR {
			totalEgress += settled + inline
		}
	}

	buckets, err := db.getBuckets(ctx, projectID, since, before)
	if err != nil {
		return nil, err
	}

	bucketsTallies := make(map[string]*[]*dbx.BucketStorageTally)
	for _, bucket := range buckets {
		storageTallies, err := storageQuery(ctx,
			dbx.BucketStorageTally_ProjectId(projectID[:]),
			dbx.BucketStorageTally_BucketName([]byte(bucket)),
			dbx.BucketStorageTally_IntervalStart(since),
			dbx.BucketStorageTally_IntervalStart(before))

		if err != nil {
			return nil, err
		}

		bucketsTallies[bucket] = &storageTallies
	}

	usage = new(console.ProjectUsage)
	usage.Egress = memory.Size(totalEgress).GB()

	// sum up storage and objects
	for _, tallies := range bucketsTallies {
		for i := len(*tallies) - 1; i > 0; i-- {
			current := (*tallies)[i]

			hours := (*tallies)[i-1].IntervalStart.Sub(current.IntervalStart).Hours()

			usage.Storage += memory.Size(current.Inline).GB() * hours
			usage.Storage += memory.Size(current.Remote).GB() * hours
			usage.ObjectCount += float64(current.ObjectCount) * hours
		}
	}

	usage.Since = since
	usage.Before = before
	return usage, nil
}

// GetBucketUsageRollups retrieves summed usage rollups for every bucket of particular project for a given period
func (db *usagerollups) GetBucketUsageRollups(ctx context.Context, projectID uuid.UUID, since, before time.Time) (_ []console.BucketUsageRollup, err error) {
	defer mon.Task()(&ctx)(&err)
	since = timeTruncateDown(since)

	buckets, err := db.getBuckets(ctx, projectID, since, before)
	if err != nil {
		return nil, err
	}

	roullupsQuery := db.db.Rebind(`SELECT SUM(settled), SUM(inline), action
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND bucket_name = ? AND interval_start >= ? AND interval_start <= ?
		GROUP BY action`)

	storageQuery := db.db.All_BucketStorageTally_By_ProjectId_And_BucketName_And_IntervalStart_GreaterOrEqual_And_IntervalStart_LessOrEqual_OrderBy_Desc_IntervalStart

	var bucketUsageRollups []console.BucketUsageRollup
	for _, bucket := range buckets {
		bucketRollup := console.BucketUsageRollup{
			ProjectID:  projectID,
			BucketName: []byte(bucket),
			Since:      since,
			Before:     before,
		}

		// get bucket_bandwidth_rollups
		rollupsRows, err := db.db.QueryContext(ctx, roullupsQuery, projectID[:], []byte(bucket), since, before)
		if err != nil {
			return nil, err
		}
		defer func() { err = errs.Combine(err, rollupsRows.Close()) }()

		// fill egress
		for rollupsRows.Next() {
			var action pb.PieceAction
			var settled, inline int64

			err = rollupsRows.Scan(&settled, &inline, &action)
			if err != nil {
				return nil, err
			}

			switch action {
			case pb.PieceAction_GET:
				bucketRollup.GetEgress += memory.Size(settled + inline).GB()
			case pb.PieceAction_GET_AUDIT:
				bucketRollup.AuditEgress += memory.Size(settled + inline).GB()
			case pb.PieceAction_GET_REPAIR:
				bucketRollup.RepairEgress += memory.Size(settled + inline).GB()
			default:
				continue
			}
		}

		bucketStorageTallies, err := storageQuery(ctx,
			dbx.BucketStorageTally_ProjectId(projectID[:]),
			dbx.BucketStorageTally_BucketName([]byte(bucket)),
			dbx.BucketStorageTally_IntervalStart(since),
			dbx.BucketStorageTally_IntervalStart(before))

		if err != nil {
			return nil, err
		}

		// fill metadata, objects and stored data
		// hours calculated from previous tallies,
		// so we skip the most recent one
		for i := len(bucketStorageTallies) - 1; i > 0; i-- {
			current := bucketStorageTallies[i]

			hours := bucketStorageTallies[i-1].IntervalStart.Sub(current.IntervalStart).Hours()

			bucketRollup.RemoteStoredData += memory.Size(current.Remote).GB() * hours
			bucketRollup.InlineStoredData += memory.Size(current.Inline).GB() * hours
			bucketRollup.MetadataSize += memory.Size(current.MetadataSize).GB() * hours
			bucketRollup.RemoteSegments += float64(current.RemoteSegmentsCount) * hours
			bucketRollup.InlineSegments += float64(current.InlineSegmentsCount) * hours
			bucketRollup.ObjectCount += float64(current.ObjectCount) * hours
		}

		bucketUsageRollups = append(bucketUsageRollups, bucketRollup)
	}

	return bucketUsageRollups, nil
}

// GetBucketTotals retrieves bucket usage totals for period of time
func (db *usagerollups) GetBucketTotals(ctx context.Context, projectID uuid.UUID, cursor console.BucketUsageCursor, since, before time.Time) (_ *console.BucketUsagePage, err error) {
	defer mon.Task()(&ctx)(&err)
	since = timeTruncateDown(since)
	search := cursor.Search + "%"

	if cursor.Limit > 50 {
		cursor.Limit = 50
	}
	if cursor.Page == 0 {
		return nil, errs.New("page can not be 0")
	}

	page := &console.BucketUsagePage{
		Search: cursor.Search,
		Limit:  cursor.Limit,
		Offset: uint64((cursor.Page - 1) * cursor.Limit),
	}

	countQuery := db.db.Rebind(`SELECT COUNT(DISTINCT bucket_name)
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND interval_start >= ? AND interval_start <= ?
		AND CAST(bucket_name as TEXT) LIKE ?`)

	countRow := db.db.QueryRowContext(ctx,
		countQuery,
		projectID[:],
		since, before,
		search)

	err = countRow.Scan(&page.TotalCount)
	if err != nil {
		return nil, err
	}
	if page.TotalCount == 0 {
		return page, nil
	}
	if page.Offset > page.TotalCount-1 {
		return nil, errs.New("page is out of range")
	}

	bucketsQuery := db.db.Rebind(`SELECT DISTINCT bucket_name
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND interval_start >= ? AND interval_start <= ?
		AND CAST(bucket_name as TEXT) LIKE ?
		ORDER BY bucket_name ASC
		LIMIT ? OFFSET ?`)

	bucketRows, err := db.db.QueryContext(ctx,
		bucketsQuery,
		projectID[:],
		since, before,
		search,
		page.Limit,
		page.Offset)

	if err != nil {
		return nil, err
	}
	defer func() { err = errs.Combine(err, bucketRows.Close()) }()

	var buckets []string
	for bucketRows.Next() {
		var bucket string
		err = bucketRows.Scan(&bucket)
		if err != nil {
			return nil, err
		}

		buckets = append(buckets, bucket)
	}

	roullupsQuery := db.db.Rebind(`SELECT SUM(settled), SUM(inline), action
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND bucket_name = ? AND interval_start >= ? AND interval_start <= ?
		GROUP BY action`)

	storageQuery := db.db.Rebind(`SELECT inline, remote, object_count
		FROM bucket_storage_tallies
		WHERE project_id = ? AND bucket_name = ? AND interval_start >= ? AND interval_start <= ?
		ORDER BY interval_start DESC
		LIMIT 1`)

	var bucketUsages []console.BucketUsage
	for _, bucket := range buckets {
		bucketUsage := console.BucketUsage{
			ProjectID:  projectID,
			BucketName: bucket,
			Since:      since,
			Before:     before,
		}

		// get bucket_bandwidth_rollups
		rollupsRows, err := db.db.QueryContext(ctx, roullupsQuery, projectID[:], []byte(bucket), since, before)
		if err != nil {
			return nil, err
		}
		defer func() { err = errs.Combine(err, rollupsRows.Close()) }()

		var totalEgress int64
		for rollupsRows.Next() {
			var action pb.PieceAction
			var settled, inline int64

			err = rollupsRows.Scan(&settled, &inline, &action)
			if err != nil {
				return nil, err
			}

			// add values for egress
			if action == pb.PieceAction_GET || action == pb.PieceAction_GET_AUDIT || action == pb.PieceAction_GET_REPAIR {
				totalEgress += settled + inline
			}
		}

		bucketUsage.Egress = memory.Size(totalEgress).GB()

		storageRow := db.db.QueryRowContext(ctx, storageQuery, projectID[:], []byte(bucket), since, before)
		if err != nil {
			return nil, err
		}

		var inline, remote, objectCount int64
		err = storageRow.Scan(&inline, &remote, &objectCount)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		// fill storage and object count
		bucketUsage.Storage = memory.Size(inline + remote).GB()
		bucketUsage.ObjectCount = objectCount

		bucketUsages = append(bucketUsages, bucketUsage)
	}

	page.PageCount = uint(page.TotalCount / uint64(cursor.Limit))
	if page.TotalCount%uint64(cursor.Limit) != 0 {
		page.PageCount++
	}

	page.BucketUsages = bucketUsages
	page.CurrentPage = cursor.Page
	return page, nil
}

// getBuckets list all bucket of certain project for given period
func (db *usagerollups) getBuckets(ctx context.Context, projectID uuid.UUID, since, before time.Time) (_ []string, err error) {
	defer mon.Task()(&ctx)(&err)
	bucketsQuery := db.db.Rebind(`SELECT DISTINCT bucket_name
		FROM bucket_bandwidth_rollups
		WHERE project_id = ? AND interval_start >= ? AND interval_start <= ?`)

	bucketRows, err := db.db.QueryContext(ctx, bucketsQuery, projectID[:], since, before)
	if err != nil {
		return nil, err
	}
	defer func() { err = errs.Combine(err, bucketRows.Close()) }()

	var buckets []string
	for bucketRows.Next() {
		var bucket string
		err = bucketRows.Scan(&bucket)
		if err != nil {
			return nil, err
		}

		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

// timeTruncateDown truncates down to the hour before to be in sync with orders endpoint
func timeTruncateDown(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}
