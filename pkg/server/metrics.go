/*
Copyright 2018 Rohith Jayawardene <gambol99@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import "github.com/prometheus/client_golang/prometheus"

var (
	admissionErrorMetrics = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "admission_error_total",
			Help: "The number of error encountered by the admission",
		},
	)
	admissionLatencyMetric = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "admission_latency_sec",
			Help: "The request latency for incoming resource authorization",
		},
	)
	authorizerLatencyMetric = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "authorizer_latency_sec",
			Help: "The authorizer latency for incoming resource authorization",
		},
		[]string{"name"},
	)
	authorizerActionMetrics = prometheus.NewCounterVec(
		prometheus.SummaryOpts{
			Name: "authorizer_action_total",
			Help: "The authorizer actions broken down by script authorizer",
		},
		[]string{"name", "action"},
	)
)

func init() {
	prometheus.MustRegister(admissionErrorMetrics)
	prometheus.MustRegister(admissionLatencyMetric)
	prometheus.MustRegister(authorizerLatencyMetric)
	prometheus.MustRegister(authorizerActionMetrics)
}
