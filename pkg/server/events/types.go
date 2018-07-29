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

package events

import (
	admission "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Event is a denial event
type Event struct {
	// Detail is the detail about the error
	Detail string
	// Object is the decoded object
	Object metav1.Object
	// Review is a reference to the review
	Review *admission.AdmissionRequest
}

// Sink is the implementation for a events consumer
type Sink interface {
	// Send is responsible is sending messages
	Send(*Event) error
}
