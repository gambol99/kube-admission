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
	"fmt"
	"regexp"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// nameRegex is a regex to validate the naming parameter
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_\*]*$`)
)

// queryBuilder is a query builer for the store's index
type queryBuilder struct {
	Interface
	// kind is the api resource kind
	kind string
	// name is the name of the resource
	name string
	// namespace is scoped namespace
	namespace string
	// within is a resource create or modified before this time duration
	within time.Duration
	// the store contract
	store *storeImpl
	// version is the object version
	version string
}

// newQueryBuilder returns a query builder for us
func newQueryBuilder(store *storeImpl) *queryBuilder {
	return &queryBuilder{store: store}
}

// Namespace sets the query namespace
func (q *queryBuilder) Namespace(name string) Interface {
	q.namespace = name
	return q
}

// Kind sets the query resource kind
func (q *queryBuilder) Kind(name string) Interface {
	q.kind = name
	return q
}

// Within set a time limit for the objects i.e all object within the last 2 minutes
func (q *queryBuilder) Within(tm time.Duration) Interface {
	q.within = tm
	return q
}

// Delete removes a object from the store
func (q *queryBuilder) Delete(name string) error {
	q.name = name
	if name == "" {
		return errors.New("no name set")
	}

	return q.store.deleteObjectStore(q)
}

// Set adds an object to the store
func (q *queryBuilder) Set(name string, o metav1.Object) error {
	q.name = name
	if name == "" {
		return errors.New("no name set")
	}

	return q.store.updateObjectStore(q, o)
}

// Has checks if the resource exists in the store
func (q *queryBuilder) Has(name string) (bool, error) {
	if name == "" {
		return false, errors.New("no name set")
	}

	if v, err := q.Get(name); err != nil {
		return false, err
	} else if v != nil {
		return true, nil
	}

	return false, nil
}

// Get retrieves a resource from the store
func (q *queryBuilder) Get(name string) (metav1.Object, error) {
	if name == "" {
		return nil, errors.New("no name set")
	}
	q.name = name

	items, err := q.store.searchObjectStore(q, queryGet)
	if err != nil {
		return nil, err
	}
	if len(items) > 1 {
		return nil, errors.New("too many results returns for query")
	}
	if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

// List retrieves a list of resources from the store
func (q *queryBuilder) List() ([]metav1.Object, error) {
	return q.store.searchObjectStore(q, queryGet)
}

// buildQuery returns the query from the builder
func (q *queryBuilder) buildQuery() (string, error) {
	// @step: first check the parameters are valid
	if err := q.IsValid(); err != nil {
		return "", err
	}

	// @step: build the query string
	var terms []string
	if q.namespace != "" {
		terms = append(terms, "+namespace:"+q.namespace)
	}
	if q.kind != "" {
		terms = append(terms, "+kind:"+q.kind)
	}
	if q.version != "" {
		terms = append(terms, "+version:"+q.version)
	}
	if q.name != "" {
		terms = append(terms, "+name:"+q.name)
	}
	if q.within != 0 {
		terms = append(terms, "+modified:>"+fmt.Sprintf("%d", time.Now().Add(-q.within).Unix()))
	}

	return strings.Join(terms, " "), nil
}

// IsValid checks if the query is valid
func (q *queryBuilder) IsValid() error {
	if q.namespace != "" && !nameRegex.MatchString(q.namespace) {
		return fmt.Errorf("namespace: %q is invalid", q.namespace)
	}
	if q.kind != "" && !nameRegex.MatchString(q.kind) {
		return fmt.Errorf("kind: %q is invalid", q.kind)
	}
	if q.name != "" && !nameRegex.MatchString(q.name) {
		return fmt.Errorf("name: %q is invalid", q.name)
	}
	if q.kind == "" {
		return errors.New("resource kind not set")
	}

	return nil
}
