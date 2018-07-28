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
	"errors"
	"fmt"
	"reflect"

	"github.com/gambol99/kube-admission/pkg/utils"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
)

// New creates and returns a resource informer
func New(ctx context.Context, config *Config) error {
	// @check the resource is supported
	if !utils.Contains(config.Resource, ResourceNames()) {
		return fmt.Errorf("resource: %s is not a supported resource type", config.Resource)
	}

	version, _ := ResourceVersions()[config.Resource]

	// @step: create the informer
	informer, err := config.Factory.ForResource(version)
	if err != nil {
		return err
	}

	svc := &resourceInformer{informer: informer, config: config, resourceVersion: version}

	return svc.start(ctx)
}

// start is responsible for running the informing and updating the caches
func (r *resourceInformer) start(ctx context.Context) error {
	r.informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(before interface{}) {
			addCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleAddObject(before)
		},
		DeleteFunc: func(before interface{}) {
			deleteCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleDeleteObject(before)
		},
		UpdateFunc: func(before, after interface{}) {
			updateCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleUpdateObject(before, after)
		},
	})

	// @step: start the shared index informer
	stopCh := make(chan struct{}, 0)
	go r.informer.Informer().Run(stopCh)

	log.Infof("synchronizing the cache for resources: %s", r.config.Resource)

	if !cache.WaitForCacheSync(stopCh, r.informer.Informer().HasSynced) {
		runtime.HandleError(fmt.Errorf("controller timed out waiting for caches to sync"))

		return fmt.Errorf("controller timed out waiting for cache sync")
	}

	// @step: wait for a signal to stop
	go func() {
		select {
		case <-ctx.Done():
			log.Infof("closing down in the informer for resource: %s", r.config.Resource)
			close(stopCh)
		}
	}()

	return nil
}

// handleAddObject is responsible for handling the deletions
func (r *resourceInformer) handleAddObject(before interface{}) {
	if r.config.AddFunc == nil {
		return
	}

	object, err := ensureMetaObject(before)
	if err != nil {
		r.handleInformerError(err)
		return
	}

	r.config.AddFunc(r.resourceVersion, object)
}

// handleDeleteObject is responsible for handling the deletions
func (r *resourceInformer) handleDeleteObject(before interface{}) {
	if r.config.DeleteFunc == nil {
		return
	}

	object, err := ensureMetaObject(before)
	if err != nil {
		r.handleInformerError(err)
		return
	}

	r.config.DeleteFunc(r.resourceVersion, object)
}

// handleUpdateObject is resposible for handling an updated object
func (r *resourceInformer) handleUpdateObject(before, after interface{}) {
	if r.config.UpdateFunc == nil {
		return
	}
	b, err := ensureMetaObject(before)
	if err != nil {
		r.handleInformerError(err)
		return
	}
	a, err := ensureMetaObject(after)
	if err != nil {
		r.handleInformerError(err)
	}

	r.config.UpdateFunc(r.resourceVersion, b, a)
}

// handleInformerError is responsible for pushing the error upstream
func (r *resourceInformer) handleInformerError(err error) {
	go func() {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("informer encountered an error")

		r.config.ErrorCh <- err
	}()
}

// ensureMetaObject checks to make sure the object is a meta object for us
func ensureMetaObject(object interface{}) (metav1.Object, error) {
	// @check the object is not nil
	if object == nil {
		return nil, errors.New("object expected: metav1.Object, got: nil")
	}

	// @check the object implements the metav.Object interface
	if _, ok := object.(metav1.Object); !ok {
		return nil, fmt.Errorf("object resource is not as expected: , got: %s", reflect.TypeOf(object).String())
	}

	return object.(metav1.Object), nil
}
