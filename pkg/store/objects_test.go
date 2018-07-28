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
	"time"

	"github.com/gambol99/kube-admission/pkg/utils/informer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddNamespace(t *testing.T) {
	s := newTestStore(t, "v1/namespaces")
	eventCh := make(chan *Event, 0)
	err := s.Register(&Listener{
		Channel:  eventCh,
		Resource: informer.ResourceVersions()["v1/namespaces"],
		Type:     Added,
	})
	require.NoError(t, err)

	ns, err := s.client.Core().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	})
	require.NotNil(t, ns)
	require.NoError(t, err)

	select {
	case e := <-eventCh:
		require.NotNil(t, e)
		assert.Equal(t, e.Type, Added)
		require.NotNil(t, e.After)
		assert.Equal(t, "test", e.After.GetName())
	case <-time.After(time.Millisecond * 100):
		t.Error("failed to update the store on namespace change")
	}
}

func TestDeletedNamespace(t *testing.T) {
	s := newTestStore(t, "v1/namespaces")
	eventCh := make(chan *Event, 0)
	err := s.Register(&Listener{
		Channel:  eventCh,
		Resource: informer.ResourceVersions()["v1/namespaces"],
		Type:     Deleted,
	})
	require.NoError(t, err)

	ns, err := s.client.Core().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	})
	require.NotNil(t, ns)
	require.NoError(t, err)

	err = s.client.Core().Namespaces().Delete("test", &metav1.DeleteOptions{})
	require.NoError(t, err)

	select {
	case e := <-eventCh:
		require.NotNil(t, e)
		assert.Equal(t, e.Type, Deleted)
		require.NotNil(t, e.After)
		assert.Equal(t, "test", e.After.GetName())
	case <-time.After(time.Millisecond * 100):
		t.Error("failed to delete the store on namespace change")
	}
}
