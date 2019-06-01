// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/testcontext"
	"storj.io/storj/lib/uplink"
	"storj.io/storj/pkg/storj"
)

// TODO: Start up test planet and call these from bash instead
func TestCBucketTests(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	planet := startTestPlanet(t, ctx)
	defer ctx.Check(planet.Shutdown)

	consoleProject := newProject(t, planet)
	consoleApikey := newAPIKey(t, ctx, planet, consoleProject.ID)
	satelliteAddr := planet.Satellites[0].Addr()
	bucketName := "TestBucket"

	envVars := []string{
		"SATELLITE_ADDR=" + satelliteAddr,
		"APIKEY=" + consoleApikey,
		"BUCKET_NAME=" + bucketName,
	}

	goUplink, err := uplink.NewUplink(ctx, testConfig)
	require.NoError(t, err)

	apikey, err := uplink.ParseAPIKey(consoleApikey)
	require.NoError(t, err)

	project, err := goUplink.OpenProject(ctx, satelliteAddr, apikey, nil)
	require.NoError(t, err)

	_, err = project.CreateBucket(ctx, bucketName, nil)
	require.NoError(t, err)

	key := storj.Key{}
	copy(key[:], []byte("abcdefghijklmnopqrstuvwxyzABCDEF"))
	bucket, err := project.OpenBucket(ctx, bucketName, nil)
	require.NoError(t, err)

	{
		runCTest(t, ctx, "bucket_test.c", envVars...)

		objectList, err := bucket.ListObjects(ctx, nil)
		require.NoError(t, err)

		require.Len(t, objectList.Items, 4)
		object, err := bucket.OpenObject(ctx, objectList.Items[0].Path)
		require.NoError(t, err)

		assert.Condition(t, func() bool {
			return time.Now().Sub(object.Meta.Modified).Seconds() < 5
		})
		// TODO: add more assertions
	}
}
