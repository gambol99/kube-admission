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

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/blevesearch/bleve"
	"github.com/prometheus/client_golang/prometheus"
)

// indexer is the service wrapper
type indexer struct {
	// store is the index interface
	store bleve.Index
}

// New returns a memory only index
func New() (Interface, error) {
	store, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}

	return &indexer{store: store}, nil
}

// NewIndexFromDisk returns an index backed to a disk
func NewIndexFromDisk(name string) (Interface, error) {
	return nil, errors.New("currently unsupported")
}

// Delete is responsible for deleting a document from the index
func (i *indexer) Delete(id string) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return i.store.Delete(id)
}

// DeleteByQuery deletes all the documents which match the query
func (i *indexer) DeleteByQuery(q string) (int, error) {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	hits, err := i.Search(q)
	if err != nil {
		return 0, err
	}
	if len(hits) <= 0 {
		return 0, nil
	}

	for index, id := range hits {
		if err := i.Delete(id); err != nil {
			return index, err
		}
	}

	return len(hits), nil
}

// Index is responsible is add a document the index
func (i *indexer) Index(id string, doc interface{}) error {
	timed := prometheus.NewTimer(indexLatency)
	defer timed.ObserveDuration()

	return i.store.Index(id, doc)
}

// Search is responsible for searching the index
func (i *indexer) Search(query string) ([]string, error) {
	timed := prometheus.NewTimer(searchLatency)
	defer timed.ObserveDuration()

	var list []string

	resp, err := i.store.Search(bleve.NewSearchRequest(bleve.NewQueryStringQuery(query)))
	if err != nil {
		return list, err
	}

	for _, x := range resp.Hits {
		list = append(list, x.ID)
	}

	return list, nil
}

// Query is responsible for searching the index
func (i *indexer) Query(search interface{}) ([]string, error) {

	if reflect.ValueOf(search).Kind() != reflect.Struct {
		return []string{}, errors.New("search must be a type struct")
	}
	var terms []string

	v := reflect.ValueOf(search).Elem()
	// @step: interatet the fields for the search struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Type().Kind() {
		case reflect.Bool, reflect.String, reflect.Uint64, reflect.Int64, reflect.Int:
		default:
			continue
		}
		value := reflect.ValueOf(field.Interface())
		if value.String() == "" {
			continue
		}
		terms = append(terms, fmt.Sprintf("+%s:%s", v.Type().Field(i).Name, value.String()))
	}

	return []string{}, nil
}

// Size returns the size of the index
func (i *indexer) Size() (uint64, error) {
	return i.store.DocCount()
}
