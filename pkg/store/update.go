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
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// indexedDocument is the struct we add to the search index
type indexedDocument struct {
	Kind      string `json:"kind"`
	Modified  int64  `json:"modified"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
}

// updateObjectStore is responsible for updating the store
func (s *storeImpl) updateObjectStore(builder *queryBuilder, object metav1.Object) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	err := func() error {
		// @step: validate the query
		query, err := builder.buildQuery()
		if err != nil {
			return err
		}

		// @step: check if the resource already exists
		hits, err := s.search.Search(query)
		if err != nil {
			return err
		}

		// @check if were got results and get the uid for the document
		var uid string
		switch len(hits) {
		case 0:
			uid, err = s.getUID(5)
			if err != nil {
				return err
			}
		case 1:
			uid = hits[0]
		default:
			return errors.New("invalid set query, received more then one result")
		}

		// @step: create a document for indexing
		if err := s.search.Index(uid, &indexedDocument{
			Kind:      builder.kind,
			Modified:  time.Now().Unix(),
			Name:      builder.name,
			Namespace: builder.namespace,
			Version:   builder.version,
		}); err != nil {
			return err
		}

		// @step: add the document to the local cache
		s.cache.Set(uid, object, 0)

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("unable to add the document into search index")

		errorCounter.Inc()
	}

	return err
}

// getUID is responsible for generateing a unique uuid
func (s *storeImpl) getUID(attempts int) (string, error) {
	for i := 0; i < attempts; i++ {
		uid, err := uuid.NewV1()
		if err != nil {
			return "", err
		}

		// @check the uid not in
		if _, found := s.cache.Get(uid.String()); !found {
			return uid.String(), nil
		}
	}

	return "", errors.New("unable to generate uuid")
}
