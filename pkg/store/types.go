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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// Store is the contract to the store
type Store interface {
	// Close releases the resources
	Close() error
	// Namespace returns operations for a namespace
	Namespace(string) Interface
	// Kind return a client request scoped to the kind
	Kind(string) Interface
	// Register is used to register an event listener to resources
	Register(*Listener) error
}

// EventType indicates the type of event
type EventType int

const (
	// Added indicates the resource was created
	Added EventType = iota
	// Updated indicates the resource was updated
	Updated
	// Deleted indicates the resources was remove
	Deleted
)

// Listener defined a upstream listener
type Listener struct {
	// Channel to recieve the events on
	Channel chan *Event
	// Type is the type of event
	Type EventType
	// Resource is what you want to listen to
	Resource schema.GroupVersionResource
}

// Event defines are to callback the origin
type Event struct {
	// Type indicates the event type
	Type EventType
	// Before was the object before
	Before metav1.Object
	// After was the object after
	After metav1.Object
	// Version is the resource version
	Version schema.GroupVersionResource
}

// Interface defines a list of action verbs
type Interface interface {
	// Delete removes a object from the store
	Delete(string) (bool, error)
	// Has checks if the resource exists in the store
	Has(string) (bool, error)
	// Get retrieves a resource from the store
	Get(string) (metav1.Object, bool, error)
	// List retrieves a list of resources from the store
	List() ([]metav1.Object, error)
	// Set adds a resource to the store
	Set(string, metav1.Object) error
	// Kind adds the api kind type to the request
	Kind(string) Interface
	// Namespace is used to set the namespace
	Namespace(string) Interface
}
