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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UKHomeOffice/policy-admission/pkg/api"
	"github.com/docker/docker/api/server"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:   "policy-admission",
		Author: "Rohith Jayawardene",
		Email:  "gambol99@gmail.com",
		Usage:  "kubernetes admission controller webhook service",

		OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options, %s\n", err)
			return err
		},

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "listen",
				Usage:  "network interface the service should listen on `INTERFACE`",
				Value:  ":8443",
				EnvVar: "LISTEN",
			},
			cli.StringFlag{
				Name:   "tls-cert",
				Usage:  "file containing the tls certificate `PATH`",
				EnvVar: "TLS_CERT",
			},
			cli.StringFlag{
				Name:   "tls-key",
				Usage:  "file containing the tls private key `PATH`",
				EnvVar: "TLS_KEY",
			},
			cli.StringFlag{
				Name:   "tls-ca",
				Usage:  "file containing the certificate authority `PATH`",
				EnvVar: "TLS_CA",
			},
			cli.StringSliceFlag{
				Name:  "authorizer",
				Usage: "enable an admission authorizer, the format is name=config_path (i.e images=config.yaml)",
			},
			cli.StringSliceFlag{
				Name:  "config",
				Usage: "adding a configurable parameter to the controller (key=value)",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "namespace to create denial events (optional as we can try and discover) `NAME`",
				EnvVar: "KUBE_NAMESPACE",
				Value:  "kube-system",
			},
			cli.BoolFlag{
				Name:   "enable-logging",
				Usage:  "indicates you wish to log the admission requests for debugging `BOOL`",
				EnvVar: "ENABLE_LOGGING",
			},
			cli.BoolTFlag{
				Name:   "enable-metrics",
				Usage:  "indicates you wish to expose the prometheus metrics `BOOL`",
				EnvVar: "ENABLE_METRICS",
			},
			cli.BoolFlag{
				Name:   "enable-events",
				Usage:  "indicates you wish to log kubernetes events on denials `BOOL`",
				EnvVar: "ENABLE_EVENTS",
			},
			cli.DurationFlag{
				Name:   "http-idle-timeout",
				Usage:  "the time duration to attempt to wrap up duplicate events `DURATION`",
				EnvVar: "HTTP_IDLE_TIMEOUT",
				Value:  10 * time.Second,
			},
			cli.DurationFlag{
				Name:   "http-read-timeout",
				Usage:  "the time duration to attempt to wrap up duplicate events `DURATION`",
				EnvVar: "HTTP_IDLE_TIMEOUT",
				Value:  10 * time.Second,
			},
			cli.DurationFlag{
				Name:   "http-write-timeout",
				Usage:  "the time duration to attempt to wrap up duplicate events `DURATION`",
				EnvVar: "HTTP_IDLE_TIMEOUT",
				Value:  10 * time.Second,
			},
			cli.BoolTFlag{
				Name:   "http-keepalive",
				Usage:  "indicates you wish for verbose logging `BOOL`",
				EnvVar: "HTTP_KEEPALIVE",
			},
			cli.BoolFlag{
				Name:   "verbose",
				Usage:  "indicates you wish for verbose logging `BOOL`",
				EnvVar: "VERBOSE",
			},
		},

		Action: func(cx *cli.Context) error {
			var authorizers []api.Authorize
			// @step: configure the authorizers
			for _, config := range cx.StringSlice("authorizer") {
				authorizer, err := configureAuthorizer(config)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[error] unable to enable authorizer: %s", err)
					os.Exit(1)
				}
				authorizers = append(authorizers, authorizer)
			}

			config := &server.Config{
				ControllerName: cx.String("controller-name"),
				EnableEvents:   cx.Bool("enable-events"),
				EnableMetrics:  cx.Bool("enable-metrics"),
				EnableLogging:  cx.Bool("enable-logging"),
				Listen:         cx.String("listen"),
				Namespace:      cx.String("namespace"),
				TLSCert:        cx.String("tls-cert"),
				TLSKey:         cx.String("tls-key"),
				TLSCA:          cx.String("tls-ca"),
				Verbose:        cx.Bool("verbose"),
			}

			// @step: create the server
			ctl, err := server.New(config, authorizers)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] unable to initialize controller, %q\n", err)
				os.Exit(1)
			}

			// @step: start the service
			if err := ctl.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "[error] unable to start controller, %q\n", err)
				os.Exit(1)
			}

			// @step setup the termination signals
			signalChannel := make(chan os.Signal)
			signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			<-signalChannel

			return nil
		},
	}

	app.Run(os.Args)

}
