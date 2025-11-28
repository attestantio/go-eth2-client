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
var globalHTTPService interface{}

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

	if os.Getenv("HTTP_ADDRESS") != "" {
		// Initialize global HTTP service for all tests to share
		initGlobalHTTPService()
		os.Exit(m.Run())
	}
}

// initGlobalHTTPService creates a single HTTP service instance that all tests will share.
// This reduces connection overhead and makes tests more efficient.
func initGlobalHTTPService() {
	if os.Getenv("HTTP_ADDRESS") == "" {
		return
	}

	ctx := context.Background()
	var service client.Service
	var err error
	if os.Getenv("HTTP_BEARER_TOKEN") != "" {
		service, err = http.New(ctx,
			http.WithTimeout(timeout),
			http.WithAddress(os.Getenv("HTTP_ADDRESS")),
			http.WithAllowDelayedStart(true),
			http.WithExtraHeaders(map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("HTTP_BEARER_TOKEN")),
			}),
		)
	} else {
		service, err = http.New(ctx,
			http.WithTimeout(timeout),
			http.WithAddress(os.Getenv("HTTP_ADDRESS")),
			http.WithAllowDelayedStart(true),
		)
	}

	if err != nil {
		// If we can't create the service, tests will fail anyway
		// Just log and continue - individual tests will handle the error
		return
	}
	globalHTTPService = service
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
	var service client.Service
	var err error
	if os.Getenv("HTTP_BEARER_TOKEN") != "" {
		service, err = http.New(ctx,
			http.WithTimeout(timeout),
			http.WithAddress(os.Getenv("HTTP_ADDRESS")),
			http.WithExtraHeaders(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("HTTP_BEARER_TOKEN"))}),
		)
	} else {
		service, err = http.New(ctx,
			http.WithTimeout(timeout),
			http.WithAddress(os.Getenv("HTTP_ADDRESS")),
		)
	}

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
