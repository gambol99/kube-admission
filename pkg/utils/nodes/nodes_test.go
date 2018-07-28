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

package nodes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreeNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestTreeSet(t *testing.T) {
	e := New()
	require.NotNil(t, e)

	for i := 0; i < 5; i++ {
		assert.NoError(t, e.Set(fmt.Sprintf("nodes/node%d", i), true))
	}
}

func TestTreeGetOK(t *testing.T) {
	e := newTestTree(t)
	cs := []struct {
		Key string
		Ok  bool
	}{
		{
			Key: "nodes",
			Ok:  true,
		},
		{
			Key: "bad",
		},
		{
			Key: "nodes/node3",
			Ok:  true,
		},
	}
	for i, c := range cs {
		_, found, err := e.Get(c.Key)
		require.NoError(t, err)
		assert.Equal(t, c.Ok, found, "case %d, expected %t, got: %t", i, c.Ok, found)
	}
}

func TestTreeList(t *testing.T) {
	e := newTestTree(t)
	cs := []struct {
		Key  string
		Size int
	}{
		{
			Key:  "nodes",
			Size: 5,
		},
		{
			Key: "bad",
		},
		{
			Key:  "namespaces",
			Size: 2,
		},
		{
			Key:  "namespaces/default/pods",
			Size: 10,
		},
	}
	for i, c := range cs {
		items, err := e.List(c.Key)
		require.NoError(t, err)
		assert.Equal(t, c.Size, len(items), "case %d, expected %d, got: %d", i, c.Size, len(items))
	}
}

func TestTreeSingleDelete(t *testing.T) {
	e := newTestTree(t)

	items, err := e.List("namespaces/default/pods")
	require.NoError(t, err)
	require.Equal(t, 10, len(items))
	e.Delete("namespaces/default/pods/pod4")
	items, err = e.List("namespaces/default/pods")
	require.NoError(t, err)
	assert.Equal(t, 9, len(items))
}

func TestTreeMultipleDelete(t *testing.T) {
	e := newTestTree(t)

	items, err := e.List("namespaces/default/pods")
	require.NoError(t, err)
	require.Equal(t, 10, len(items))
	e.Delete("namespaces/default/pods/pod4")
	e.Delete("namespaces/default/pods/pod9")
	items, err = e.List("namespaces/default/pods")
	require.NoError(t, err)
	assert.Equal(t, 8, len(items))
}

func TestTreeDeleteRoot(t *testing.T) {
	e := newTestTree(t)

	items, err := e.List("namespaces/default/pods")
	require.NoError(t, err)
	require.Equal(t, 10, len(items))
	e.Delete("namespaces/default/pods")
	items, err = e.List("namespaces/default/pods")
	require.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

func TestTreeDeleteNamespace(t *testing.T) {
	e := newTestTree(t)

	items, err := e.List("namespaces")
	require.NoError(t, err)
	require.Equal(t, 2, len(items))
	e.Delete("namespaces/default")
	items, err = e.List("namespaces/default/pods")
	assert.Equal(t, 0, len(items))
	items, err = e.List("namespaces")
	require.NoError(t, err)
	assert.Equal(t, 1, len(items))
}

func newTestTree(t *testing.T) *NodeTree {
	e := New()
	require.NotNil(t, e)

	// add some nodes
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("nodes/node%d", i)
		e.Set(name, true)
	}

	// add some pods
	for _, x := range []string{"default", "test"} {
		for i := 0; i < 10; i++ {
			name := fmt.Sprintf("namespaces/%s/pods/pod%d", x, i)
			e.Set(name, true)
		}
	}

	// add some services
	for _, x := range []string{"default", "test"} {
		for i := 0; i < 2; i++ {
			name := fmt.Sprintf("namespaces/%s/services/svc%d", x, i)
			e.Set(name, true)
		}
	}

	return e
}
