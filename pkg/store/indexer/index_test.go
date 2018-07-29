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
	"fmt"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type document struct {
	Kind      string `json:"kind"`
	Modified  int64  `json:"modified"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func newUUID() string {
	id, _ := uuid.NewV1()

	return id.String()
}

func newTestIndex(t *testing.T) *indexer {
	i, err := New()
	require.NotNil(t, i)
	require.NoError(t, err)

	namespaces := []string{"default", "test", "frontend"}
	services := []string{"svc1", "svc2"}
	pods := []string{"test1", "test2"}

	for _, x := range namespaces {
		require.NoError(t, i.Index(newUUID(), &document{Kind: "namespace", Name: x, Modified: time.Now().Unix()}))
	}

	for _, n := range namespaces {
		for _, p := range pods {
			require.NoError(t, i.Index(newUUID(), &document{Kind: "pod", Name: p, Namespace: n, Modified: time.Now().Unix()}))
		}
	}

	for _, n := range namespaces {
		for _, p := range services {
			require.NoError(t, i.Index(newUUID(), &document{Kind: "service", Name: p, Namespace: n, Modified: time.Now().Unix()}))
		}
	}

	return i.(*indexer)
}

func TestNew(t *testing.T) {
	i, err := New()
	assert.NotNil(t, i)
	assert.NoError(t, err)
}

func TestNewIndexFromDisk(t *testing.T) {
	i, err := NewIndexFromDisk("none")
	assert.Nil(t, i)
	assert.Error(t, err)
}

func TestSize(t *testing.T) {
	i := newTestIndex(t)
	size, err := i.Size()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0xf), size)
}

func TestIndex(t *testing.T) {
	_ = newTestIndex(t)
}

func TestDelete(t *testing.T) {
	i := newTestIndex(t)
	hits, err := i.Search("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.NotEmpty(t, hits)
	require.Equal(t, 1, len(hits))

	require.NoError(t, i.Delete(hits[0]))
	hits, err = i.Search("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.Empty(t, hits)
	require.Equal(t, 0, len(hits))
}

func TestDeleteNamespaceByQuery(t *testing.T) {
	i := newTestIndex(t)
	query := "+namespace:default"
	hits, err := i.Search(query)
	require.NoError(t, err)
	require.NotEmpty(t, hits)
	require.Equal(t, 4, len(hits))

	num, err := i.DeleteByQuery(query)
	require.NoError(t, err)
	require.Equal(t, 4, num)
	hits, err = i.Search(query)
	require.NoError(t, err)
	require.Empty(t, hits)
	require.Equal(t, 0, len(hits))
}

func TestQuery(t *testing.T) {
	cs := []struct {
		Query         string
		ExpectedCount int
	}{
		{
			Query:         "+kind:pod[s]",
			ExpectedCount: 6,
		},
		{
			Query:         "+namespace:default",
			ExpectedCount: 4,
		},
		{
			Query:         "+namespace:default +kind:pod[s]",
			ExpectedCount: 2,
		},
		{
			Query:         "+namespace:default +kind:pod[s] +name=test1",
			ExpectedCount: 1,
		},
		{
			Query: "+namespace:default +kind:pod[s] +name=test0",
		},
		{
			Query:         "+namespace:default +modified:>" + fmt.Sprintf("%d", time.Now().Add(-1*time.Minute).Unix()),
			ExpectedCount: 4,
		},
	}
	i := newTestIndex(t)
	for x, c := range cs {
		resp, err := i.Search(c.Query)
		require.NoError(t, err)
		assert.Equal(t, c.ExpectedCount, len(resp), "case %d, query: %q, expected: %d, got: %d",
			x, c.Query, c.ExpectedCount, len(resp))
	}
}
