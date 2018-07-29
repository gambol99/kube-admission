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

package utils

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// requireID is a request unique id
type requestID string

const requestIDName requestID = "uuid"

// GetID returns the request transaction id from the context
func GetID(c context.Context) string {
	v := c.Value(requestIDName)
	if v == nil {
		return ""
	}

	return v.(string)
}

// SetID sets the uuid of the context
func SetID(c context.Context, id string) context.Context {
	return context.WithValue(c, requestIDName, id)
}

// Contains checks if a item exists
func Contains(v string, list []string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}

	return false
}

// Unique removes any duplicates from the list
func Unique(e []string) []string {
	var list []string

	items := make(map[string]bool, 0)
	for _, x := range e {
		if _, found := items[x]; found {
			continue
		}
		items[x] = true
		list = append(list, x)
	}

	return list
}

// GetKubernetesClient returns a kubernetes api client for us
func GetKubernetesClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
