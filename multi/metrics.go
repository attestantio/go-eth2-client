// Copyright © 2021, 2024 Attestant Limited.
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

package multi

import (
	"context"

	"github.com/attestantio/go-eth2-client/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	connectionsMetric *prometheus.GaugeVec
	stateMetric       *prometheus.GaugeVec
)

func registerMetrics(ctx context.Context, monitor metrics.Service) error {
	if connectionsMetric != nil {
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
	connectionsMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "consensusclient",
		Subsystem: "multi",
		Name:      "connections",
		Help:      "Number of connections",
	}, []string{"name", "state"})
	if err := prometheus.Register(connectionsMetric); err != nil {
		return errors.Wrap(err, "failed to register connections")
	}
	stateMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "consensusclient",
		Subsystem: "multi",
		Name:      "connection_state",
		Help:      "The state of the client connection (active/inactive)",
	}, []string{"name", "server", "state"})
	if err := prometheus.Register(stateMetric); err != nil {
		return errors.Wrap(err, "failed to register connection_state")
	}

	return nil
}

func (s *Service) setProviderStateMetric(_ context.Context, server string, state string) {
	if stateMetric == nil {
		return
	}

	switch state {
	case "active":
		stateMetric.WithLabelValues(s.name, server, "active").Set(1)
		stateMetric.WithLabelValues(s.name, server, "inactive").Set(0)
	case "inactive":
		stateMetric.WithLabelValues(s.name, server, "active").Set(0)
		stateMetric.WithLabelValues(s.name, server, "inactive").Set(1)
	}
}

func (s *Service) setConnectionsMetric(_ context.Context, active int, inactive int) {
	if connectionsMetric == nil {
		return
	}

	connectionsMetric.WithLabelValues(s.name, "active").Set(float64(active))
	connectionsMetric.WithLabelValues(s.name, "inactive").Set(float64(inactive))
}
