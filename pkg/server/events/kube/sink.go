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
	"fmt"

	"github.com/gambol99/kube-admission/pkg/server/events"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type sink struct {
	// client is the kubernete client
	client kubernetes.Interface
	// name is the controller name
	name string
}

// New returns a kubernetes
func New(client kubernetes.Interface, name string) (events.Sink, error) {
	return &sink{client: client}, nil
}

// Send is responsible for sending the event into the kubernete events
func (s *sink) Send(e *events.Event) error {
	message := fmt.Sprintf("Denied in namespace: '%s', e.Object: '%s', reason: %s",
		e.Object.GetNamespace(),
		e.Object.GetGenerateName(),
		e.Detail)

	_, err := s.client.CoreV1().Events(e.Object.GetNamespace()).Create(&core.Event{
		Message: message,
		Reason:  "Forbidden",
		Source:  core.EventSource{Component: s.name},
		Type:    "Warning",
	})

	return err
}
