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
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func newTestQueryBuilder(t *testing.T) *queryBuilder {
	return newQueryBuilder(newTestStore(t))
}

func TestQueryBuilderBad(t *testing.T) {
	node := &unstructured.Unstructured{}
	cs := []struct {
		Action func()
	}{
		{
			Action: func() {
				c := newTestQueryBuilder(t).Namespace("%%.d").Kind("services")
				assert.Error(t, c.Set("test", node))
			},
		},
		{
			Action: func() {
				c := newTestQueryBuilder(t).Namespace("test/../").Kind("services")
				assert.Error(t, c.Set("test", node))
			},
		},
		{
			Action: func() {
				c := newTestQueryBuilder(t).Namespace("default").Kind("services/test")
				assert.Error(t, c.Set("test", node))
			},
		},
	}
	for _, c := range cs {
		c.Action()
	}
}

func TestBuildQuery(t *testing.T) {
	node := &unstructured.Unstructured{}
	cs := []struct {
		Client   func() *queryBuilder
		Expected string
	}{
		{
			Client: func() *queryBuilder {
				c := newTestQueryBuilder(t).Kind("namespaces")
				c.List()

				return c.(*queryBuilder)
			},
			Expected: "+kind:namespaces",
		},
		{
			Client: func() *queryBuilder {
				return newTestQueryBuilder(t).Namespace("default").Kind("services").(*queryBuilder)
			},
			Expected: "+namespace:default +kind:services",
		},
		{
			Client: func() *queryBuilder {
				c := newTestQueryBuilder(t).Namespace("default").Kind("services")
				c.Set("test", node)

				return c.(*queryBuilder)
			},
			Expected: "+namespace:default +kind:services +name:test",
		},
		{
			Client: func() *queryBuilder {
				c := newTestQueryBuilder(t).Kind("nodes")
				c.Set("test", node)

				return c.(*queryBuilder)
			},
			Expected: "+kind:nodes +name:test",
		},
	}
	for i, c := range cs {
		query, err := c.Client().buildQuery()
		require.NoError(t, err)
		assert.Equal(t, c.Expected, query, "case %d, expected: %s, got: %s", i, c.Expected, query)
	}
}
