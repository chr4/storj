// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package uplink

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	"github.com/vivint/infectious"
	"github.com/zeebo/errs"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/internal/memory"
	"storj.io/storj/pkg/eestream"
	"storj.io/storj/pkg/encryption"
	"storj.io/storj/pkg/identity"
	"storj.io/storj/pkg/metainfo/kvmetainfo"
	"storj.io/storj/pkg/peertls/tlsopts"
	ecclient "storj.io/storj/pkg/storage/ec"
	"storj.io/storj/pkg/storage/segments"
	"storj.io/storj/pkg/storage/streams"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/pkg/transport"
	"storj.io/storj/uplink/metainfo"
)

// RSConfig is a configuration struct that keeps details about default
// redundancy strategy information
type RSConfig struct {
	MaxBufferMem     memory.Size `help:"maximum buffer memory (in bytes) to be allocated for read buffers" default:"4MiB" hidden:"true"`
	ErasureShareSize memory.Size `help:"the size of each new erasure share in bytes" default:"256B" hidden:"true"`
	MinThreshold     int         `help:"the minimum pieces required to recover a segment. k." releaseDefault:"29" devDefault:"4" hidden:"true"`
	RepairThreshold  int         `help:"the minimum safe pieces before a repair is triggered. m." releaseDefault:"35" devDefault:"6" hidden:"true"`
	SuccessThreshold int         `help:"the desired total pieces for a segment. o." releaseDefault:"80" devDefault:"8" hidden:"true"`
	MaxThreshold     int         `help:"the largest amount of pieces to encode to. n." releaseDefault:"130" devDefault:"10" hidden:"true"`
}

// EncryptionConfig is a configuration struct that keeps details about
// encrypting segments
type EncryptionConfig struct {
	EncryptionKey string `help:"the root key for encrypting the data which will be stored in KeyFilePath" setup:"true"`
	KeyFilepath   string `help:"the path to the file which contains the root key for encrypting the data"`
	DataType      int    `help:"Type of encryption to use for content and metadata (1=AES-GCM, 2=SecretBox)" default:"1"`
	PathType      int    `help:"Type of encryption to use for paths (0=Unencrypted, 1=AES-GCM, 2=SecretBox)" default:"1"`
}

// ClientConfig is a configuration struct for the uplink that controls how
// to talk to the rest of the network.
type ClientConfig struct {
	APIKey         string        `default:"" help:"the api key to use for the satellite" noprefix:"true"`
	SatelliteAddr  string        `releaseDefault:"127.0.0.1:7777" devDefault:"127.0.0.1:10000" help:"the address to use for the satellite" noprefix:"true"`
	MaxInlineSize  memory.Size   `help:"max inline segment size in bytes" default:"4KiB"`
	SegmentSize    memory.Size   `help:"the size of a segment in bytes" default:"64MiB"`
	RequestTimeout time.Duration `help:"timeout for request" default:"0h0m20s"`
	DialTimeout    time.Duration `help:"timeout for dials" default:"0h0m20s"`
}

// Config uplink configuration
type Config struct {
	Client ClientConfig
	RS     RSConfig
	Enc    EncryptionConfig
	TLS    tlsopts.Config
}

var (
	mon = monkit.Package()

	// Error is the errs class of standard End User Client errors
	Error = errs.Class("Uplink configuration error")
)

// GetMetainfo returns an implementation of storj.Metainfo
func (c Config) GetMetainfo(ctx context.Context, identity *identity.FullIdentity) (db storj.Metainfo, ss streams.Store, err error) {
	defer mon.Task()(&ctx)(&err)

	tlsOpts, err := tlsopts.NewOptions(identity, c.TLS)
	if err != nil {
		return nil, nil, err
	}

	// ToDo: Handle Versioning for Uplinks here

	tc := transport.NewClientWithTimeouts(tlsOpts, transport.Timeouts{
		Request: c.Client.RequestTimeout,
		Dial:    c.Client.DialTimeout,
	})

	if c.Client.SatelliteAddr == "" {
		return nil, nil, errors.New("satellite address not specified")
	}

	metainfo, err := metainfo.NewClient(ctx, tc, c.Client.SatelliteAddr, c.Client.APIKey)
	if err != nil {
		return nil, nil, Error.New("failed to connect to metainfo service: %v", err)
	}

	ec := ecclient.NewClient(tc, c.RS.MaxBufferMem.Int())
	fc, err := infectious.NewFEC(c.RS.MinThreshold, c.RS.MaxThreshold)
	if err != nil {
		return nil, nil, Error.New("failed to create erasure coding client: %v", err)
	}
	rs, err := eestream.NewRedundancyStrategy(eestream.NewRSScheme(fc, c.RS.ErasureShareSize.Int()), c.RS.RepairThreshold, c.RS.SuccessThreshold)
	if err != nil {
		return nil, nil, Error.New("failed to create redundancy strategy: %v", err)
	}

	maxEncryptedSegmentSize, err := encryption.CalcEncryptedSize(c.Client.SegmentSize.Int64(), c.GetEncryptionScheme())
	if err != nil {
		return nil, nil, Error.New("failed to calculate max encrypted segment size: %v", err)
	}
	segments := segments.NewSegmentStore(metainfo, ec, rs, c.Client.MaxInlineSize.Int(), maxEncryptedSegmentSize)

	blockSize := c.GetEncryptionScheme().BlockSize
	if int(blockSize)%c.RS.ErasureShareSize.Int()*c.RS.MinThreshold != 0 {
		err = Error.New("EncryptionBlockSize must be a multiple of ErasureShareSize * RS MinThreshold")
		return nil, nil, err
	}

	key, err := LoadEncryptionKey(c.Enc.KeyFilepath)
	if err != nil {
		return nil, nil, Error.Wrap(err)
	}

	streams, err := streams.NewStreamStore(segments, c.Client.SegmentSize.Int64(), key,
		int(blockSize), storj.Cipher(c.Enc.DataType), c.Client.MaxInlineSize.Int(),
	)
	if err != nil {
		return nil, nil, Error.New("failed to create stream store: %v", err)
	}

	return kvmetainfo.New(metainfo, streams, segments, key, blockSize, rs, c.Client.SegmentSize.Int64()), streams, nil
}

// GetRedundancyScheme returns the configured redundancy scheme for new uploads
func (c Config) GetRedundancyScheme() storj.RedundancyScheme {
	return storj.RedundancyScheme{
		Algorithm:      storj.ReedSolomon,
		ShareSize:      c.RS.ErasureShareSize.Int32(),
		RequiredShares: int16(c.RS.MinThreshold),
		RepairShares:   int16(c.RS.RepairThreshold),
		OptimalShares:  int16(c.RS.SuccessThreshold),
		TotalShares:    int16(c.RS.MaxThreshold),
	}
}

// GetPathCipherSuite returns the cipher suite used for path encryption for bucket objects
func (c Config) GetPathCipherSuite() storj.CipherSuite {
	return storj.Cipher(c.Enc.PathType).ToCipherSuite()
}

// GetEncryptionScheme returns the configured encryption scheme for new uploads
// Blocksize should align with the stripe size therefore multiples of stripes
// should fit in every encryption block. Instead of lettings users configure this
// multiple value, we hardcode stripesPerBlock as 2 for simplicity.
func (c Config) GetEncryptionScheme() storj.EncryptionScheme {
	const stripesPerBlock = 2
	return storj.EncryptionScheme{
		Cipher:    storj.Cipher(c.Enc.DataType),
		BlockSize: c.GetRedundancyScheme().StripeSize() * stripesPerBlock,
	}
}

// GetSegmentSize returns the segment size set in uplink config
func (c Config) GetSegmentSize() memory.Size {
	return c.Client.SegmentSize
}

// LoadEncryptionKey loads the encryption key stored in the file pointed by
// filepath.
//
// An error is file is not found or there is an I/O error.
func LoadEncryptionKey(filepath string) (key *storj.Key, error error) {
	if filepath == "" {
		return &storj.Key{}, nil
	}

	rawKey, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return storj.NewKey(rawKey)
}
