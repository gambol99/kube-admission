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
	"testing"

	"github.com/stretchr/testify/assert"
	admission "k8s.io/api/admission/v1beta1"
)

func TestIsValidConfig(t *testing.T) {
	hh := func(_ context.Context) (string, error) {
		return "", nil
	}
	ha := func(_ context.Context, _ *admission.AdmissionReview) (*admission.AdmissionResponse, error) {
		return nil, nil
	}

	cs := []struct {
		Config Config
		OK     bool
	}{
		{},
		{Config: Config{Listen: ":100"}},
		{Config: Config{Listen: ":100", HealthHandlerFunc: hh}},
		{Config: Config{Listen: ":100", HealthHandlerFunc: hh, AdmissionHandlerFunc: ha}, OK: true},
	}
	for i, c := range cs {
		err := c.Config.IsValid()
		switch c.OK {
		case true:
			assert.NoError(t, err, "case %d, did not expect an error", i)
		default:
			assert.Error(t, err, "case %d, expected an error", i)
		}
	}
}

func newTestConfig() *Config {
	c := &Config{
		EnableMetrics: true,
		Listen:        ":11000",
		HealthHandlerFunc: func(_ context.Context) (string, error) {
			return "ok", nil
		},
		AdmissionHandlerFunc: func(_ context.Context, _ *admission.AdmissionReview) (*admission.AdmissionResponse, error) {
			return &admission.AdmissionResponse{
				Allowed: true,
			}, nil
		},
	}

	return c
}
