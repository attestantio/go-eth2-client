// Copyright © 2020 - 2026 Attestant Limited.
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

package http_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

// timeout for tests.
var timeout = 5 * time.Minute

// Global HTTP service instance shared across all tests to reduce connection overhead.
var globalHTTPService any

// testCoordinator controls how many tests can run concurrently to avoid overwhelming the endpoint.
// This is configured via HTTP_TEST_CONCURRENCY (default: 1 for sequential execution).
var testCoordinator *semaphore.Weighted

func TestMain(m *testing.M) {
	if logLevel := os.Getenv("HTTP_DEBUG_LOG_ENABLED"); strings.ToLower(logLevel) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	// Initialize test coordinator to limit concurrent test execution.
	// Default is 1 (sequential) to avoid overwhelming the endpoint.
	// Set HTTP_TEST_CONCURRENCY to allow more concurrent tests.
	concurrency := int64(1) // default: run tests sequentially
	if concurrencyStr := os.Getenv("HTTP_TEST_CONCURRENCY"); concurrencyStr != "" {
		if parsed, err := strconv.ParseInt(concurrencyStr, 10, 64); err == nil && parsed > 0 {
			concurrency = parsed
		}
	}
	testCoordinator = semaphore.NewWeighted(concurrency)

	// On the validation devnet a fork-gated test that skips is a bug rather than a
	// correct outcome, so this turns those skips into failures. See requireGloas.
	requireGloas = strings.EqualFold(os.Getenv("HTTP_REQUIRE_GLOAS"), "true")

	if os.Getenv("HTTP_ADDRESS") != "" {
		// Initialize global HTTP service for all tests to share
		initGlobalHTTPService()
		os.Exit(m.Run())
	}
}

// newTestService creates an HTTP service against HTTP_ADDRESS, authenticating with
// HTTP_BEARER_TOKEN when one is set. customSpecSupport picks which branch the SSZ
// decoders take, so a caller that needs a mode other than the suite's default gets it
// without restating the address, timeout and token handling.
func newTestService(ctx context.Context,
	customSpecSupport bool,
	params ...http.Parameter,
) (
	client.Service,
	error,
) {
	// Custom spec support is required for the interim Glamsterdam devnet's minimal
	// preset; it is a no-op superset for mainnet-preset endpoints (see ADR-0003).
	parameters := []http.Parameter{
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
		http.WithCustomSpecSupport(customSpecSupport),
	}

	if token := os.Getenv("HTTP_BEARER_TOKEN"); token != "" {
		parameters = append(parameters, http.WithExtraHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}))
	}

	return http.New(ctx, append(parameters, params...)...)
}

// initGlobalHTTPService creates a single HTTP service instance that all tests will share.
// This reduces connection overhead and makes tests more efficient.
func initGlobalHTTPService() {
	if os.Getenv("HTTP_ADDRESS") == "" {
		return
	}

	ctx := context.Background()

	// If we can't create the service, tests will fail anyway; leave the global nil and
	// let the individual tests handle the error.
	if service, err := newTestService(ctx, true, http.WithAllowDelayedStart(true)); err == nil {
		globalHTTPService = service
	}
}

// testService returns an HTTP service for testing.
// It returns the global shared service if available, otherwise creates a new one.
// Tests should use this function instead of creating their own service instances.
//
// This function also acquires a test coordination semaphore to limit concurrent
// test execution, preventing endpoint overload. The semaphore is automatically
// released when the test completes via t.Cleanup().
func testService(ctx context.Context, t *testing.T) any {
	// Acquire test coordinator semaphore to limit concurrent tests
	if testCoordinator != nil {
		if err := testCoordinator.Acquire(ctx, 1); err != nil {
			t.Fatalf("Failed to acquire test coordinator: %v", err)
		}
		// Release the semaphore when the test completes
		t.Cleanup(func() {
			testCoordinator.Release(1)
		})
	}

	if globalHTTPService != nil {
		return globalHTTPService
	}

	// Fallback: create a new service if global service is not available
	service, err := newTestService(ctx, true)
	if err != nil {
		t.Fatalf("Failed to create HTTP service: %v", err)
	}

	return service
}

// mustParseRoot is used for testing.
func mustParseRoot(input string) *phase0.Root {
	root, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic("invalid root")
	}
	if len(root) != phase0.RootLength {
		panic("invalid length root")
	}

	var res phase0.Root
	copy(res[:], root)

	return &res
}

// mustParseSignature is used for testing.
func mustParseSignature(input string) *phase0.BLSSignature {
	sig, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic("invalid signature")
	}
	if len(sig) != phase0.SignatureLength {
		panic("invalid length signature")
	}

	var res phase0.BLSSignature
	copy(res[:], sig)

	return &res
}

// mustParsePubKey is used for testing.
func mustParsePubKey(input string) *phase0.BLSPubKey {
	pubKey, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic("invalid public key")
	}
	if len(pubKey) != phase0.PublicKeyLength {
		panic("invalid length public key")
	}

	var res phase0.BLSPubKey
	copy(res[:], pubKey)

	return &res
}
