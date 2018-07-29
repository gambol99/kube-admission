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

package indexer

import "github.com/prometheus/client_golang/prometheus"

var (
	deleteLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "indexer_delete_latency_sec",
			Help: "The latency on get operations to the store",
		},
	)
	indexLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "indexer_index_latency_sec",
			Help: "The latency on index operations in the index",
		},
	)
	searchLatency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "indexer_search_latency_sec",
			Help: "The latency on search operations in the index",
		},
	)
)

func init() {
	prometheus.MustRegister(deleteLatency)
	prometheus.MustRegister(indexLatency)
	prometheus.MustRegister(searchLatency)
}
