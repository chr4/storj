// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedb

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/skyrings/skyring-common/tools/uuid"
	"github.com/zeebo/errs"

	"storj.io/storj/internal/currency"
	"storj.io/storj/satellite/console"
	dbx "storj.io/storj/satellite/satellitedb/dbx"
)

type usercredits struct {
	db *dbx.DB
}

// GetCreditUsage returns the total amount of referral a user has made based on user id, total available credits, and total used credits based on user id
func (c *usercredits) GetCreditUsage(ctx context.Context, userID uuid.UUID, expirationEndDate time.Time) (*console.UserCreditUsage, error) {
	usageRows, err := c.db.DB.QueryContext(ctx, c.db.Rebind(`SELECT a.used_credit, b.available_credit, c.referred
		FROM (SELECT SUM(credits_used_in_cents) AS used_credit FROM user_credits WHERE user_id = ?) AS a,
		(SELECT SUM(credits_earned_in_cents - credits_used_in_cents) AS available_credit FROM user_credits WHERE expires_at > ? AND user_id = ?) AS b,
		(SELECT count(id) AS referred FROM user_credits WHERE user_credits.referred_by = ?) AS c;`), userID[:], expirationEndDate, userID[:], userID[:])
	if err != nil {
		return nil, errs.Wrap(err)
	}

	usage := console.UserCreditUsage{}

	for usageRows.Next() {

		var (
			usedCreditInCents      sql.NullInt64
			availableCreditInCents sql.NullInt64
			referred               sql.NullInt64
		)
		err = usageRows.Scan(&usedCreditInCents, &availableCreditInCents, &referred)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		usage.Referred += referred.Int64
		usage.UsedCredits = usage.UsedCredits.Add(currency.Cents(int(usedCreditInCents.Int64)))
		usage.AvailableCredits = usage.AvailableCredits.Add(currency.Cents(int(availableCreditInCents.Int64)))
	}

	return &usage, nil
}

// Create insert a new record of user credit
func (c *usercredits) Create(ctx context.Context, userCredit console.UserCredit) (*console.UserCredit, error) {
	credit, err := c.db.Create_UserCredit(ctx,
		dbx.UserCredit_UserId(userCredit.UserID[:]),
		dbx.UserCredit_OfferId(userCredit.OfferID),
		dbx.UserCredit_CreditsEarnedInCents(userCredit.CreditsEarned.Cents()),
		dbx.UserCredit_ExpiresAt(userCredit.ExpiresAt),
		dbx.UserCredit_Create_Fields{
			ReferredBy: dbx.UserCredit_ReferredBy(userCredit.ReferredBy[:]),
		},
	)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return convertDBCredit(credit)
}

// UpdateAvailableCredits updates user's available credits based on their spending and the time of their spending
func (c *usercredits) UpdateAvailableCredits(ctx context.Context, creditsToCharge int, id uuid.UUID, expirationEndDate time.Time) (remainingCharge int, err error) {
	tx, err := c.db.Open(ctx)
	if err != nil {
		return creditsToCharge, errs.Wrap(err)
	}

	availableCredits, err := tx.All_UserCredit_By_UserId_And_ExpiresAt_Greater_And_CreditsUsedInCents_Less_CreditsEarnedInCents_OrderBy_Asc_ExpiresAt(ctx,
		dbx.UserCredit_UserId(id[:]),
		dbx.UserCredit_ExpiresAt(expirationEndDate),
	)
	if err != nil {
		return creditsToCharge, errs.Wrap(errs.Combine(err, tx.Rollback()))
	}
	if len(availableCredits) == 0 {
		return creditsToCharge, errs.Combine(errs.New("No available credits"), tx.Commit())
	}

	values := make([]interface{}, len(availableCredits)*2)
	rowIds := make([]interface{}, len(availableCredits))

	remainingCharge = creditsToCharge
	for i, credit := range availableCredits {
		if remainingCharge == 0 {
			break
		}

		creditsForUpdateInCents := credit.CreditsEarnedInCents - credit.CreditsUsedInCents

		if remainingCharge < creditsForUpdateInCents {
			creditsForUpdateInCents = remainingCharge
		}

		values[i%2] = credit.Id
		values[(i%2 + 1)] = creditsForUpdateInCents
		rowIds[i] = credit.Id

		remainingCharge -= creditsForUpdateInCents
	}

	values = append(values, rowIds...)

	var statement string
	switch t := c.db.Driver().(type) {
	case *sqlite3.SQLiteDriver:
		statement = generateQuery(len(availableCredits), false)
	case *pq.Driver:
		statement = generateQuery(len(availableCredits), true)
	default:
		return creditsToCharge, errs.New("Unsupported database %t", t)
	}

	_, err = tx.Tx.ExecContext(ctx, c.db.Rebind(`UPDATE user_credits SET
			credits_used_in_cents = CASE `+statement), values...)
	if err != nil {
		return creditsToCharge, errs.Wrap(errs.Combine(err, tx.Rollback()))
	}
	return remainingCharge, errs.Wrap(tx.Commit())
}

func generateQuery(totalRows int, toInt bool) (query string) {
	whereClause := `WHERE id IN (`
	condition := `WHEN id=? THEN ? `
	if toInt {
		condition = `WHEN id=? THEN ?::int `
	}

	for i := 0; i < totalRows; i++ {
		query += condition

		if i == totalRows-1 {
			query += ` END ` + whereClause + ` ?);`
			break
		}
		whereClause += `?, `
	}

	return query
}

func convertDBCredit(userCreditDBX *dbx.UserCredit) (*console.UserCredit, error) {
	if userCreditDBX == nil {
		return nil, errs.New("userCreditDBX parameter is nil")
	}

	userID, err := bytesToUUID(userCreditDBX.UserId)
	if err != nil {
		return nil, err
	}

	referredByID, err := bytesToUUID(userCreditDBX.ReferredBy)
	if err != nil {
		return nil, err
	}

	return &console.UserCredit{
		ID:            userCreditDBX.Id,
		UserID:        userID,
		OfferID:       userCreditDBX.OfferId,
		ReferredBy:    referredByID,
		CreditsEarned: currency.Cents(userCreditDBX.CreditsEarnedInCents),
		CreditsUsed:   currency.Cents(userCreditDBX.CreditsUsedInCents),
		ExpiresAt:     userCreditDBX.ExpiresAt,
		CreatedAt:     userCreditDBX.CreatedAt,
	}, nil
}
