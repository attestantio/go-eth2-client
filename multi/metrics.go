// Copyright © 2021 Attestant Limited.
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

var providersMetric *prometheus.GaugeVec
var providerStateMetric *prometheus.GaugeVec

func registerMetrics(ctx context.Context, monitor metrics.Service) error {
	if providersMetric != nil {
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

func registerPrometheusMetrics(ctx context.Context) error {
	providersMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eth2client",
		Subsystem: "multi",
		Name:      "providers_total",
		Help:      "Number of providers",
	}, []string{"state"})
	if err := prometheus.Register(providersMetric); err != nil {
		return errors.Wrap(err, "failed to register providers_total")
	}
	providerStateMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "eth2client",
		Subsystem: "multi",
		Name:      "provider_state",
		Help:      "State of provider",
	}, []string{"provider"})
	if err := prometheus.Register(providerStateMetric); err != nil {
		return errors.Wrap(err, "failed to register provider_state")
	}

	return nil
}

func setProviderStateMetric(ctx context.Context, provider string, state string) {
	if providerStateMetric != nil {
		if state == "active" {
			providerStateMetric.WithLabelValues(provider).Set(1)
		} else {
			providerStateMetric.WithLabelValues(provider).Set(0)
		}
	}
}

func setProvidersMetric(ctx context.Context, state string, count int) {
	if providersMetric != nil {
		providersMetric.WithLabelValues(state).Set(float64(count))
	}
}
