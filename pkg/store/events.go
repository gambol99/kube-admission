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
	"github.com/gambol99/kube-admission/pkg/store/informer"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// handleObject is called by an informer when an resource has been changed, added or deleted
func (s *storeImpl) handleObject(version schema.GroupVersionResource, before, object metav1.Object, eventType EventType) {
	fields := log.Fields{
		"group":     version.Group,
		"kind":      version.Resource,
		"name":      object.GetName(),
		"namespace": object.GetNamespace(),
		"resource":  object.GetResourceVersion(),
		"version":   version.Version,
	}
	if eventType == Deleted {
		deleteCounter.Inc()
		if err := s.Namespace(object.GetNamespace()).Kind(version.Resource).Delete(object.GetName()); err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("unable to delete from the store")
		}
	} else {
		updateCounter.Inc()
		if err := s.Namespace(object.GetNamespace()).Kind(version.Resource).Set(object.GetName(), object); err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("unable to update or create resource in the store")
		}
	}

	s.handleEventListeners(version, before, object, eventType)
}

// handleEventListeners is responsible for handling the event listeners
func (s *storeImpl) handleEventListeners(version schema.GroupVersionResource, before, object metav1.Object, eventType EventType) {

	// @step: check if anyone is listening to this resource
	listeners, found := s.listeners[informer.NiceVersion(version)]
	if !found {
		return
	}

	// @step: fire of the events to the listeners
	event := &Event{After: object, Before: before, Type: eventType, Version: version}

	for _, listener := range listeners {
		if listener.Type != eventType {
			continue
		}
		go func() { listener.Channel <- event }()
	}
}
