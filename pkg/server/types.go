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
	"sync"

	"github.com/gambol99/kube-admission/pkg/authorize"
)

// Config is the configuration for the service
type Config struct {
	// EnableMetrics indicates we should expose prometheus metrics
	EnableMetrics bool `yaml:"enable-metrics" json:"enable-metrics"`
	// EnableLogging indicates we should enable http logging
	EnableLogging bool `yaml:"enable-logging" json:"enable-logging"`
	// Listen in the interface to listen on
	Listen string `yaml:"listen" json:"listen"`
	// TLSCert is the certificate to use for the service
	TLSCert string `yaml:"tls-cert" json:"tls-cert"`
	// TLSCA is the certificate authority
	TLSCA string `json:"tls-ca" yaml:"tls-ca"`
	// TLSPrivateKey is the private key to use
	TLSPrivateKey string `yaml:"tls-private-key" json:"tls-private-key"`
}

// Admission is the wrapper for the admission controller
type Admission struct {
	sync.RWMutex

	// authorizers is a collection of authorization scripts
	authorizers []authorize.Authorizer
}
