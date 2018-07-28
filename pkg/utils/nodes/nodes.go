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
	"strings"
)

// New returns a nodetree
func New() *NodeTree {
	return &NodeTree{
		root: NewEntity("root"),
	}
}

// Set is responsible for adding a node
func (n *NodeTree) Set(key string, o interface{}) error {
	cursor := n.root
	for _, x := range strings.Split(key, "/") {
		if entry, found := cursor.Entities[x]; found {
			cursor = entry
			continue
		}
		entry := NewEntity(x)
		entry.Parent = cursor

		cursor.Add(entry)
		cursor = entry
	}
	cursor.Value = o

	return nil
}

// Delete is responsible removing a node
func (n *NodeTree) Delete(key string) (bool, error) {
	entry := n.find(key)
	if entry == nil {
		return false, nil
	}
	// @step: delete children if has any
	if entry.Size() > 0 {
		n.deleteAll(entry)
	}

	// @step: delete from the parent
	if entry.Parent != nil {
		// @check if we have any children
		entry.Parent.Delete(entry.Name)
	}

	return true, nil
}

// Has checks if the resource exists
func (n *NodeTree) Has(key string) (bool, error) {
	entry := n.find(key)
	if entry == nil {
		return false, nil
	}

	return true, nil
}

// Get is resposible for retrieving a node
func (n *NodeTree) Get(key string) (interface{}, bool, error) {
	if entry := n.find(key); entry != nil {
		return entry.Value, true, nil
	}

	return nil, false, nil
}

// List of a search of the nodes
func (n *NodeTree) List(key string) ([]interface{}, error) {
	var list []interface{}

	entry := n.find(key)
	if entry == nil {
		return list, nil
	}

	return entry.Values(), nil
}

// deleteAll is responseible for deleting all the entries
func (n *NodeTree) deleteAll(entry *Entity) {
	for k, v := range entry.Entities {
		if len(v.Entities) > 0 {
			n.deleteAll(v)
			return
		}
		entry.Delete(k)
	}
}

// find searchs for the entity
func (n *NodeTree) find(key string) *Entity {
	cursor := n.root
	for _, x := range strings.Split(key, "/") {
		entry, found := cursor.Entities[x]
		if found {
			cursor = entry
			continue
		}

		return nil
	}

	return cursor
}
