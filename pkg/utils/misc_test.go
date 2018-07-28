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

package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	items := []string{"a", "b", "b", "c", "c", "d"}
	values := Unique(items)
	assert.NotEmpty(t, values)
	assert.Equal(t, 4, len(values))
	assert.Equal(t, []string{"a", "b", "c", "d"}, values)
}

func TestContains(t *testing.T) {
	items := []string{"a", "b", "c", "d"}
	assert.True(t, Contains("a", items))
	assert.True(t, Contains("c", items))
	assert.True(t, Contains("d", items))
	assert.False(t, Contains("z", items))
}

func TestSetID(t *testing.T) {
	ctx := SetID(context.Background(), "test")
	assert.Equal(t, "test", GetID(ctx))
}

func TestGetID(t *testing.T) {
	ctx := SetID(context.Background(), "test")
	assert.Equal(t, "test", GetID(ctx))
}
