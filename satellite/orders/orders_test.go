// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package orders_test

import (
	"testing"
	"time"

	"github.com/skyrings/skyring-common/tools/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/memory"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/internal/testrand"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/satellite/orders"
	"storj.io/storj/uplink"
)

func TestSendingReceivingOrders(t *testing.T) {
	// test happy path
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 6, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		planet.Satellites[0].Audit.Service.Loop.Stop()
		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.Pause()
		}

		expectedData := testrand.Bytes(50 * memory.KiB)

		redundancy := noLongTailRedundancy(planet)
		err := planet.Uplinks[0].UploadWithConfig(ctx, planet.Satellites[0], &redundancy, "testbucket", "test/path", expectedData)
		require.NoError(t, err)

		sumBeforeSend := 0
		for _, storageNode := range planet.StorageNodes {
			infos, err := storageNode.DB.Orders().ListUnsent(ctx, 10)
			require.NoError(t, err)
			sumBeforeSend += len(infos)
		}
		require.NotZero(t, sumBeforeSend)

		sumUnsent := 0
		sumArchived := 0

		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.TriggerWait()

			infos, err := storageNode.DB.Orders().ListUnsent(ctx, 10)
			require.NoError(t, err)
			sumUnsent += len(infos)

			archivedInfos, err := storageNode.DB.Orders().ListArchived(ctx, sumBeforeSend)
			require.NoError(t, err)
			sumArchived += len(archivedInfos)
		}

		require.Zero(t, sumUnsent)
		require.Equal(t, sumBeforeSend, sumArchived)
	})
}

func TestUnableToSendOrders(t *testing.T) {
	// test sending when satellite is unavailable
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 6, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		planet.Satellites[0].Audit.Service.Loop.Stop()
		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.Pause()
		}

		expectedData := testrand.Bytes(50 * memory.KiB)

		redundancy := noLongTailRedundancy(planet)
		err := planet.Uplinks[0].UploadWithConfig(ctx, planet.Satellites[0], &redundancy, "testbucket", "test/path", expectedData)
		require.NoError(t, err)

		sumBeforeSend := 0
		for _, storageNode := range planet.StorageNodes {
			infos, err := storageNode.DB.Orders().ListUnsent(ctx, 10)
			require.NoError(t, err)
			sumBeforeSend += len(infos)
		}
		require.NotZero(t, sumBeforeSend)

		err = planet.StopPeer(planet.Satellites[0])
		require.NoError(t, err)

		sumUnsent := 0
		sumArchived := 0
		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.TriggerWait()

			infos, err := storageNode.DB.Orders().ListUnsent(ctx, 10)
			require.NoError(t, err)
			sumUnsent += len(infos)

			archivedInfos, err := storageNode.DB.Orders().ListArchived(ctx, sumBeforeSend)
			require.NoError(t, err)
			sumArchived += len(archivedInfos)
		}

		require.Zero(t, sumArchived)
		require.Equal(t, sumBeforeSend, sumUnsent)
	})
}

func TestUploadDownloadBandwidth(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 6, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		hourBeforeTest := time.Now().UTC().Add(-time.Hour)

		planet.Satellites[0].Audit.Service.Loop.Stop()
		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.Pause()
		}

		expectedData := testrand.Bytes(50 * memory.KiB)

		redundancy := noLongTailRedundancy(planet)
		err := planet.Uplinks[0].UploadWithConfig(ctx, planet.Satellites[0], &redundancy, "testbucket", "test/path", expectedData)
		require.NoError(t, err)

		data, err := planet.Uplinks[0].Download(ctx, planet.Satellites[0], "testbucket", "test/path")
		require.NoError(t, err)
		require.Equal(t, expectedData, data)

		var expectedBucketBandwidth int64
		expectedStorageBandwidth := make(map[storj.NodeID]int64)
		for _, storageNode := range planet.StorageNodes {
			infos, err := storageNode.DB.Orders().ListUnsent(ctx, 10)
			require.NoError(t, err)
			if len(infos) > 0 {
				for _, info := range infos {
					expectedBucketBandwidth += info.Order.Amount
					expectedStorageBandwidth[storageNode.ID()] += info.Order.Amount
				}
			}
		}

		for _, storageNode := range planet.StorageNodes {
			storageNode.Storage2.Sender.Loop.TriggerWait()
		}

		projects, err := planet.Satellites[0].DB.Console().Projects().GetAll(ctx)
		require.NoError(t, err)

		ordersDB := planet.Satellites[0].DB.Orders()

		bucketBandwidth, err := ordersDB.GetBucketBandwidth(ctx, projects[0].ID, []byte("testbucket"), hourBeforeTest, time.Now().UTC())
		require.NoError(t, err)
		require.Equal(t, expectedBucketBandwidth, bucketBandwidth)

		for _, storageNode := range planet.StorageNodes {
			nodeBandwidth, err := ordersDB.GetStorageNodeBandwidth(ctx, storageNode.ID(), hourBeforeTest, time.Now().UTC())
			require.NoError(t, err)
			require.Equal(t, expectedStorageBandwidth[storageNode.ID()], nodeBandwidth)
		}
	})
}

func noLongTailRedundancy(planet *testplanet.Planet) uplink.RSConfig {
	redundancy := planet.Uplinks[0].GetConfig(planet.Satellites[0]).RS
	redundancy.SuccessThreshold = redundancy.MaxThreshold
	return redundancy
}

func TestSplitBucketIDInvalid(t *testing.T) {
	var testCases = []struct {
		name     string
		bucketID []byte
	}{
		{"invalid, not valid UUID", []byte("not UUID string/bucket1")},
		{"invalid, not valid UUID, no bucket", []byte("not UUID string")},
		{"invalid, no project, no bucket", []byte("")},
	}
	for _, tt := range testCases {
		tt := tt // avoid scopelint error, ref: https://github.com/golangci/golangci-lint/issues/281
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := orders.SplitBucketID(tt.bucketID)
			assert.NotNil(t, err)
			assert.Error(t, err)
		})
	}
}

func TestSplitBucketIDValid(t *testing.T) {
	var testCases = []struct {
		name               string
		project            string
		bucketName         string
		expectedBucketName string
	}{
		{"valid, no bucket, no objects", "bb6218e3-4b4a-4819-abbb-fa68538e33c0", "", ""},
		{"valid, with bucket", "bb6218e3-4b4a-4819-abbb-fa68538e33c0", "testbucket", "testbucket"},
		{"valid, with object", "bb6218e3-4b4a-4819-abbb-fa68538e33c0", "testbucket/foo/bar.txt", "testbucket"},
	}
	for _, tt := range testCases {
		tt := tt // avoid scopelint error, ref: https://github.com/golangci/golangci-lint/issues/281
		t.Run(tt.name, func(t *testing.T) {
			expectedProjectID, err := uuid.Parse(tt.project)
			assert.NoError(t, err)
			bucketID := expectedProjectID.String() + "/" + tt.bucketName

			actualProjectID, actualBucketName, err := orders.SplitBucketID([]byte(bucketID))
			assert.NoError(t, err)
			assert.Equal(t, actualProjectID, expectedProjectID)
			assert.Equal(t, actualBucketName, []byte(tt.expectedBucketName))
		})
	}
}
