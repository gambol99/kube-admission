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
	"fmt"
	"time"

	"github.com/gambol99/kube-admission/pkg/utils"
	"github.com/gambol99/kube-admission/pkg/utils/informer"
	"github.com/gambol99/kube-admission/pkg/utils/nodes"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// storeImpl is the implementation of the store
type storeImpl struct {
	// client is the kubernetes client to use
	client kubernetes.Interface
	// informerErrorCh is a channel to recieve errors from the informers
	informerErrorCh chan error
	// informerCancels cancelations
	informerCancels []context.CancelFunc
	// factory is the shared informer factory interface
	factory informers.SharedInformerFactory
	// listeners is a list of event listeners
	listeners map[string][]*Listener
	// tree is the underlining nodes tree
	tree *nodes.NodeTree
}

// New creates and returns a resource store: we recieve a kubernete clients and a list
// of resource type to watch and add into the store
func New(client kubernetes.Interface, resources []string) (Store, error) {
	// @step: create a the store service
	svc := &storeImpl{
		client:          client,
		factory:         informers.NewSharedInformerFactoryWithOptions(client, 30*time.Second),
		informerErrorCh: make(chan error, 0),
		listeners:       make(map[string][]*Listener, 0),
		tree:            nodes.New(),
	}

	// @step: iterate the resources required and create informers
	for _, resource := range utils.Unique(resources) {
		log.Infof("attmepting to create a informer for resource: %s", resource)

		// @step: create a context for the informer to run under
		ctx, cancel := context.WithCancel(context.Background())

		// @step: create a resource information to watch for changes
		if err := informer.New(ctx, &informer.Config{
			Factory:  svc.factory,
			Resource: resource,
			ErrorCh:  svc.informerErrorCh,
			AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
				svc.handleObject(version, nil, object, Added)
			},
			DeleteFunc: func(version schema.GroupVersionResource, object metav1.Object) {
				svc.handleObject(version, nil, object, Deleted)
			},
			UpdateFunc: func(version schema.GroupVersionResource, before, after metav1.Object) {
				svc.handleObject(version, before, after, Updated)
			},
		}); err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"resource": resource,
			}).Error("failed to create resource informer")

			return nil, err
		}
		svc.informerCancels = append(svc.informerCancels, cancel)
	}

	return svc, nil
}

// Register records a callback for upstream recievers
func (s *storeImpl) Register(listener *Listener) error {
	if listener.Channel == nil {
		return errors.New("listener event channel is nil")
	}

	key := s.versionKey(listener.Resource)
	s.listeners[key] = append(s.listeners[key], listener)

	return nil
}

// Namespace returns a namespece scoped client request
func (s *storeImpl) Namespace(name string) Interface {
	request := &request{store: s}

	return request.Namespace(name)
}

// Kind returns a kind scope client request
func (s *storeImpl) Kind(name string) Interface {
	request := &request{store: s}

	return request.Kind(name)
}

// Close is responsible for releasing the resources
func (s *storeImpl) Close() error {
	for _, cancel := range s.informerCancels {
		cancel()
	}

	return nil
}

// versionKey returns a key for the resource
func (s *storeImpl) versionKey(v schema.GroupVersionResource) string {
	return fmt.Sprintf("%s/%s/%s", v.Group, v.Version, v.Resource)
}
