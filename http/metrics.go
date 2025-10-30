// Copyright © 2023, 2024 Attestant Limited.
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

package http

import (
	"context"
	"errors"
	"regexp"

	"github.com/attestantio/go-eth2-client/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsMetric *prometheus.CounterVec
	stateMetric    *prometheus.GaugeVec
)

func registerMetrics(ctx context.Context, monitor metrics.Service) error {
	if requestsMetric != nil {
		// Already registered.
		return nil
	}

	if monitor == nil {
		// No monitor.
		return nil
	}

	if monitor.Presenter() == "prometheus" {
		return registerPrometheusMetrics(ctx)
	}

	return nil
}

func registerPrometheusMetrics(_ context.Context) error {
	requestsMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "consensusclient",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Number of requests",
	}, []string{"server", "method", "endpoint", "result"})
	if err := prometheus.Register(requestsMetric); err != nil {
		return errors.Join(errors.New("failed to register requests_total"), err)
	}

	stateMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "consensusclient",
		Subsystem: "http",
		Name:      "connection_state",
		Help:      "The state of the client connection (active/synced/inactive)",
	}, []string{"server", "state"})
	if err := prometheus.Register(stateMetric); err != nil {
		return errors.Join(errors.New("failed to register state"), err)
	}

	return nil
}

func (s *Service) monitorGetComplete(_ context.Context, endpoint string, result string) {
	if requestsMetric == nil {
		return
	}

	requestsMetric.WithLabelValues(s.address, "GET", reduceEndpoint(endpoint), result).Inc()
}

func (s *Service) monitorPostComplete(_ context.Context, endpoint string, result string) {
	if requestsMetric == nil {
		return
	}

	requestsMetric.WithLabelValues(s.address, "POST", reduceEndpoint(endpoint), result).Inc()
}

type templateReplacement struct {
	pattern     *regexp.Regexp
	replacement []byte
}

var endpointTemplates = []*templateReplacement{
	{
		pattern: regexp.MustCompile(
			"/(blinded_blocks|blob_sidecars|blocks|headers|sync_committee)/(0x[0-9a-fA-F]{64}|[0-9]+|head|genesis|finalized)",
		),
		replacement: []byte("/$1/{block_id}"),
	},
	{
		pattern:     regexp.MustCompile("/bootstrap/0x[0-9a-fA-F]{64}"),
		replacement: []byte("/bootstrap/{block_root}"),
	},
	{
		pattern:     regexp.MustCompile("/duties/(attester|proposer|sync)/[0-9]+"),
		replacement: []byte("/duties/$1/{epoch}"),
	},
	{
		pattern:     regexp.MustCompile("/peers/[0-9a-zA-Z]+"),
		replacement: []byte("/peers/{peer_id}"),
	},
	{
		pattern:     regexp.MustCompile("/rewards/attestations/[0-9]+"),
		replacement: []byte("/rewards/attestations/{epoch}"),
	},
	{
		pattern:     regexp.MustCompile("/states/(0x[0-9a-fA-F]{64}|[0-9]+|head|genesis|finalized)"),
		replacement: []byte("/states/{state_id}"),
	},
	{
		pattern:     regexp.MustCompile("/validators/(0x[0-9a-fA-F]{64}|[0-9]+)"),
		replacement: []byte("/validators/{validator_id}"),
	},
}

// reduceEndpoint reduces an endpoint to its template.
func reduceEndpoint(in string) string {
	out := []byte(in)
	for _, template := range endpointTemplates {
		out = template.pattern.ReplaceAll(out, template.replacement)
	}

	return string(out)
}

func (s *Service) monitorState(state string) {
	if stateMetric == nil {
		return
	}

	switch state {
	case "synced":
		stateMetric.WithLabelValues(s.address, "synced").Set(1)
		stateMetric.WithLabelValues(s.address, "active").Set(0)
		stateMetric.WithLabelValues(s.address, "inactive").Set(0)
	case "active":
		stateMetric.WithLabelValues(s.address, "synced").Set(0)
		stateMetric.WithLabelValues(s.address, "active").Set(1)
		stateMetric.WithLabelValues(s.address, "inactive").Set(0)
	case "inactive":
		stateMetric.WithLabelValues(s.address, "synced").Set(0)
		stateMetric.WithLabelValues(s.address, "active").Set(0)
		stateMetric.WithLabelValues(s.address, "inactive").Set(1)
	default:
		// Unknown state, do nothing
	}
}
