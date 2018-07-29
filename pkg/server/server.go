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

package server

import (
	"context"
)

const (
	// The request has been accepted
	actionAccepted = "accept"
	// The request has been refused
	actionDenied = "deny"
	// The request has cause an error
	actionErrored = "error"
)

// New creates and returns an admission controller
func New(c *Config) (*Admission, error) {

	return &Admission{}, nil
}

// health provides information on the health of the service
func (a *Admission) health() ([]byte, error) {
	return []byte("ok"), nil
}

// Run is responsible for starting the admission controller
func (a *Admission) Run(ctc context.Context) error {

	return nil
}
