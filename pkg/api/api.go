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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// New creates and returns an API server
func New(config *Config) (*Server, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	svc := &Server{config: config}

	// @step: create the http router
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	if config.EnableLogging {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.GET("/health", svc.healthHandler)
	e.POST("/", svc.admissionHandler)
	if config.EnableMetrics {
		e.GET("/metrics", func(ctx echo.Context) error {
			prometheus.Handler().ServeHTTP(ctx.Response().Writer, ctx.Request())
			return nil
		})
	}
	svc.engine = e

	return svc, nil
}

// Run is responsible for starting the API server up
func (a *Server) Run(ctx context.Context) error {
	// @step: create the http server
	hs := &http.Server{
		Addr:         a.config.Listen,
		Handler:      a.engine,
		IdleTimeout:  a.config.IdleTimeout,
		ReadTimeout:  a.config.ReadTimeout,
		WriteTimeout: a.config.WriteTimeout,
	}

	listener, err := net.Listen("tcp", a.config.Listen)
	if err != nil {
		return err
	}

	// @step: load any tls certificates
	if a.config.TLSCert != "" && a.config.TLSPrivateKey != "" {
		tlsConfig := &tls.Config{
			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			// https://github.com/golang/go/tree/master/src/crypto/elliptic
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			// Use modern tls mode https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
			NextProtos: []string{"http/1.1", "h2"},
			// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
			MinVersion: tls.VersionTLS12,
			// These ciphersuites support Forward Secrecy: https://en.wikipedia.org/wiki/Forward_secrecy
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}

		cert, err := tls.LoadX509KeyPair(a.config.TLSCert, a.config.TLSPrivateKey)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}

		if a.config.TLSCACertificate != "" {
			caCert, caCertErr := ioutil.ReadFile(a.config.TLSCACertificate)
			if caCertErr != nil {
				return err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.ClientCAs = caCertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		hs.TLSConfig = tlsConfig

		listener = tls.NewListener(listener, tlsConfig)
	}

	go func() {
		if err := hs.Serve(listener); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to create the http server")
		}
	}()

	return nil
}

// Handler returns the http handler for the service
func (a *Server) Handler() http.Handler {
	return a.engine
}
