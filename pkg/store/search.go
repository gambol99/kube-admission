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

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// searchObjectStore is resposible for handling a search of the store for objects
func (s *storeImpl) searchObjectStore(builder *queryBuilder, qtype queryType) ([]metav1.Object, error) {
	var list []metav1.Object

	// @check if the query is valid
	query, err := builder.buildQuery()
	if err != nil {
		return list, err
	}

	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	// @step: query the index for the objects
	hits, err := s.searchStoreIndex(query)
	if err != nil {
		return list, err
	}

	// @step: lookup the resources from the cache
	for _, x := range hits {
		if o, found := s.cache.Get(x); !found {
			log.WithFields(log.Fields{"query": query}).Warn("cache key was not found in local store cache")
			errorCounter.Inc()
		} else {
			list = append(list, o.(metav1.Object))
		}
	}

	return list, nil
}

// searchStoreIndex is responsible for searching the search index
func (s *storeImpl) searchStoreIndex(query string) ([]string, error) {
	hits, err := s.search.Search(query)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"query": query,
		}).Error("failed to query the search index")

		errorCounter.Inc()
	}

	return hits, err
}
