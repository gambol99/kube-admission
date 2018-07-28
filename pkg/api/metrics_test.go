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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	hs := makeTestAPIServer(t, nil)
	require.NotNil(t, hs)

	resp, err := http.Get(hs.URL + "/metrics")
	require.NoError(t, err)
	require.NotNil(t, resp)

	content, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(content), "go_gc_duration_seconds")
}
