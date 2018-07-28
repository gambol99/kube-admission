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
	"sync"

	"github.com/armon/go-radix"
)

// Keystore is used to store the resources
type Keystore struct {
	sync.RWMutex
	// tree is the underlining resources
	tree *radix.Tree
}

// New creates and returns a keystore which implements the Verb interface
func New() *Keystore {
	return &Keystore{tree: radix.New()}
}

// Delete removes from the keys store tree
func (k *Keystore) Delete(key string) {
	k.Lock()
	defer k.Unlock()

	k.tree.Delete(key)
}

// Dump returns a map of the keys
func (k *Keystore) Dump() map[string]interface{} {
	return k.tree.ToMap()
}

// Get retrieves an resource from the store
func (k *Keystore) Get(key string) interface{} {
	k.RLock()
	defer k.RUnlock()

	v, found := k.tree.Get(key)
	if !found {
		return nil
	}

	return v
}

// Has checks if the store has the key
func (k *Keystore) Has(key string) bool {
	k.RLock()
	defer k.RUnlock()

	_, found := k.tree.Get(key)

	return found
}

// List returns a list of nodes from the tree
func (k *Keystore) List(key string) []interface{} {
	k.RLock()
	defer k.RUnlock()

	var list []interface{}
	k.tree.WalkPrefix(key, func(key string, value interface{}) bool {
		list = append(list, value)

		return false
	})

	return list
}

// Set adds a object to the store
func (k *Keystore) Set(key string, o interface{}) {
	k.Lock()
	defer k.Unlock()

	k.tree.Insert(key, o)
}
