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

package kube

import (
	"testing"

	"github.com/gambol99/kube-admission/pkg/server/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNew(t *testing.T) {
	s, err := New(fake.NewSimpleClientset(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestSend(t *testing.T) {
	client := fake.NewSimpleClientset()
	s, err := New(client, "test")
	require.NoError(t, err)
	require.NotNil(t, s)

	object := &unstructured.Unstructured{}
	object.SetNamespace("default")

	resp, err := client.Core().Events(object.GetNamespace()).List(metav1.ListOptions{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, len(resp.Items))

	event := &events.Event{Detail: "this is a event", Object: object}
	require.NoError(t, s.Send(event))

	resp, err = client.Core().Events(object.GetNamespace()).List(metav1.ListOptions{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, len(resp.Items))
}
