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
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newTestStore(t *testing.T, resources ...string) *storeImpl {
	var list []string

	list = append(list, resources...)
	s, err := New(fake.NewSimpleClientset(), list)
	require.NotNil(t, s)
	require.NoError(t, err)

	return s.(*storeImpl)
}

func TestNew(t *testing.T) {
	client := fake.NewSimpleClientset()
	s, err := New(client, []string{})
	assert.NotNil(t, s)
	assert.NoError(t, err)
}

func TestStoreList(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 10; i++ {
		node := &unstructured.Unstructured{}
		node.SetName(fmt.Sprintf("nothing%d", i))

		s.Kind("nothing").Set(node.GetName(), node)
	}

	for i := 0; i < 10; i++ {
		node := &unstructured.Unstructured{}
		node.SetName(fmt.Sprintf("node%d", i))

		s.Kind("nodes").Set(node.GetName(), node)
	}

	items, err := s.Kind("nodes").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.NotEmpty(t, items)
	assert.Equal(t, 10, len(items))
}

func TestStoreGet(t *testing.T) {
	s := newTestStore(t)

	node := &unstructured.Unstructured{}
	node.SetName("node1")
	s.Kind("nodes").Set("node1", node)

	v, found, err := s.Kind("nodes").Get("node1")
	require.NoError(t, err)
	require.NotNil(t, v)
	require.True(t, found)
	assert.Equal(t, "node1", v.GetName())
}

func TestStoreSet(t *testing.T) {
	s := newTestStore(t)

	node := &unstructured.Unstructured{}
	node.SetName("node1")
	s.Kind("nodes").Set("node1", node)

	v, found, err := s.Kind("nodes").Get("node1")
	require.NoError(t, err)
	require.NotNil(t, v)
	require.True(t, found)
	assert.Equal(t, "node1", v.GetName())
}

func TestStoreActions(t *testing.T) {
	s := newTestStore(t)
	cs := []struct {
		Actions func()
		Checks  func(int)
	}{
		{
			Actions: func() {
				s.Kind("nodes").Set("node1", &unstructured.Unstructured{})
				s.Kind("nodes").Set("node2", &unstructured.Unstructured{})
			},
			Checks: func(i int) {
				found, err := s.Kind("nodes").Has("node2")
				assert.True(t, found)
				assert.NoError(t, err)
				items, err := s.Kind("nodes").List()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(items))
			},
		},
		{
			Actions: func() {
				s.Kind("nodes").Delete("node1")
			},
			Checks: func(i int) {
				found, err := s.Kind("nodes").Has("node1")
				assert.False(t, found)
				assert.NoError(t, err)
				items, err := s.Kind("nodes").List()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(items))
			},
		},
		{
			Actions: func() {
				s.Kind("namespaces").Set("default", &unstructured.Unstructured{})
			},
			Checks: func(i int) {
				found, err := s.Kind("namespaces").Has("default")
				assert.True(t, found)
				assert.NoError(t, err)
			},
		},
		{
			Actions: func() {
				s.Kind("namespaces").Set("default", &unstructured.Unstructured{})
				for _, x := range []string{"test0", "test1"} {
					s.Namespace("default").Kind("pods").Set(x, &unstructured.Unstructured{})
				}
			},
			Checks: func(i int) {
				found, err := s.Kind("namespaces").Has("default")
				require.True(t, found)
				require.NoError(t, err)
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(items))
			},
		},
		{
			Actions: func() {
				s.Namespace("default").Kind("pods").Delete("test0")
			},
			Checks: func(i int) {
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(items))
			},
		},
		{
			Actions: func() {
				s.Kind("namespaces").Delete("default")
			},
			Checks: func(i int) {
				found, err := s.Kind("namespaces").Has("default")
				require.False(t, found)
				require.NoError(t, err)
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 0, len(items))
			},
		},
	}
	for i, c := range cs {
		c.Actions()
		c.Checks(i)
	}
}
