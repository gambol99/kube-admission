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

package authorize

import (
	"context"

	"github.com/gambol99/kube-admission/pkg/store"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Authorizer is a wrapper to a scripted authorizer
type Authorizer interface {
	// Name is the name of authorizer
	Name() string
	// Description is a short description of what the authorizer does
	Description() string
	// Filter is the filter for the authorizer i.e. what it wants to recieve
	Filter() *[]Filter
	// Admit is called to run the authorizer
	Admit(context.Context, metav1.Object, store.Store) (*Decision, error)
}

// Filter is used to control what a authorizer wishes to see
type Filter struct {
	// Group is the API group of the resource
	Group string `yaml:"group" json:"group"`
	// Kind is the kind of resource
	Kind string `yaml:"kind" json:"kind"`
	// Operations is the type of operations i.e. DELETE, UPDATE, CREATE etc
	Operations []string `yaml:"operations" json:"operations"`
	// Version is the API version of the resource
	Version string `yaml:"version" json:"version"`
}

// Reason is the reason why something was rejected
type Reason struct {
	// Field is the offending field
	Field string
	// Message is a human readable reason
	Message string
	// Value is value which was rejected
	Value interface{}
}

// Decision is the decision of the authorizer
type Decision struct {
	// Permitted indicates the resource is permitted
	Permitted bool
	// Reason is the reason why it was rejected
	Reasons []Reason
}
