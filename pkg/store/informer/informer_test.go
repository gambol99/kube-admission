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

package informer

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNew(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shared := informers.NewSharedInformerFactoryWithOptions(fake.NewSimpleClientset(), 0)

	err := New(ctx, &Config{
		Factory:  shared,
		Resource: "v1/namespaces",
	})
	assert.NoError(t, err)
}

func TestNewUnknownResource(t *testing.T) {
	err := New(context.TODO(), &Config{Resource: "unkwown"})
	assert.Error(t, err)
}

func TestMultipleInformers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	client := fake.NewSimpleClientset()

	doneCh := make(chan struct{}, 0)
	errorCh := make(chan error, 0)
	factory := informers.NewSharedInformerFactoryWithOptions(client, 0)

	// @step: create the namespace informer
	assert.NoError(t, New(ctx, &Config{
		Factory:  factory,
		Resource: "v1/namespaces",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
		},
	}))

	assert.NoError(t, New(ctx, &Config{
		Factory:  factory,
		Resource: "v1/pods",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "test_pod", object.GetName())
			doneCh <- struct{}{}
		},
	}))

	// @step: add the namespace and pod
	_, err := client.Core().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NoError(t, err)

	_, err = client.Core().Pods("default").Create(&v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test_pod"}})
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func TestInformerCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	doneCh := make(chan struct{}, 0)
	errorCh := make(chan error, 0)
	client := newTestInformer(t, &Config{
		Resource: "v1/namespaces",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
			doneCh <- struct{}{}
		},
		ErrorFunc: func(version schema.GroupVersionResource, err error) {
			errorCh <- err
		},
	}, ctx)

	// @step: add a namespace and check we get a update
	ns, err := client.Core().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NotNil(t, ns)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func TestInformerUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	var updated int
	doneCh := make(chan struct{}, 0)
	errorCh := make(chan error, 0)
	client := newTestInformer(t, &Config{
		Resource: "v1/namespaces",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			updated++
		},
		UpdateFunc: func(version schema.GroupVersionResource, before, after metav1.Object) {
			updated++
			if updated == 3 {
				require.NotNil(t, before)
				require.NotNil(t, after)
				assert.Equal(t, "default", before.GetName())
				annotations := after.GetAnnotations()
				require.NotNil(t, annotations)
				assert.Equal(t, "default", after.GetName())
				assert.Equal(t, "test", after.GetAnnotations()["test"])
				doneCh <- struct{}{}
			}
		},
	}, ctx)

	// @step: add a namespace and check we get a update
	ns, err := client.Core().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NotNil(t, ns)
	require.NoError(t, err)
	ns.SetAnnotations(map[string]string{"test": "test"})
	_, err = client.Core().Namespaces().Update(ns)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func TestInformerDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	doneCh := make(chan struct{}, 0)
	errorCh := make(chan error, 0)
	client := newTestInformer(t, &Config{
		Resource: "v1/namespaces",
		DeleteFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
			doneCh <- struct{}{}
		},
	}, ctx)

	// @step: add a namespace and check we get a update
	_, err := client.Core().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NoError(t, err)
	require.NoError(t, client.Core().Namespaces().Delete("default", &metav1.DeleteOptions{}))

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func newTestInformer(t *testing.T, c *Config, ctx context.Context) kubernetes.Interface {
	client := fake.NewSimpleClientset()
	c.Factory = informers.NewSharedInformerFactoryWithOptions(client, 0)
	err := New(ctx, c)
	require.NoError(t, err)

	return client
}
