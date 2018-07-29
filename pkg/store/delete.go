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
)

// deleteObjectStore is responsible for deleting an object from the store
func (s *storeImpl) deleteObjectStore(builder *queryBuilder) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	err := func() error {
		// @check if the query is valid
		query, err := builder.buildQuery()
		if err != nil {
			return err
		}
		// @step: search for object matching the query
		items, err := s.searchStoreIndex(query)
		if err != nil {
			return err
		}
		if len(items) <= 0 {
			return nil
		}
		// @step: delete the items from the search index
		if err := s.deleteSearchIndex(query); err != nil {
			return err
		}
		// @step: delete items from the cache
		for _, key := range items {
			s.cache.Delete(key)
		}

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to delete from object store")

		errorCounter.Inc()
	}

	return err
}

// deleteSearchIndex is responsible for deleting object from the index
func (s *storeImpl) deleteSearchIndex(query string) error {
	_, err := s.search.DeleteByQuery(query)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"query": query,
		}).Error("failed to remove items from search index")

		errorCounter.Inc()
	}

	return err
}
