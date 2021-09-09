// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prysmgrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	client "github.com/attestantio/go-eth2-client"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Service is an Ethereum 2 client service.
type Service struct {
	// Hold the initialising context to allow for streams to use it.
	ctx context.Context

	// Client connection.
	conn    *grpc.ClientConn
	address string
	timeout time.Duration

	maxPageSize int32

	// Various information from the node that never changes once we have it.
	spec                          map[string]interface{}
	genesisTime                   *time.Time
	genesisValidatorsRoot         []byte
	slotDuration                  *time.Duration
	slotsPerEpoch                 *uint64
	farFutureEpoch                *spec.Epoch
	targetAggregatorsPerCommittee *uint64
	genesisForkVersion            []byte

	// Event handlers.
	beaconChainHeadUpdatedMutex    sync.RWMutex
	beaconChainHeadUpdatedHandlers []client.BeaconChainHeadUpdatedHandler

	// The standard API commonly uses validator indices, and the prysm API commonly uses public keys.
	// We keep a mapping of index to public keys to avoid repeated lookups.
	indexMap   map[spec.ValidatorIndex]spec.BLSPubKey
	indexMapMu sync.RWMutex
}

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new Ethereum 2 client service, connecting with Prysm GRPC.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "prysmgrpc").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	grpcOpts := []grpc.DialOption{
		// Maximum receive value 256 MB
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(256 * 1024 * 1024)),
	}

	if parameters.tls {
		credentials, err := tlsCredentials(ctx, parameters.clientCert, parameters.clientKey, parameters.caCert)
		if err != nil {
			return nil, errors.Wrap(err, "problem with TLS credentials")
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	dialCtx, cancel := context.WithTimeout(ctx, parameters.timeout)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, parameters.address, grpcOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial connection")
	}

	s := &Service{
		ctx:         ctx,
		conn:        conn,
		address:     parameters.address,
		timeout:     parameters.timeout,
		maxPageSize: 250, // Prysm default.
		indexMap:    make(map[spec.ValidatorIndex]spec.BLSPubKey),
	}

	// Obtain the node version to confirm the connection is good.
	if _, err := s.NodeVersion(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to confirm node connection")
	}

	// Obtain the page size.
	if maxPageSize, err := s.obtainMaxPageSize(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to obtain largest page size")
	} else {
		s.maxPageSize = maxPageSize
		log.Trace().Int32("max_page_size", maxPageSize).Msg("Set maximum page size")
	}

	// Close the service on context done.
	go func(s *Service) {
		<-ctx.Done()
		log.Trace().Msg("Context done; closing connection")
		s.close()
	}(s)

	return s, nil
}

// Name provides the name of the service.
func (s *Service) Name() string {
	return "Prysm (gRPC)"
}

// Address provides the address for the connection.
func (s *Service) Address() string {
	return s.address
}

// Close the service, freeing up resources.
func (s *Service) close() {
	if err := s.conn.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close connection")
	}
}

func (s *Service) obtainMaxPageSize(ctx context.Context) (int32, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)
	if conn == nil {
		return -1, errors.New("failed to obtain beacon chain client")
	}

	validatorsReq := &ethpb.ListValidatorsRequest{
		PageSize: 9999999,
	}

	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	_, err := conn.ListValidators(opCtx, validatorsReq)
	cancel()
	if err == nil {
		// Max page size is > 9999999, that'll do for us
		return 9999999, nil
	}
	if !strings.Contains(err.Error(), "Requested page size 9999999 can not be greater than max size ") {
		return -1, errors.New("failed to obtain message with max size")
	}

	re := regexp.MustCompile(`^.*Requested page size 9999999 can not be greater than max size ([0-9]+)`)
	res := re.FindStringSubmatch(err.Error())
	if len(res) != 2 {
		return -1, errors.New("unexpected error response; cannot parse for max page size")
	}
	maxPageSize, err := strconv.ParseInt(res[1], 10, 32)
	if err != nil {
		return -1, errors.New("invalid value for max page size")
	}
	return int32(maxPageSize), nil
}

// tlsCredentials composes a set of transport credentials given optional client and CA certificates.
func tlsCredentials(ctx context.Context, clientCert []byte, clientKey []byte, caCert []byte) (credentials.TransportCredentials, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	if clientCert != nil && clientKey != nil {
		clientPair, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load client keypair")
		}
		tlsCfg.Certificates = []tls.Certificate{clientPair}
	}

	if caCert != nil {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(caCert) {
			return nil, errors.New("failed to add CA certificate")
		}
		tlsCfg.RootCAs = cp
	}

	return credentials.NewTLS(tlsCfg), nil
}
