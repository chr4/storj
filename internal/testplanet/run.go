// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package testplanet

import (
	"strconv"
	"strings"
	"testing"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"storj.io/storj/internal/dbutil/pgutil"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/satellitedb"
	"storj.io/storj/satellite/satellitedb/satellitedbtest"
)

// Run runs testplanet in multiple configurations.
func Run(t *testing.T, config Config, test func(t *testing.T, ctx *testcontext.Context, planet *Planet)) {
	schemaSuffix := pgutil.CreateRandomTestingSchemaName(8)
	t.Log("schema-suffix ", schemaSuffix)

	for _, satelliteDB := range satellitedbtest.Databases() {
		satelliteDB := satelliteDB
		t.Run(satelliteDB.Name, func(t *testing.T) {
			t.Parallel()

			ctx := testcontext.New(t)
			defer ctx.Cleanup()

			if satelliteDB.URL == "" {
				t.Skipf("Database %s connection string not provided. %s", satelliteDB.Name, satelliteDB.Message)
			}

			planetConfig := config
			planetConfig.Reconfigure.NewBootstrapDB = nil
			planetConfig.Reconfigure.NewSatelliteDB = func(log *zap.Logger, index int) (satellite.DB, error) {
				schema := strings.ToLower(t.Name() + "-satellite/" + strconv.Itoa(index) + "-" + schemaSuffix)
				db, err := satellitedb.New(log, pgutil.ConnstrWithSchema(satelliteDB.URL, schema))
				if err != nil {
					t.Fatal(err)
				}

				err = db.CreateSchema(schema)
				if err != nil {
					t.Fatal(err)
				}

				return &satelliteSchema{
					DB:     db,
					schema: schema,
				}, nil
			}

			planet, err := NewCustom(zaptest.NewLogger(t), planetConfig)
			if err != nil {
				t.Fatal(err)
			}
			defer ctx.Check(planet.Shutdown)

			planet.Start(ctx)

			// make sure nodes are refreshed in db
			planet.Satellites[0].Discovery.Service.Refresh.TriggerWait()

			test(t, ctx, planet)
		})
	}
}

// satelliteSchema closes database and drops the associated schema
type satelliteSchema struct {
	satellite.DB
	schema string
}

func (db *satelliteSchema) Close() error {
	return errs.Combine(
		db.DB.DropSchema(db.schema),
		db.DB.Close(),
	)
}
