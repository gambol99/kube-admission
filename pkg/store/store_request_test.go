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
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func newStoreRequest(t *testing.T) *request {
	return &request{
		store: newTestStore(t),
	}
}

func TestClientRequestNamingOK(t *testing.T) {
	node := &unstructured.Unstructured{}
	cs := []struct {
		Action func()
	}{
		{
			Action: func() {
				c := newStoreRequest(t).Namespace("%%.d").Kind("services")
				assert.Error(t, c.Set("test", node))
			},
		},
		{
			Action: func() {
				c := newStoreRequest(t).Namespace("test/../").Kind("services")
				assert.Error(t, c.Set("test", node))
			},
		},
		{
			Action: func() {
				c := newStoreRequest(t).Namespace("default").Kind("services/test")
				assert.Error(t, c.Set("test", node))
			},
		},
	}
	for _, c := range cs {
		c.Action()
	}
}

func TestClientRequest(t *testing.T) {
	node := &unstructured.Unstructured{}
	cs := []struct {
		Client   func() *request
		Expected string
	}{
		{
			Client: func() *request {
				return newStoreRequest(t).Namespace("default").Kind("services").(*request)
			},
			Expected: "namespaces/default/services",
		},
		{
			Client: func() *request {
				c := newStoreRequest(t).Namespace("default").Kind("services")
				c.Set("test", node)

				return c.(*request)
			},
			Expected: "namespaces/default/services/test",
		},
		{
			Client: func() *request {
				c := newStoreRequest(t).Kind("nodes")
				c.Set("test", node)

				return c.(*request)
			},
			Expected: "nodes/test",
		},
	}
	for i, c := range cs {
		key := c.Client().buildKey()
		assert.Equal(t, c.Expected, key, "case %d, expected: %s, got: %s", i, c.Expected, key)
	}
}
