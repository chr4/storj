// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/currency"
	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testrand"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/console"
	"storj.io/storj/satellite/rewards"
	"storj.io/storj/satellite/satellitedb/satellitedbtest"
)

func TestUsercredits(t *testing.T) {
	satellitedbtest.Run(t, func(t *testing.T, db satellite.DB) {
		ctx := testcontext.New(t)
		defer ctx.Cleanup()

		consoleDB := db.Console()

		user, referrer, offer := setupData(ctx, t, db)
		randomID := testrand.UUID()

		// test foreign key constraint for inserting a new user credit entry with randomID
		var invalidUserCredits = []console.UserCredit{
			{
				UserID:        randomID,
				OfferID:       offer.ID,
				ReferredBy:    referrer.ID,
				CreditsEarned: currency.Cents(100),
				ExpiresAt:     time.Now().UTC().AddDate(0, 1, 0),
			},
			{
				UserID:        user.ID,
				OfferID:       10,
				ReferredBy:    referrer.ID,
				CreditsEarned: currency.Cents(100),
				ExpiresAt:     time.Now().UTC().AddDate(0, 1, 0),
			},
			{
				UserID:        user.ID,
				OfferID:       offer.ID,
				ReferredBy:    randomID,
				CreditsEarned: currency.Cents(100),
				ExpiresAt:     time.Now().UTC().AddDate(0, 1, 0),
			},
		}

		for _, ivc := range invalidUserCredits {
			_, err := consoleDB.UserCredits().Create(ctx, ivc)
			require.Error(t, err)
		}

		type result struct {
			remainingCharge int
			usage           console.UserCreditUsage
			hasErr          bool
		}

		var validUserCredits = []struct {
			userCredit     console.UserCredit
			chargedCredits int
			expected       result
		}{
			{
				userCredit: console.UserCredit{
					UserID:        user.ID,
					OfferID:       offer.ID,
					ReferredBy:    referrer.ID,
					CreditsEarned: currency.Cents(100),
					ExpiresAt:     time.Now().UTC().AddDate(0, 1, 0),
				},
				chargedCredits: 120,
				expected: result{
					remainingCharge: 20,
					usage: console.UserCreditUsage{
						AvailableCredits: currency.Cents(0),
						UsedCredits:      currency.Cents(100),
						Referred:         0,
					},
					hasErr: false,
				},
			},
			{
				// simulate a credit that's already expired
				userCredit: console.UserCredit{
					UserID:        user.ID,
					OfferID:       offer.ID,
					ReferredBy:    referrer.ID,
					CreditsEarned: currency.Cents(100),
					ExpiresAt:     time.Now().UTC().AddDate(0, 0, -5),
				},
				chargedCredits: 60,
				expected: result{
					remainingCharge: 60,
					usage: console.UserCreditUsage{
						AvailableCredits: currency.Cents(0),
						UsedCredits:      currency.Cents(100),
						Referred:         0,
					},
					hasErr: true,
				},
			},
			{
				// simulate a credit that's not expired
				userCredit: console.UserCredit{
					UserID:        user.ID,
					OfferID:       offer.ID,
					ReferredBy:    referrer.ID,
					CreditsEarned: currency.Cents(100),
					ExpiresAt:     time.Now().UTC().AddDate(0, 0, 5),
				},
				chargedCredits: 80,
				expected: result{
					remainingCharge: 0,
					usage: console.UserCreditUsage{
						AvailableCredits: currency.Cents(20),
						UsedCredits:      currency.Cents(180),
						Referred:         0,
					},
					hasErr: false,
				},
			},
		}

		for i, vc := range validUserCredits {
			_, err := consoleDB.UserCredits().Create(ctx, vc.userCredit)
			require.NoError(t, err)

			{
				remainingCharge, err := consoleDB.UserCredits().UpdateAvailableCredits(ctx, vc.chargedCredits, vc.userCredit.UserID, time.Now().UTC())
				if vc.expected.hasErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, vc.expected.remainingCharge, remainingCharge)
			}

			{
				usage, err := consoleDB.UserCredits().GetCreditUsage(ctx, vc.userCredit.UserID, time.Now().UTC())
				require.NoError(t, err)
				require.Equal(t, vc.expected.usage, *usage)
			}

			{
				referred, err := consoleDB.UserCredits().GetCreditUsage(ctx, referrer.ID, time.Now().UTC())
				require.NoError(t, err)
				require.Equal(t, int64(i+1), referred.Referred)
			}
		}
	})
}

func setupData(ctx context.Context, t *testing.T, db satellite.DB) (user *console.User, referrer *console.User, offer *rewards.Offer) {
	consoleDB := db.Console()
	offersDB := db.Rewards()

	// create user
	userPassHash := testrand.Bytes(8)
	referrerPassHash := testrand.Bytes(8)

	var err error

	// create an user
	user, err = consoleDB.Users().Insert(ctx, &console.User{
		FullName:     "John Doe",
		Email:        "john@mail.test",
		PasswordHash: userPassHash,
		Status:       console.Active,
	})
	require.NoError(t, err)

	//create an user as referrer
	referrer, err = consoleDB.Users().Insert(ctx, &console.User{
		FullName:     "referrer",
		Email:        "referrer@mail.test",
		PasswordHash: referrerPassHash,
		Status:       console.Active,
	})
	require.NoError(t, err)

	// create offer
	offer, err = offersDB.Create(ctx, &rewards.NewOffer{
		Name:                      "test",
		Description:               "test offer 1",
		AwardCredit:               currency.Cents(100),
		InviteeCredit:             currency.Cents(50),
		AwardCreditDurationDays:   60,
		InviteeCreditDurationDays: 30,
		RedeemableCap:             50,
		ExpiresAt:                 time.Now().UTC().Add(time.Hour * 1),
		Status:                    rewards.Active,
		Type:                      rewards.Referral,
	})
	require.NoError(t, err)

	return user, referrer, offer
}
