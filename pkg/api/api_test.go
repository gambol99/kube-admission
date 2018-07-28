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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	a, err := New(newTestConfig())
	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestHandler(t *testing.T) {
	a, err := New(newTestConfig())
	require.NoError(t, err)
	require.NotNil(t, a)
	assert.NotNil(t, a.Handler())
}

func TestRun(t *testing.T) {
	a, err := New(newTestConfig())
	require.NoError(t, err)
	require.NotNil(t, a)
	err = a.Run(context.TODO())
	assert.NoError(t, err)
}

// makeTestAPIServer creates and returns a API server for us
func makeTestAPIServer(t *testing.T, config *Config) *httptest.Server {
	if config == nil {
		config = newTestConfig()
	}
	require.NotNil(t, config)

	svc, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, svc)

	return httptest.NewServer(svc.Handler())
}
