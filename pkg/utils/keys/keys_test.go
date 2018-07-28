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

package keys

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeystore(t *testing.T) {
	assert.NotNil(t, New())
}

func TestKeystoreSet(t *testing.T) {
	k := New()
	require.NotNil(t, k)
	k.Set("core/v1/namespace", true)

	v := k.Get("core/v1/namespace")
	assert.NotNil(t, v)
}

func TestKeystoreDelete(t *testing.T) {
	k := New()
	key := "core/v1/namespace"
	require.NotNil(t, k)
	k.Set(key, true)

	v := k.Get(key)
	assert.NotNil(t, v)

	k.Delete(key)

	v = k.Get(key)
	assert.Nil(t, v)
}

func TestKeystoreList(t *testing.T) {
	k := New()
	require.NotNil(t, k)

	prefix := "core/v1/"
	for _, x := range []string{"nodes", "namespace"} {
		for i := 0; i < 10; i++ {
			k.Set(fmt.Sprintf("%s/%s/%d", prefix, x, i), true)
		}
	}

	items := k.List(prefix + "/nodes")
	assert.NotEmpty(t, items)
	assert.Equal(t, 10, len(items))

	items = k.List(prefix + "/")
	assert.Equal(t, 20, len(items))
}

func TestKeystoreGet(t *testing.T) {
	k := New()
	require.NotNil(t, k)
	k.Set("core/v1/namespace", true)

	v := k.Get("core/v1/namespace")
	assert.NotNil(t, v)
}
