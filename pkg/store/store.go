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
	"context"
	"errors"
	"time"

	"github.com/gambol99/kube-admission/pkg/store/indexer"
	"github.com/gambol99/kube-admission/pkg/store/informer"
	"github.com/gambol99/kube-admission/pkg/utils"

	pcache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// storeImpl is the implementation of the store
type storeImpl struct {
	// cache is local cache for object
	cache *pcache.Cache
	// client is the kubernetes client to use
	client kubernetes.Interface
	// factory is the shared informer factory interface
	factory informers.SharedInformerFactory
	// listeners is a list of event listeners
	listeners map[string][]*Listener
	// search is the search index for the resources
	search indexer.Interface
}

// New creates and returns a resource store: we recieve a kubernete clients and a list
// of resource type to watch and add into the store
func New(client kubernetes.Interface, resources []string) (Store, error) {
	// @step: create a search index
	search, err := indexer.New()
	if err != nil {
		return nil, err
	}

	// @step: create a the store service
	svc := &storeImpl{
		cache:     pcache.New(10*time.Minute, 5*time.Minute),
		client:    client,
		factory:   informers.NewSharedInformerFactoryWithOptions(client, 30*time.Second),
		listeners: make(map[string][]*Listener, 0),
		search:    search,
	}

	// @step: iterate the resources required and create informers
	for _, resource := range utils.Unique(resources) {
		log.Infof("attmepting to create a informer for resource: %s", resource)

		// @step: create a resource information to watch for changes
		if err := informer.New(context.TODO(), &informer.Config{
			Factory:  svc.factory,
			Resource: resource,

			// add actions methods
			AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
				svc.handleObject(version, nil, object, Added)
			},
			DeleteFunc: func(version schema.GroupVersionResource, object metav1.Object) {
				svc.handleObject(version, nil, object, Deleted)
			},
			UpdateFunc: func(version schema.GroupVersionResource, before, after metav1.Object) {
				svc.handleObject(version, before, after, Updated)
			},
			// add the downstream error method
			ErrorFunc: func(version schema.GroupVersionResource, err error) {
				log.WithFields(log.Fields{
					"error":   err.Error(),
					"kind":    version.Resource,
					"version": version.Version,
				}).Error("resource informer has encountered an error")

				errorCounter.Inc()
			},
		}); err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"resource": resource,
			}).Error("failed to create resource informer")

			return nil, err
		}
	}

	return svc, nil
}

// Register records a callback for upstream recievers
func (s *storeImpl) Register(listener *Listener) error {
	if listener.Channel == nil {
		return errors.New("listener event channel is nil")
	}

	key := informer.NiceVersion(listener.Resource)
	s.listeners[key] = append(s.listeners[key], listener)

	return nil
}

// Namespace returns a namespece scoped client request
func (s *storeImpl) Namespace(name string) Interface {
	q := newQueryBuilder(s)
	q.namespace = name

	return q
}

// Kind returns a kind scope client request
func (s *storeImpl) Kind(name string) Interface {
	q := newQueryBuilder(s)
	q.kind = name

	return q
}

// Close is responsible for releasing the resources
func (s *storeImpl) Close() error {
	return nil
}
