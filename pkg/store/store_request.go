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
	"fmt"
	"regexp"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// validNameFilter is a regex to validate the naming parameter
	valueNameFilter = regexp.MustCompile(`^[a-zA-Z0-9\-_]*$`)
)

// request is used to hold the scope of a request
type request struct {
	// kind is the api resource kind
	kind string
	// name is the name of the resource
	name string
	// namespace is scoped namespace
	namespace string
	// store is the underlining store references
	store *storeImpl
	// version is the object version
	version string
}

// Delete removes a object from the store
func (c *request) Delete(name string) (bool, error) {
	c.name = name
	if err := c.isValidRequest(); err != nil {
		return false, err
	}
	n := time.Now()
	defer deleteLatency.Observe(time.Since(n).Seconds())

	return c.store.tree.Delete(c.buildKey())
}

// Has checks if the resource exists in the store
func (c *request) Has(name string) (bool, error) {
	c.name = name
	if err := c.isValidRequest(); err != nil {
		return false, err
	}
	n := time.Now()
	defer hasLatency.Observe(time.Since(n).Seconds())

	return c.store.tree.Has(c.buildKey())
}

// Get retrieves a resource from the store
func (c *request) Get(name string) (metav1.Object, bool, error) {
	c.name = name
	if err := c.isValidRequest(); err != nil {
		return nil, false, err
	}
	n := time.Now()
	defer getLatency.Observe(time.Since(n).Seconds())

	entry, found, err := c.store.tree.Get(c.buildKey())
	if err != nil || !found {
		return nil, false, err
	}

	return entry.(metav1.Object), true, nil
}

// List retrieves a list of resources from the store
func (c *request) List() ([]metav1.Object, error) {
	var list []metav1.Object

	if err := c.isValidRequest(); err != nil {
		return list, err
	}
	n := time.Now()
	defer listLatency.Observe(time.Since(n).Seconds())

	items, err := c.store.tree.List(c.buildKey())
	if err != nil {
		return list, err
	}
	for _, x := range items {
		list = append(list, x.(metav1.Object))
	}

	return list, nil
}

// Set adds a resource to the store
func (c *request) Set(name string, o metav1.Object) error {
	c.name = name
	if err := c.isValidRequest(); err != nil {
		return err
	}
	n := time.Now()
	defer setLatency.Observe(time.Since(n).Seconds())

	return c.store.tree.Set(c.buildKey(), o)
}

// Kind adds the api kind type to the request
func (c *request) Kind(name string) Interface {
	c.kind = name

	return c
}

// Namespaces add the namespace
func (c *request) Namespace(name string) Interface {
	c.namespace = name

	return c
}

// isValidRequest checks each of the parameters
func (c *request) isValidRequest() error {
	if c.namespace != "" && !valueNameFilter.MatchString(c.namespace) {
		return fmt.Errorf("namespace: %q is invalid", c.namespace)
	}
	if c.kind != "" && !valueNameFilter.MatchString(c.kind) {
		return fmt.Errorf("kind: %q is invalid", c.kind)
	}
	if c.name != "" && !valueNameFilter.MatchString(c.name) {
		return fmt.Errorf("name: %q is invalid", c.name)
	}

	return nil
}

func (c *request) buildKey() string {
	var paths []string
	if c.namespace != "" {
		paths = append(paths, []string{"namespaces", c.namespace}...)
	}
	if c.kind != "" {
		paths = append(paths, c.kind)
	}
	if c.version != "" {
		paths = append(paths, c.version)
	}
	if c.name != "" {
		paths = append(paths, c.name)
	}

	return strings.Join(paths, "/")
}
