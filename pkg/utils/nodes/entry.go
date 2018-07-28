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

// NewEntity creates and returns a new entity
func NewEntity(name string) *Entity {
	return &Entity{
		Name:     name,
		Entities: make(map[string]*Entity, 0),
	}
}

// Get returns a child entity
func (e *Entity) Get(name string) (*Entity, bool) {
	e.RLock()
	defer e.RUnlock()

	if v, found := e.Entities[name]; found {
		return v, true
	}

	return nil, false
}

// Add a entity to the list of children
func (e *Entity) Add(entity *Entity) {
	e.Lock()
	defer e.Unlock()

	e.Entities[entity.Name] = entity
}

// Delete removes a child entity
func (e *Entity) Delete(name string) {
	e.Lock()
	defer e.Unlock()

	delete(e.Entities, name)
}

// List returns a list of named children
func (e *Entity) List() []string {
	e.RLock()
	defer e.RUnlock()

	var list []string
	for k, _ := range e.Entities {
		list = append(list, k)
	}

	return list
}

// Size returns the size of the entries
func (e *Entity) Size() int {
	e.RLock()
	defer e.RUnlock()

	return len(e.Entities)
}

// Values returns the values from the entity
func (e *Entity) Values() []interface{} {
	e.RLock()
	defer e.RUnlock()

	var list []interface{}
	for _, v := range e.Entities {
		list = append(list, v.Value)
	}

	return list
}
