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

package api

import (
	"context"
	"time"

	"github.com/labstack/echo"
	admission "k8s.io/api/admission/v1beta1"
)

// Config defines an API server
type Config struct {
	// EnableLogging indicates we should enable http logging
	EnableLogging bool
	// EnableMetrics indicates we should expose prometheus metrics
	EnableMetrics bool
	// Listen in the interface to listen on
	Listen string
	// IdleTimeout is the timeout in idle connections
	IdleTimeout time.Duration
	// KeepAlive toogle keepalive connections
	KeepAlive time.Duration
	// ListenLimit is the overflow limit
	ListenLimit int
	// ReadTimeout is the http server read timeout
	ReadTimeout time.Duration
	// TLSCACertificate is the certificate authority
	TLSCACertificate string
	// TLSCert is the certificate to use for the service
	TLSCert string
	// TLSPrivateKey is the private key to use
	TLSPrivateKey string
	// WriteTimeout is the http write timeout
	WriteTimeout time.Duration

	// AdmissionHandlerFunc is the admission request handler
	AdmissionHandlerFunc func(context.Context, *admission.AdmissionReview) (*admission.AdmissionResponse, error)
	// HealthHandlerFunc returns the health of the service
	HealthHandlerFunc func(context.Context) (string, error)
}

// Server is the interface for an API service
type Server struct {
	// config is the configuration for the service
	config *Config
	// engine is the echo handler
	engine *echo.Echo
}
