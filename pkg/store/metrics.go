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

package store

import "github.com/prometheus/client_golang/prometheus"

var (
	deleteCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_delete_counter",
			Help: "A counter or the delete operations in the store",
		},
	)
	deleteLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "delete_latency_sec",
			Help: "The latency on delete operations to the store",
		},
	)
	getLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "get_latency_sec",
			Help: "The latency on get operations to the store",
		},
	)
	hasLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "has_latency_sec",
			Help: "The latency on has operations to the store",
		},
	)
	listLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "list_latency_sec",
			Help: "The latency on list operations to the store",
		},
	)
	setLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "set_latency_sec",
			Help: "The latency on set operations to the store",
		},
	)
	updateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "store_update_counter",
			Help: "A counter or the update and add operations in the store",
		},
	)
)

func init() {
	prometheus.MustRegister(deleteCounter)
	prometheus.MustRegister(deleteLatency)
	prometheus.MustRegister(getLatency)
	prometheus.MustRegister(hasLatency)
	prometheus.MustRegister(listLatency)
	prometheus.MustRegister(setLatency)
	prometheus.MustRegister(updateCounter)
}
