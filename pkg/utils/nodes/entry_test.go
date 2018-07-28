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

package nodes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntity(t *testing.T) {
	e := NewEntity("test")
	assert.NotNil(t, e)
	assert.Equal(t, "test", e.Name)
}

func TestAdd(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	e.Add(NewEntity("node"))
}

func TestGet(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	e.Add(NewEntity("node"))
	entry, found := e.Get("node")
	assert.NotNil(t, entry)
	assert.True(t, found)
	assert.Equal(t, "node", entry.Name)
}

func TestGetMultiple(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	for i := 0; i < 10; i++ {
		e.Add(NewEntity(fmt.Sprintf("node%d", i)))
	}

	entry, found := e.Get("node7")
	require.NotNil(t, entry)
	assert.True(t, found)
	assert.Equal(t, "node7", entry.Name)
}

func TestList(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	for i := 0; i < 10; i++ {
		e.Add(NewEntity(fmt.Sprintf("node%d", i)))
	}

	entities := e.List()
	assert.NotEmpty(t, entities)
	assert.Equal(t, 10, len(entities))
}

func TestDelete(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	e.Add(NewEntity("test"))
	entry, found := e.Get("test")
	require.NotNil(t, entry)
	require.True(t, found)

	e.Delete("test")
	entry, found = e.Get("test")
	assert.Nil(t, entry)
	assert.False(t, found)
}

func TestSize(t *testing.T) {
	e := NewEntity("root")
	require.NotNil(t, e)

	for i := 0; i < 10; i++ {
		e.Add(NewEntity(fmt.Sprintf("nodes/node%d", i)))
	}

	assert.Equal(t, 10, e.Size())
}
